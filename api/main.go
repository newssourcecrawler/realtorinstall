package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

func main() {
	// 1. Ensure data folder exists (SQLite)
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	// 2. Initialize SQLite-backed repo + service
	dbPath := "data/properties.db"
	propRepo, err := repos.NewSQLitePropertyRepo(dbPath)
	if err != nil {
		panic(err)
	}
	propSvc := services.NewPropertyService(propRepo)

	// 3. Create Gin router
	r := gin.Default()
	// Allow CORS from any origin (for local dev):
	r.Use(cors.Default())

	// 4. Register property‐related routes
	//    GET  /properties         → list all
	//    POST /properties         → create new
	//    (You can add PUT/DELETE later)
	r.GET("/properties", func(c *gin.Context) {
		ps, err := propSvc.ListProperties(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ps)
	})
	r.POST("/properties", func(c *gin.Context) {
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

	// 5. (Optionally add handlers for installments, buyers, etc.)

	// 6. Start server on port 8080
	r.Run(":8080")
}
