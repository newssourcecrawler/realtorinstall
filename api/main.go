package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
)

func main() {
	// 1. Ensure data folder exists (for SQLite)
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("mkdir data: %w", err))
	}

	// 2. Initialize repos
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

	// 3. Construct services
	propSvc := apiServices.NewPropertyService(propRepo, pricingRepo)
	buyerSvc := apiServices.NewBuyerService(buyerRepo)

	// 4. Build Gin router
	router := gin.Default()
	router.Use(cors.Default())

	// ====== PROPERTY ROUTES ======

	// GET /properties
	router.GET("/properties", func(c *gin.Context) {
		props, svcErr := propSvc.ListProperties(context.Background())
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, props)
	})

	// POST /properties
	router.POST("/properties", func(c *gin.Context) {
		var p models.Property
		if bindErr := c.BindJSON(&p); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}
		// Set CreatedAt / LastModified here, if not already set:
		if p.CreatedAt.IsZero() {
			p.CreatedAt = time.Now().UTC()
		}
		if p.LastModified.IsZero() {
			p.LastModified = time.Now().UTC()
		}
		p.Deleted = false

		id, svcErr := propSvc.CreateProperty(context.Background(), p)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /properties/:id  → update an existing property
	router.PUT("/properties/:id", func(c *gin.Context) {
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
		// Overwrite any JSON‐sent ID so service uses path param:
		p.ID = id64
		p.LastModified = time.Now().UTC()

		svcErr := propSvc.UpdateProperty(context.Background(), id64, p)
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

	// DELETE /properties/:id  → soft-delete a property
	router.DELETE("/properties/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
			return
		}
		svcErr := propSvc.DeleteProperty(context.Background(), id64)
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

	// ====== BUYER ROUTES ======

	// GET /buyers
	router.GET("/buyers", func(c *gin.Context) {
		bs, svcErr := buyerSvc.ListBuyers(context.Background())
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, bs)
	})

	// POST /buyers
	router.POST("/buyers", func(c *gin.Context) {
		var b models.Buyer
		if bindErr := c.BindJSON(&b); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}
		// Stamp timestamps:
		if b.CreatedAt.IsZero() {
			b.CreatedAt = time.Now().UTC()
		}
		if b.LastModified.IsZero() {
			b.LastModified = time.Now().UTC()
		}
		b.Deleted = false

		id, svcErr := buyerSvc.CreateBuyer(context.Background(), b)
		if svcErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /buyers/:id  → update a buyer
	router.PUT("/buyers/:id", func(c *gin.Context) {
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
		b.ID = id64
		b.LastModified = time.Now().UTC()

		svcErr := buyerSvc.UpdateBuyer(context.Background(), id64, b)
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

	// DELETE /buyers/:id  → soft-delete a buyer
	router.DELETE("/buyers/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id64, parseErr := strconv.ParseInt(idStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid buyer ID"})
			return
		}
		svcErr := buyerSvc.DeleteBuyer(context.Background(), id64)
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

	// 6. Graceful shutdown on Ctrl+C
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
