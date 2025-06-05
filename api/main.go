package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/handlers"
	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
)

func main() {
	// 1. Ensure data folder exists
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("mkdir data: %w", err))
	}

	// 2. Initialize repositories (one per domain)
	propRepo, _ := apiRepos.NewSQLitePropertyRepo("data/properties.db")
	pricingRepo, _ := apiRepos.NewSQLiteLocationPricingRepo("data/pricing.db")
	buyerRepo, _ := apiRepos.NewSQLiteBuyerRepo("data/buyers.db")
	userRepo, _ := apiRepos.NewSQLiteUserRepo("data/users.db")
	planRepo, _ := apiRepos.NewSQLiteInstallmentPlanRepo("data/plans.db")
	instRepo, _ := apiRepos.NewSQLiteInstallmentRepo("data/installments.db")
	payRepo, _ := apiRepos.NewSQLitePaymentRepo("data/payments.db")

	// New repositories for additional domains:
	salesRepo, _ := apiRepos.NewSQLiteSalesRepo("data/sales.db")
	introRepo, _ := apiRepos.NewSQLiteIntroductionsRepo("data/introductions.db")
	lettingsRepo, _ := apiRepos.NewSQLiteLettingsRepo("data/lettings.db")
	commissionRepo, _ := apiRepos.NewSQLiteCommissionRepo("data/commissions.db")

	// 3. Construct services
	authSvc := apiServices.NewAuthService(userRepo, "YOUR_VERY_LONG_SECRET_KEY", time.Hour*24)
	propSvc := apiServices.NewPropertyService(propRepo, userRepo, pricingRepo)
	//buyerSvc := apiServices.NewBuyerService(buyerRepo, userRepo)
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

	// Assume ReportService aggregates from multiple repos:
	reportSvc := apiServices.NewReportService(
		commissionRepo,
		salesRepo,
		pricingRepo,
		// add any other repos needed for reporting
	)

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
	router.Use(AuthMiddleware("YOUR_JWT_SECRET"))

	// 6. Authentication routes
	router.POST("/login", authH.Login)
	//router.POST("/register", authH.Register)
	router.POST("/register", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), authH.Register)

	// 7. User CRUD routes
	router.GET("/users", userH.List)
	router.POST("/users", userH.Create)
	router.PUT("/users/:id", userH.Update)
	//router.DELETE("/users/:id", userH.Delete)
	router.DELETE("/users/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), usersH.Delete)

	// 8. Property routes
	router.GET("/properties", propH.List)
	router.POST("/properties", propH.Create)
	router.PUT("/properties/:id", propH.Update)
	router.DELETE("/properties/:id", propH.Delete)

	// 9. Buyer routes
	router.GET("/buyers", buyerH.List)
	router.POST("/buyers", buyerH.Create)
	router.PUT("/buyers/:id", buyerH.Update)
	router.DELETE("/buyers/:id", buyerH.Delete)

	// 10. Pricing routes
	router.GET("/pricing", priceH.List)
	router.POST("/pricing", priceH.Create)
	router.PUT("/pricing/:id", priceH.Update)
	router.DELETE("/pricing/:id", priceH.Delete)

	// 11. Sales routes
	router.GET("/sales", salesH.List)
	router.POST("/sales", salesH.Create)
	router.PUT("/sales/:id", salesH.Update)
	//router.DELETE("/sales/:id", salesH.Delete)
	router.DELETE("/sales/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), salesH.Delete)

	// 12. Introduction routes
	router.GET("/introductions", introH.List)
	router.POST("/introductions", introH.Create)
	router.PUT("/introductions/:id", introH.Update)
	//router.DELETE("/introductions/:id", introH.Delete)
	router.DELETE("/introductions/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), introH.Delete)

	// 13. Lettings routes
	router.GET("/lettings", lettingsH.List)
	router.POST("/lettings", lettingsH.Create)
	router.PUT("/lettings/:id", lettingsH.Update)
	//router.DELETE("/lettings/:id", lettingsH.Delete)
	router.DELETE("/lettings/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), lettingsH.Delete)

	// 14. Plan routes
	router.GET("/plans", planH.List)
	router.POST("/plans", planH.Create)
	router.PUT("/plans/:id", planH.Update)
	//router.DELETE("/plans/:id", planH.Delete)
	router.DELETE("/plans/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), planH.Delete)

	// 15. Installment routes
	router.GET("/installments", instH.List)
	router.GET("/installments/plan/:planId", instH.ListByPlan)
	router.POST("/installments", instH.Create)
	router.PUT("/installments/:id", instH.Update)
	//router.DELETE("/installments/:id", instH.Delete)
	router.DELETE("/installments/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), instH.Delete)

	// 16. Payment routes
	router.GET("/payments", payH.List)
	router.POST("/payments", payH.Create)
	router.PUT("/payments/:id", payH.Update)
	//router.DELETE("/payments/:id", payH.Delete)
	router.DELETE("/payments/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), payH.Delete)

	// 17. Commission routes
	router.GET("/commissions", commissionH.List)
	router.POST("/commissions", commissionH.Create)
	router.PUT("/commissions/:id", commissionH.Update)
	//router.DELETE("/commissions/:id", commissionH.Delete)
	router.DELETE("/commissions/:id", AuthMiddleware(authSvc), AuthorizeRoles("admin", "manager"), commissionsH.Delete)

	// 18. Reporting routes
	router.GET("/reports/commissions/beneficiary", reportH.CommissionsByBeneficiary)

	// 19. Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		fmt.Println("API listening on http://localhost:8080")
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
func AuthMiddleware(authSvc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if !strings.HasPrefix(hdr, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(hdr, "Bearer ")
		claims, err := authSvc.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		if claims.LicenseExp < time.Now().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"error": "license expired"})
			return
		}
		c.Set("currentUser", claims.UserID)
		c.Set("currentTenant", claims.TenantID)
		c.Next()
	}
}
