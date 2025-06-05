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

type IntroductionsHandler struct {
	svc *services.IntroductionsService
}

func NewIntroductionsHandler(svc *services.IntroductionsService) *IntroductionsHandler {
	return &IntroductionsHandler{svc: svc}
}

func (h *IntroductionsHandler) List(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	list, err := h.svc.ListIntroductions(context.Background(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *IntroductionsHandler) Create(c *gin.Context) {
	var intro models.Introductions
	if err := c.BindJSON(&intro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	id, err := h.svc.CreateIntroduction(context.Background(), tenantID, currentUser, intro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *IntroductionsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid introduction ID"})
		return
	}
	var intro models.Introductions
	if err := c.BindJSON(&intro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.UpdateIntroduction(context.Background(), tenantID, currentUser, id64, intro); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "introduction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *IntroductionsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid introduction ID"})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	if err := h.svc.DeleteIntroduction(context.Background(), tenantID, currentUser, id64); err != nil {
		if err == repos.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "introduction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
