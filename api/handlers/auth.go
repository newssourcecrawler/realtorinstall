package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"github.com/newssourcecrawler/realtorinstall/api/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register expects JSON:
//
//	{
//	  "user": { "username": "...", "first_name": "...", "last_name": "...", "role": "...", ... },
//	  "password": "rawPasswordValue"
//	}
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		User     models.User `json:"user"`
		Password string      `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	currentUser := c.GetString("currentUser")
	id, err := h.svc.Register(context.Background(), tenantID, currentUser, req.User, req.Password)
	if err != nil {
		// If username taken or other errors, return 400 or 500 accordingly
		switch err {
		case repos.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case repos.ErrInvalidRegistration:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Login expects JSON:
//
//	{
//	  "username": "someName",
//	  "password": "someRawPassword"
//	}
//
// On success, returns { "token": "<JWT string>" }.
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("currentTenant")
	token, err := h.svc.Login(context.Background(), tenantID, req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
