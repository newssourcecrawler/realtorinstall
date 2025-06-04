package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
)

func main() {
	// 1. Ensure data folder exists
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("mkdir data: %w", err))
	}

	// 2. Initialize repositories
	propRepo, err := apiRepos.NewSQLitePropertyRepo("data/properties.db")
	if err != nil {
		panic(fmt.Errorf("open property repo: %w", err))
	}
	pricingRepo, err := apiRepos.NewSQLiteLocationPricingRepo("data/pricing.db")
	if err != nil {
		panic(fmt.Errorf("open pricing repo: %w", err))
	}
	buyerRepo, err := apiRepos.NewSQLiteBuyerRepo("data/buyers.db")
	if err != nil {
		panic(fmt.Errorf("open buyer repo: %w", err))
	}
	userRepo, err := apiRepos.NewSQLiteUserRepo("data/users.db")
	if err != nil {
		panic(fmt.Errorf("open user repo: %w", err))
	}

	// 3. Construct services
	authSvc := apiServices.NewAuthService(userRepo)
	propSvc := apiServices.NewPropertyService(propRepo, pricingRepo, userRepo)
	buyerSvc := apiServices.NewBuyerService(buyerRepo, userRepo)

	// 4. Build Gin router with CORS and authentication middleware
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(AuthMiddleware("YOUR_JWT_SECRET"))

	// PROPERTY ROUTES

	// GET /properties
	router.GET("/properties", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		props, svcErr := propSvc.ListProperties(c.Request.Context(), tenantID)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, props)
	})

	// POST /properties
	router.POST("/properties", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		var p models.Property
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, svcErr := propSvc.CreateProperty(c.Request.Context(), tenantID, currentUser, p)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /properties/:id
	router.PUT("/properties/:id", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
			return
		}
		var p models.Property
		if bindErr := c.BindJSON(&p); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}
		svcErr := propSvc.UpdateProperty(c.Request.Context(), tenantID, currentUser, id64, p)
		if svcErr != nil {
			if svcErr == repos.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// DELETE /properties/:id
	router.DELETE("/properties/:id", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
			return
		}
		svcErr := propSvc.DeleteProperty(c.Request.Context(), tenantID, currentUser, id64)
		if svcErr != nil {
			if svcErr == repos.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// BUYER ROUTES

	// GET /buyers
	router.GET("/buyers", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		bs, svcErr := buyerSvc.ListBuyers(c.Request.Context(), tenantID)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, bs)
	})

	// POST /buyers
	router.POST("/buyers", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		var b models.Buyer
		if bindErr := c.BindJSON(&b); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}
		id, svcErr := buyerSvc.CreateBuyer(c.Request.Context(), tenantID, currentUser, b)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /buyers/:id
	router.PUT("/buyers/:id", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid buyer ID"})
			return
		}
		var b models.Buyer
		if bindErr := c.BindJSON(&b); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}
		svcErr := buyerSvc.UpdateBuyer(c.Request.Context(), tenantID, currentUser, id64, b)
		if svcErr != nil {
			if svcErr == repos.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// DELETE /buyers/:id
	router.DELETE("/buyers/:id", func(c *gin.Context) {
		tenantID := c.GetString("currentTenant")
		currentUser := c.GetString("currentUser")
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid buyer ID"})
			return
		}
		svcErr := buyerSvc.DeleteBuyer(c.Request.Context(), tenantID, currentUser, id64)
		if svcErr != nil {
			if svcErr == repos.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// 5. Start HTTP server with graceful shutdown
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

	// 6. Graceful shutdown
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

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("currentUser", claims.UserID)
		c.Set("currentTenant", claims.TenantID)
		c.Next()
	}
}
