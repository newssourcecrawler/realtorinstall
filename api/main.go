package main

import (
	"context"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
	"github.com/newssourcecrawler/realtorinstall/internal/models"
	intRepos "github.com/newssourcecrawler/realtorinstall/internal/repos"
)

func main() {
	// 1. Ensure data folder exists (for SQLite files)
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	// 2. Initialize the API's PropertyRepo (api/repos/sqlite_property_repo.go)
	propRepo, err := repos.NewSQLitePropertyRepo("data/properties.db")
	if err != nil {
		panic(err)
	}

	// 3. Initialize the internal LocationPricingRepo (internal/repos/sqlite_locationpricing_repo.go)
	pricingRepo, err := intRepos.NewSQLiteLocationPricingRepo("data/pricing.db")
	if err != nil {
		panic(err)
	}

	// 4. Construct the PropertyService with both repos
	propSvc := services.NewPropertyService(propRepo, pricingRepo)

	// 5. Create Gin router and enable CORS
	r := gin.Default()
	r.Use(cors.Default())

	// 6. GET /properties → list all properties
	r.GET("/properties", func(c *gin.Context) {
		props, err := propSvc.ListProperties(context.Background())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, props)
	})

	// 7. POST /properties → create a new property
	r.POST("/properties", func(c *gin.Context) {
		var p models.Property
		if err := c.BindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		id, err := propSvc.CreateProperty(context.Background(), p)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"id": id})
	})

	// 8. Start the server on port 8080
	r.Run(":8080")
}
