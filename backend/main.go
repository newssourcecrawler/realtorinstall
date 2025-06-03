// backend/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// 1) Import the Wails core package so CreateApp is available:
	"github.com/wailsapp/wails/v2"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	// Our internal packages:
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
	"github.com/newssourcecrawler/realtorinstall/internal/services"
	"github.com/newssourcecrawler/realtorinstall/internal/utils"
)

func main() {
	// 1. Load config file
	cfg, err := utils.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Ensure the data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data folder: %v", err)
	}

	// 3. Initialize all SQLite repos (no encryption for MVP)
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

	// 4. Create service layer instances
	propertyService := services.NewPropertyService(propertyRepo, pricingRepo)
	buyerService := services.NewBuyerService(buyerRepo)
	planService := services.NewPlanService(planRepo, installmentRepo)
	installmentService := services.NewInstallmentService(installmentRepo, paymentRepo)
	paymentService := services.NewPaymentService(paymentRepo, installmentService)
	pricingService := services.NewPricingService(pricingRepo)
	reportService := services.NewReportService(installmentRepo)

	// 5. Build and run the Wails app, exporting our services to JS
	app := wails.CreateApp(&options.App{
		Title:  "Realtor Installment Assistant",
		Width:  1200,
		Height: 800,
		OnStartup: func(ctx context.Context) {
			// any startup logic goes here
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
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				TitleVisibility:            mac.TitleVisibilityHidden,
			},
		},
		Linux: &linux.Options{
			Icon: "frontend/assets/icon.png",
		},
		Assets: nil, // let Wails bundle the frontend/build folder
	})

	err = app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
