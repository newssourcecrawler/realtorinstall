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

	salesRepo, _ := apiRepos.NewSQLiteSalesRepo("data/sales.db")
	introRepo, _ := apiRepos.NewSQLiteIntroductionsRepo("data/introductions.db")
	lettingsRepo, _ := apiRepos.NewSQLiteLettingsRepo("data/lettings.db")
	commissionRepo, _ := apiRepos.NewSQLiteCommissionRepo("data/commissions.db")

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

	reportSvc := apiServices.NewReportService(
		commissionRepo,
		salesRepo,
		pricingRepo,
		// add other repos if needed
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
	router.Use(AuthMiddleware(authSvc, userRepo))

	// 6. Authentication routes
	router.POST("/login", authH.Login)
	router.POST("/register",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		authH.Register,
	)

	// 7. User CRUD routes
	router.GET("/users", userH.List)
	router.POST("/users", userH.Create)
	router.PUT("/users/:id", userH.Update)
	router.DELETE("/users/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		userH.Delete,
	)

	// 8. Property routes
	router.GET("/properties", propH.List)
	router.POST("/properties",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		propH.Create,
	)
	router.PUT("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		propH.Update,
	)
	router.DELETE("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		propH.Delete,
	)

	// 9. Buyer routes
	router.GET("/buyers", buyerH.List)
	router.POST("/buyers",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		buyerH.Create,
	)
	router.PUT("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		buyerH.Update,
	)
	router.DELETE("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		buyerH.Delete,
	)

	// 10. Pricing routes
	router.GET("/pricing", priceH.List)
	router.POST("/pricing",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		priceH.Create,
	)
	router.PUT("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		priceH.Update,
	)
	router.DELETE("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		priceH.Delete,
	)

	// 11. Sales routes
	router.GET("/sales", salesH.List)
	router.POST("/sales",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		salesH.Create,
	)
	router.PUT("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		salesH.Update,
	)
	router.DELETE("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		salesH.Delete,
	)

	// 12. Introduction routes
	router.GET("/introductions", introH.List)
	router.POST("/introductions",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		introH.Create,
	)
	router.PUT("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		introH.Update,
	)
	router.DELETE("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		introH.Delete,
	)

	// 13. Lettings routes
	router.GET("/lettings", lettingsH.List)
	router.POST("/lettings",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		lettingsH.Create,
	)
	router.PUT("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		lettingsH.Update,
	)
	router.DELETE("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		lettingsH.Delete,
	)

	// 14. Plan routes
	router.GET("/plans", planH.List)
	router.POST("/plans",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		planH.Create,
	)
	router.PUT("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		planH.Update,
	)
	router.DELETE("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		planH.Delete,
	)

	// 15. Installment routes
	router.GET("/installments", instH.List)
	router.GET("/installments/plan/:planId", instH.ListByPlan)
	router.POST("/installments",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		instH.Create,
	)
	router.PUT("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		instH.Update,
	)
	router.DELETE("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		instH.Delete,
	)

	// 16. Payment routes
	router.GET("/payments", payH.List)
	router.POST("/payments",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		payH.Create,
	)
	router.PUT("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		payH.Update,
	)
	router.DELETE("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		payH.Delete,
	)

	// 17. Commission routes
	router.GET("/commissions", commissionH.List)
	router.POST("/commissions",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		commissionH.Create,
	)
	router.PUT("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		commissionH.Update,
	)
	router.DELETE("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		AuthorizeRoles(userRepo, "admin", "manager"),
		commissionH.Delete,
	)

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
		c.Next()
	}
}

// AuthorizeRoles checks that the logged‐in user's role is one of the allowed list.
func AuthorizeRoles(userRepo apiRepos.UserRepo, allowed ...string) gin.HandlerFunc {
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
