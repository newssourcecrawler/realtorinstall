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

func (h *ReportHandler) CommissionsByBeneficiary(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.CommissionByBeneficiary(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
