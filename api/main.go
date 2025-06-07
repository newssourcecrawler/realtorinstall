package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"encoding/csv"
	"io"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/newssourcecrawler/realtorinstall/api/handlers"
	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
	"github.com/newssourcecrawler/realtorinstall/dbmigrations"
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

	// 5. Build Gin router with CORS + JWT middleware
	router := gin.Default()
	router.Use(cors.Default())

	// 6. Authentication routes
	router.POST("/login", authH.Login)
	router.Use(AuthMiddleware(authSvc, userRepo))
	router.POST("/register",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "register_user"),
		authH.Register,
	)

	// 7. User CRUD routes
	router.GET("/users", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_user"),
		userH.List,
	)
	router.POST("/users", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_user"),
		userH.Create,
	)
	router.PUT("/users/:id", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_user"),
		userH.Update,
	)
	router.DELETE("/users/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_user"),
		userH.Delete,
	)

	// 8. Property routes
	router.GET("/properties", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_property"),
		propH.List,
	)
	router.POST("/properties",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_property"),
		propH.Create,
	)
	router.PUT("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_property"),
		propH.Update,
	)
	router.DELETE("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_property"),
		propH.Delete,
	)

	// 9. Buyer routes
	router.GET("/buyers",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_buyer"),
		buyerH.List,
	)
	router.POST("/buyers",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_buyer"),
		buyerH.Create,
	)
	router.PUT("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_buyer"),
		buyerH.Update,
	)
	router.DELETE("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_buyer"),
		buyerH.Delete,
	)

	// 10. Pricing routes
	router.GET("/pricing",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_pricing"),
		priceH.List,
	)
	router.POST("/pricing",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_pricing"),
		priceH.Create,
	)
	router.PUT("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_pricing"),
		priceH.Update,
	)
	router.DELETE("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_pricing"),
		priceH.Delete,
	)

	// 11. Sales routes
	router.GET("/sales", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_sale"),
		salesH.List,
	)
	router.POST("/sales",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		salesH.Create,
	)
	router.PUT("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_sale"),
		salesH.Update,
	)
	router.DELETE("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_sale"),
		salesH.Delete,
	)

	// 12. Introduction routes
	router.GET("/introductions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_introduction"),
		introH.List,
	)
	router.POST("/introductions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_introduction"),
		introH.Create,
	)
	router.PUT("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_introduction"),
		introH.Update,
	)
	router.DELETE("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_introduction"),
		introH.Delete,
	)

	// 13. Lettings routes
	router.GET("/lettings",
		RequirePermission(userRepo, "view_lettings"),
		lettingsH.List
	)
	router.POST("/lettings",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Create,
	)
	router.PUT("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Update,
	)
	router.DELETE("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Delete,
	)

	// 14. Plan routes
	router.GET("/plans", 
		RequirePermission(userRepo, "view_plans"),
		planH.List
	)
	router.POST("/plans",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Create,
	)
	router.PUT("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Update,
	)
	router.DELETE("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Delete,
	)

	// 15. Installment routes
	router.GET("/installments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments"),
		instH.List,
	)
	router.GET("/installments/plan/:planId",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments_byplan"),
		instH.ListByPlan,
	)
	router.POST("/installments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_installments"),
		instH.Create,
	)
	router.PUT("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_installments"),
		instH.Update,
	)
	router.DELETE("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_installments"),
		instH.Delete,
	)

	// 16. Payment routes
	router.GET("/payments", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_payments"),
		payH.List,
	)
	router.POST("/payments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_payments"),
		payH.Create,
	)
	router.PUT("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_payments"),
		payH.Update,
	)
	router.DELETE("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_payments"),
		payH.Delete,
	)

	// 17. Commission routes
	router.GET("/commissions", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_commission"),
		commissionH.List,
	)
	router.POST("/commissions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_commission"),
		commissionH.Create,
	)
	router.PUT("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_commission"),
		commissionH.Update,
	)
	router.DELETE("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_commission"),
		commissionH.Delete,
	)

	// 18. Reporting routes
	router.GET("/reports/commissions/beneficiary",
		RequirePermission(userRepo, "view_commissions_report"),
		reportH.CommissionByBeneficiary,
	)

	router.GET("/reports/installments/outstanding",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments_report"),
		reportH.OutstandingInstallmentsByPlan,
	)

	router.GET("/reports/sales/monthly",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_sales_report"),
		reportH.MonthlySalesVolume,
	)

	router.GET("/reports/lettings/rentroll",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_lettings_report"),
		reportH.ActiveLettingsRentRoll,
	)

	router.GET("/reports/properties/top-payments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_property_payments_report"),
		reportH.TopPropertiesByPaymentVolume,
	)

	router.GET("/export/properties.csv", ExportPropertiesCSV(db))

	router.GET("/export/properties.csv", ImportPropertiesCSV(db))

	// 19. Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8443",
		Handler: router,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	  }
	  go func() {
		fmt.Println("API listening on https://localhost:8443")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

// AuthMiddleware extracts tenantID + userID from JWT in Authorization header.
func AuthMiddleware(authSvc *apiServices.AuthService, userRepo apiRepos.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if !strings.HasPrefix(hdr, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(hdr, "Bearer ")
		claims, err := authSvc.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// Check license expiry if you embedded LicenseExp in JWT (optional)
		// if claims.LicenseExp < time.Now().Unix() {
		//     c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "license expired"})
		//     return
		// }

		c.Set("currentUser", claims.UserID)
		c.Set("currentTenant", claims.TenantID)
		c.Set("perms", claims.Permissions)
		c.Next()
	}
}

// RequirePermission checks that the logged‐in user's role is one of the allowed list.
func RequirePermission(userRepo apiRepos.UserRepo, allowed ...string) gin.HandlerFunc {
	isAllowed := func(role string) bool {
		for _, r := range allowed {
			if r == role {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		userID := c.GetInt64("currentUser")
		tenantID := c.GetString("currentTenant")
		user, err := userRepo.GetByID(context.Background(), tenantID, userID)
		if err != nil || user.Deleted || !isAllowed(user.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := dbmigrations.ApplyMigrations(db); err != nil {
		return nil, err
	}
	return db, nil
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
            var id int64; var addr string; var price float64
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