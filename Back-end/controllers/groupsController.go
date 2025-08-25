package controllers

import (
	keycloak "Go-API-T/Keycloak"
	"Go-API-T/middlewere"
	//"Go-API-T/services"
	//"crypto/rand"
	//"fmt"
	//"math/big"

	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v5"
	//"golang.org/x/crypto/bcrypt"
	//"strconv"
)

func GroupsController(rg *gin.RouterGroup, handler *HandlerAPI, mw *middlewere.Middleware) {
	group := rg.Group("/group")

	group.POST("/createGroup", mw.RequireAuth(), handler.createGroup)
}

func (h *HandlerAPI) createGroup(c *gin.Context) {
	var jsonData struct {
		Name string
	}

	accessToken := c.Request.Header.Get("Access-Token")

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	err := h.clientKC.CreateGroup(c.Request.Context(), keycloak.CreateGroupParams{Name: jsonData.Name}, accessToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Group created successfully",
	})
}

func (h *HandlerAPI) deleteGroup(c *gin.Context){
	
}