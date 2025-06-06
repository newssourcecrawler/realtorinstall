package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/services"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) TotalCommissionByBeneficiary(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.TotalCommissionByBeneficiary(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) OutstandingInstallmentsByPlan(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.OutstandingInstallmentsByPlan(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) MonthlySalesVolume(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.MonthlySalesVolume(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) ActiveLettingsRentRoll(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.ActiveLettingsRentRoll(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) TopPropertiesByPaymentVolume(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.TopPropertiesByPaymentVolume(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
