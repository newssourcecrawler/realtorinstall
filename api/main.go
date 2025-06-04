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

	// 2. Initialize repositories
	propRepo, _ := apiRepos.NewSQLitePropertyRepo("data/properties.db")
	pricingRepo, _ := apiRepos.NewSQLiteLocationPricingRepo("data/pricing.db")
	buyerRepo, _ := apiRepos.NewSQLiteBuyerRepo("data/buyers.db")
	userRepo, _ := apiRepos.NewSQLiteUserRepo("data/users.db")
	planRepo, _ := apiRepos.NewSQLiteInstallmentPlanRepo("data/plans.db")
	instRepo, _ := apiRepos.NewSQLiteInstallmentRepo("data/installments.db")
	payRepo, _ := apiRepos.NewSQLitePaymentRepo("data/payments.db")

	// 3. Construct services
	authSvc := apiServices.NewAuthService(userRepo)
	propSvc := apiServices.NewPropertyService(propRepo, pricingRepo, userRepo)
	buyerSvc := apiServices.NewBuyerService(buyerRepo, userRepo)
	pricingSvc := apiServices.NewPricingService(pricingRepo)
	planSvc := apiServices.NewPlanService(planRepo, instRepo)
	instSvc := apiServices.NewInstallmentService(instRepo, payRepo)
	paySvc := apiServices.NewPaymentService(payRepo, instSvc)

	// 4. Instantiate handlers
	authH := handlers.NewAuthHandler(authSvc)
	propH := handlers.NewPropertyHandler(propSvc)
	buyerH := handlers.NewBuyerHandler(buyerSvc)
	priceH := handlers.NewPricingHandler(pricingSvc)
	planH := handlers.NewPlanHandler(planSvc)
	instH := handlers.NewInstallmentHandler(instSvc)
	payH := handlers.NewPaymentHandler(paySvc)

	// 5. Build Gin router with CORS + JWT middleware
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(AuthMiddleware("YOUR_JWT_SECRET"))

	// 6. Register Authentication routes
	router.POST("/login", authH.Login)
	router.POST("/register", authH.Register)

	// 7. Register Property routes
	router.GET("/properties", propH.List)
	router.POST("/properties", propH.Create)
	router.PUT("/properties/:id", propH.Update)
	router.DELETE("/properties/:id", propH.Delete)

	// 8. Register Buyer routes
	router.GET("/buyers", buyerH.List)
	router.POST("/buyers", buyerH.Create)
	router.PUT("/buyers/:id", buyerH.Update)
	router.DELETE("/buyers/:id", buyerH.Delete)

	// 9. Register Pricing routes
	router.GET("/pricing", priceH.List)
	router.POST("/pricing", priceH.Create)
	router.PUT("/pricing/:id", priceH.Update)
	router.DELETE("/pricing/:id", priceH.Delete)

	// 10. Register Plan routes
	router.GET("/plans", planH.List)
	router.POST("/plans", planH.Create)
	router.PUT("/plans/:id", planH.Update)
	router.DELETE("/plans/:id", planH.Delete)

	// 11. Register Installment routes
	router.GET("/installments", instH.List)
	router.GET("/installments/plan/:planId", instH.ListByPlan)
	router.POST("/installments", instH.Create)
	router.PUT("/installments/:id", instH.Update)
	router.DELETE("/installments/:id", instH.Delete)

	// 12. Register Payment routes
	router.GET("/payments", payH.List)
	router.POST("/payments", payH.Create)
	router.PUT("/payments/:id", payH.Update)
	router.DELETE("/payments/:id", payH.Delete)

	// 13. Start HTTP server with graceful shutdown
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

	// 14. Graceful shutdown on SIGINT
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
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := apiServices.ParseToken(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("currentUser", claims.UserID)
		c.Set("currentTenant", claims.TenantID)
		c.Next()
	}
}
