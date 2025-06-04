package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	apiRepos "github.com/newssourcecrawler/realtorinstall/api/repos"
	apiServices "github.com/newssourcecrawler/realtorinstall/api/services"
	"github.com/newssourcecrawler/realtorinstall/internal/models"
	intRepos "github.com/newssourcecrawler/realtorinstall/internal/repos"
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
	pricingRepo, err := intRepos.NewSQLiteLocationPricingRepo("data/pricing.db")
	if err != nil {
		panic(fmt.Errorf("open pricing repo: %w", err))
	}

	// 3. Construct services
	propSvc := apiServices.NewPropertyService(propRepo, pricingRepo)

	// Buyer service (you must have a BuyerRepo under api/repos and service under api/services)
	buyerRepo, err := apiRepos.NewSQLiteBuyerRepo("data/buyers.db")
	if err != nil {
		panic(fmt.Errorf("open buyer repo: %w", err))
	}
	buyerSvc := apiServices.NewBuyerService(buyerRepo)

	// 4. Build Gin router
	router := gin.Default()
	router.Use(cors.Default())

	// ====== PROPERTY ROUTES ======

	// GET /properties
	router.GET("/properties", func(c *gin.Context) {
		props, err := propSvc.ListProperties(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, props)
	})

	// POST /properties
	router.POST("/properties", func(c *gin.Context) {
		var p models.Property
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := propSvc.CreateProperty(context.Background(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /properties/:id  → update an existing property
	router.PUT("/properties/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var p models.Property
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Ensure the ID field in the body matches the path param (optional):
		p.ID = 0 // ignore any ID in JSON; let service use idParam
		err := propSvc.UpdateProperty(context.Background(), idParam, p)
		if err != nil {
			if err == apiServices.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// DELETE /properties/:id  → delete a property
	router.DELETE("/properties/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		err := propSvc.DeleteProperty(context.Background(), idParam)
		if err != nil {
			if err == apiServices.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// ====== BUYER ROUTES ======

	// GET /buyers
	router.GET("/buyers", func(c *gin.Context) {
		bs, err := buyerSvc.ListBuyers(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, bs)
	})

	// POST /buyers
	router.POST("/buyers", func(c *gin.Context) {
		var b models.Buyer
		if err := c.BindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := buyerSvc.CreateBuyer(context.Background(), b)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// PUT /buyers/:id  → update a buyer
	router.PUT("/buyers/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var b models.Buyer
		if err := c.BindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		b.ID = 0 // ignore any ID from JSON; service will use idParam
		err := buyerSvc.UpdateBuyer(context.Background(), idParam, b)
		if err != nil {
			if err == apiServices.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// DELETE /buyers/:id  → delete a buyer
	router.DELETE("/buyers/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		err := buyerSvc.DeleteBuyer(context.Background(), idParam)
		if err != nil {
			if err == apiServices.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// 5. Create http.Server with Gin as handler
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 6. Run server in a goroutine
	go func() {
		fmt.Println("API listening on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Listen error: %v\n", err)
			os.Exit(1)
		}
	}()

	// 7. Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down API server...")

	// 8. Attempt graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server forced to shutdown: %v\n", err)
	}

	fmt.Println("API server stopped.")
}
