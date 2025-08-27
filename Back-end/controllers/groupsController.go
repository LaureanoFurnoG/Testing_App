package controllers

import (
	keycloak "Go-API-T/Keycloak"
	"Go-API-T/initializers"
	"Go-API-T/middlewere"
	"Go-API-T/models"
	//"fmt"
	//"strings"

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
	group.DELETE("/deleteGroup", mw.RequireAuth(), handler.deleteGroup)

}

func (h *HandlerAPI) createGroup(c *gin.Context) {
	var jsonData struct {
		Name string
	}
	accessToken := c.Request.Header.Get("Access-Token") //temporal
	//accessToken := c.GetHeader("Authorization")
	//tokenString := strings.TrimPrefix(accessToken, "Bearer ")

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	groupID, err := h.clientKC.CreateGroup(c.Request.Context(), keycloak.CreateGroupParams{Name: jsonData.Name}, accessToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	group := models.Groups{KeycloakID: groupID}

	createG := initializers.DB.Create(&group)

	if createG.Error != nil {
		c.JSON(400, gin.H{"error": createG.Error})
		return
	}
	var userF models.Users
	userKeycloak, err := h.clientKC.UserInfo(c.Request.Context(), accessToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	userFind := initializers.DB.First(&userF, "keycloak_id = ?", userKeycloak.ID)

	if userFind.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Missing",
		})
		return
	}

	groupRelation := models.GroupsRelation{Idgroup: group.ID, Iduser: userF.ID, Accepted: true}

	createRelation := initializers.DB.Create(&groupRelation)

	if createRelation.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Missing",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group created successfully",
	})
}

func (h *HandlerAPI) deleteGroup(c *gin.Context) {
	var jsonData struct {
		GroupID int
	}

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var group models.Groups
	var GroupsRelation models.GroupsRelation

	groupFound := initializers.DB.First(&group, "id = ?", jsonData.GroupID)

	if groupFound.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	groupRleatioNFound := initializers.DB.First(&GroupsRelation, "idgroup = ?", jsonData.GroupID)

	if groupRleatioNFound.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group relation Missing",
		})
		return
	}

	err := h.clientKC.DeleteGroup(c.Request.Context(), group.KeycloakID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	groupRelationDelete := initializers.DB.Delete(&group, jsonData.GroupID)

	if groupRelationDelete.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}
	groupDelete := initializers.DB.Delete(&GroupsRelation, GroupsRelation.ID)

	if groupDelete.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group deleted successfully",
	})
}

func (h *HandlerAPI) inviteGroup(c *gin.Context) {

}
