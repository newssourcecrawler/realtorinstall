// api/handlers/property.go
package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/services"
	intModels "github.com/newssourcecrawler/realtorinstall/internal/models"
)

type PropertyHandler struct {
	svc *services.PropertyService
}

func NewPropertyHandler(svc *services.PropertyService) *PropertyHandler {
	return &PropertyHandler{svc: svc}
}

func (h *PropertyHandler) List(c *gin.Context) {
	ps, err := h.svc.ListProperties(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ps)
}

func (h *PropertyHandler) Create(c *gin.Context) {
	var p intModels.Property
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.svc.CreateProperty(context.Background(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
