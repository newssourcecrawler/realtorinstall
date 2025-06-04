package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
)

type PricingHandler struct {
	svc *services.PricingService
}

func NewPricingHandler(svc *services.PricingService) *PricingHandler {
	return &PricingHandler{svc: svc}
}

func (h *PricingHandler) List(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	list, err := h.svc.ListLocationPricings(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *PricingHandler) Create(c *gin.Context) {
	var lp models.LocationPricing
	if err := c.BindJSON(&lp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	id, err := h.svc.CreateLocationPricing(context.Background(), tenantID, currentUser, lp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *PricingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pricing ID"})
		return
	}
	var lp models.LocationPricing
	if err := c.BindJSON(&lp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.UpdateLocationPricing(context.Background(), tenantID, currentUser, id64, lp); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pricing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *PricingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pricing ID"})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.DeleteLocationPricing(context.Background(), tenantID, currentUser, id64); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pricing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
