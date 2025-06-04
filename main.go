// main.go (project root)
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
	"github.com/newssourcecrawler/realtorinstall/internal/utils"
)

func main() {
	// 1. Load config (must include DatabasePath)
	cfg, err := utils.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Ensure data directory exists (SQLite files live here)
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data folder: %v", err)
	}

	// 3. Initialize all SQLite repos
	propertyRepo, err := repos.NewSQLitePropertyRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open PropertyRepo: %v", err)
	}
	buyerRepo, err := repos.NewSQLiteBuyerRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open BuyerRepo: %v", err)
	}
	planRepo, err := repos.NewSQLitePlanRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open PlanRepo: %v", err)
	}
	installmentRepo, err := repos.NewSQLiteInstallmentRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open InstallmentRepo: %v", err)
	}
	paymentRepo, err := repos.NewSQLitePaymentRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open PaymentRepo: %v", err)
	}
	pricingRepo, err := repos.NewSQLiteLocationPricingRepo(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open LocationPricingRepo: %v", err)
	}

	// 4. Create service instances
	propertyService := services.NewPropertyService(propertyRepo, pricingRepo)
	buyerService := services.NewBuyerService(buyerRepo)
	planService := services.NewPlanService(planRepo, installmentRepo)
	installmentService := services.NewInstallmentService(installmentRepo, paymentRepo)
	paymentService := services.NewPaymentService(paymentRepo, installmentService)
	pricingService := services.NewPricingService(pricingRepo)
	reportService := services.NewReportService(installmentRepo)

	// 5. Create an AssetServer pointing at ./frontend
	//assetServer, err := assetserver.NewAssetServer(
	//	"frontend",                     // directory to serve
	//	assetopts.Options{},            // default options
	//	false,                          // false = do not embed (dev mode)
	//	assetserver.Logger(nil),        // no custom logger
	//	assetserver.RuntimeAssets(nil), // no runtime assets override
	//)
	if err != nil {
		log.Fatalf("Failed to create AssetServer: %v", err)
	}

	// 6. Build and run the Wails app
	err = wails.Run(&options.App{
		Title:  "Realtor Installment Assistant",
		Width:  1200,
		Height: 800,
		OnStartup: func(ctx context.Context) {
			// (Optional startup logic here)
		},
		Bind: []interface{}{
			propertyService,
			buyerService,
			planService,
			installmentService,
			paymentService,
			pricingService,
			reportService,
		},
		//Assets: assetServer, // use the AssetServer we just created
		// Optional: if you want a custom HTTP handler (e.g., for API endpoints),
		// you can add a Middleware. By default, none is needed:
		//Middleware: []http.Handler{},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
