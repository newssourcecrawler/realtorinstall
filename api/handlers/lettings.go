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

type LettingsHandler struct {
	svc *services.LettingsService
}

func NewLettingsHandler(svc *services.LettingsService) *LettingsHandler {
	return &LettingsHandler{svc: svc}
}

func (h *LettingsHandler) List(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	list, err := h.svc.ListLettings(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *LettingsHandler) Create(c *gin.Context) {
	var lt models.Lettings
	if err := c.BindJSON(&lt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	id, err := h.svc.CreateLetting(context.Background(), tenantID, currentUser, lt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *LettingsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid letting ID"})
		return
	}
	var lt models.Lettings
	if err := c.BindJSON(&lt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.UpdateLetting(context.Background(), tenantID, currentUser, id64, lt); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "letting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *LettingsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid letting ID"})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.DeleteLetting(context.Background(), tenantID, currentUser, id64); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "letting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
