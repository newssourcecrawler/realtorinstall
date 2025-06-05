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

type CommissionHandler struct {
	svc *services.CommissionService
}

func NewCommissionHandler(svc *services.CommissionService) *CommissionHandler {
	return &CommissionHandler{svc: svc}
}

func (h *CommissionHandler) List(c *gin.Context) {
	tenantID := c.GetString("currentTenant")

	filterType := c.Query("transaction_type")

	beneficiaryIDStr := c.Query("beneficiary_id")
	var beneficiaryID int64
	if beneficiaryIDStr != "" {
		if id, err := strconv.ParseInt(beneficiaryIDStr, 10, 64); err == nil {
			beneficiaryID = id
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary_id"})
			return
		}
	}

	list, err := h.svc.ListCommissions(context.Background(), tenantID, filterType, beneficiaryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CommissionHandler) Create(c *gin.Context) {
	var cm models.Commission
	if err := c.BindJSON(&cm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	id, err := h.svc.CreateCommission(context.Background(), tenantID, currentUser, cm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *CommissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commission ID"})
		return
	}
	var cm models.Commission
	if err := c.BindJSON(&cm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.UpdateCommission(context.Background(), tenantID, currentUser, id64, cm); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "commission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *CommissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commission ID"})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.DeleteCommission(context.Background(), tenantID, currentUser, id64); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "commission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
