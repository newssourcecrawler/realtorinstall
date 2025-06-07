package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/newssourcecrawler/realtorinstall/api/handlers"
	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
)

func main() {
	// 1. Ensure data folder exists
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("mkdir data: %w", err))
	}

	userDB, _ := openDB("data/users.db")
	commissionDB, _ := openDB("data/commissions.db")
	propDB, _ := openDB("data/properties.db")
	pricingDB, _ := openDB("data/pricing.db")
	buyerDB, _ := openDB("data/buyers.db")
	planDB, _ := openDB("data/plans.db")
	instDB, _ := openDB("data/installments.db")
	payDB, _ := openDB("data/payments.db")
	salesDB, _ := openDB("data/sales.db")
	introDB, _ := openDB("data/introductions.db")
	lettingsDB, _ := openDB("data/lettings.db")

	// 2. Initialize repositories (one per domain)
	userRepo, _ := apiRepos.NewSQLiteUserRepo(userDB)
	commissionRepo, _ := apiRepos.NewSQLiteCommissionRepo(commissionDB)
	propRepo, _ := apiRepos.NewSQLitePropertyRepo(propDB)
	pricingRepo, _ := apiRepos.NewSQLiteLocationPricingRepo(pricingDB)
	buyerRepo, _ := apiRepos.NewSQLiteBuyerRepo(buyerDB)
	planRepo, _ := apiRepos.NewSQLiteInstallmentPlanRepo(planDB)
	instRepo, _ := apiRepos.NewSQLiteInstallmentRepo(instDB)
	payRepo, _ := apiRepos.NewSQLitePaymentRepo(payDB)
	salesRepo, _ := apiRepos.NewSQLiteSalesRepo(salesDB)
	introRepo, _ := apiRepos.NewSQLiteIntroductionsRepo(introDB)
	lettingsRepo, _ := apiRepos.NewSQLiteLettingsRepo(lettingsDB)

	// 3. Construct services
	jwtSecret := os.Getenv("APP_JWT_SECRET")
	if jwtSecret == "" {
		panic("APP_JWT_SECRET must be set")
	}
	authSvc := apiServices.NewAuthService(userRepo, jwtSecret, time.Hour*24)

	propSvc := apiServices.NewPropertyService(propRepo, userRepo, pricingRepo)
	buyerSvc := apiServices.NewBuyerService(buyerRepo)
	pricingSvc := apiServices.NewPricingService(pricingRepo)
	planSvc := apiServices.NewPlanService(planRepo, instRepo)
	instSvc := apiServices.NewInstallmentService(instRepo, payRepo)
	paySvc := apiServices.NewPaymentService(payRepo, instSvc)
	userSvc := apiServices.NewUserService(userRepo)
	salesSvc := apiServices.NewSalesService(salesRepo)
	introSvc := apiServices.NewIntroductionsService(introRepo, salesRepo, lettingsRepo, userRepo)
	lettingsSvc := apiServices.NewLettingsService(lettingsRepo)
	commissionSvc := apiServices.NewCommissionService(commissionRepo, salesRepo, lettingsRepo, introRepo, userRepo)

	reportSvc := apiServices.NewReportService(commissionRepo)

	// 4. Instantiate handlers
	authH := handlers.NewAuthHandler(authSvc)
	propH := handlers.NewPropertyHandler(propSvc)
	buyerH := handlers.NewBuyerHandler(buyerSvc)
	priceH := handlers.NewPricingHandler(pricingSvc)
	planH := handlers.NewPlanHandler(planSvc)
	instH := handlers.NewInstallmentHandler(instSvc)
	payH := handlers.NewPaymentHandler(paySvc)
	userH := handlers.NewUserHandler(userSvc)
	salesH := handlers.NewSalesHandler(salesSvc)
	introH := handlers.NewIntroductionsHandler(introSvc)
	lettingsH := handlers.NewLettingsHandler(lettingsSvc)
	commissionH := handlers.NewCommissionHandler(commissionSvc)
	reportH := handlers.NewReportHandler(reportSvc)

	// 19. Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:      ":8443",
		Handler:   router,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	go func() {
		fmt.Println("API listening on https://localhost:8443")
		if err := srv.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Listen error: %v\n", err)
			os.Exit(1)
		}
	}()

	// 20. Graceful shutdown on SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server forced to shutdown: %v\n", err)
	}
	fmt.Println("API server stopped.")
}

func ExportPropertiesCSV(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, address, base_price_usd FROM properties WHERE deleted=0`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Disposition", "attachment;filename=properties.csv")
		w.Header().Set("Content-Type", "text/csv")
		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Header
		writer.Write([]string{"ID", "Address", "BasePriceUSD"})

		for rows.Next() {
			var id int64
			var addr string
			var price float64
			rows.Scan(&id, &addr, &price)
			writer.Write([]string{
				fmt.Sprint(id),
				addr,
				fmt.Sprintf("%.2f", price),
			})
		}
	}
}

func ImportPropertiesCSV(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()
		reader := csv.NewReader(file)
		// Skip header
		if _, err := reader.Read(); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		tx, _ := db.Begin()
		stmt, _ := tx.Prepare(`INSERT INTO properties(address, base_price_usd) VALUES(?,?)`)
		defer stmt.Close()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 400)
				return
			}
			stmt.Exec(record[1], record[2])
		}
		tx.Commit()
		w.WriteHeader(http.StatusCreated)
	}
}

func ExportPropertiesXLSX(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := excelize.NewFile()
		rows, _ := db.Query(`SELECT id, address, base_price_usd FROM properties WHERE deleted=0`)
		defer rows.Close()
		sheet := "Properties"
		f.SetSheetName("Sheet1", sheet)
		f.SetCellValue(sheet, "A1", "ID")
		f.SetCellValue(sheet, "B1", "Address")
		f.SetCellValue(sheet, "C1", "BasePriceUSD")
		rowNum := 2
		for rows.Next() {
			var id int64
			var addr string
			var price float64
			rows.Scan(&id, &addr, &price)
			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), id)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), addr)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), price)
			rowNum++
		}
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=properties.xlsx")
		f.Write(w)
	}
}
