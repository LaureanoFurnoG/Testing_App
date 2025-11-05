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
	group.POST("/inviteGroup", mw.RequireAuth(), handler.inviteGroup)
	group.PATCH("/acceptInvitation", mw.RequireAuth(), handler.acceptInvitation)
	group.DELETE("/declineGroup", mw.RequireAuth(), handler.declineGroup)

	group.DELETE("/deleteGroup", mw.RequireAuth(), handler.deleteGroup)

	group.GET("/showAllGroups", mw.RequireAuth(), handler.showAllGroups)

}

func (h *HandlerAPI) createGroup(c *gin.Context) {
	var jsonData struct {
		Name string
	}
	accessToken := c.GetHeader("Access-Token")

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

	groupDelete := initializers.DB.Delete(&group, jsonData.GroupID)

	if groupDelete.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}
	groupRelationDelete := initializers.DB.Delete(&GroupsRelation, GroupsRelation.ID)

	if groupRelationDelete.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group relation Missing",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group deleted successfully",
	})
}

func (h *HandlerAPI) inviteGroup(c *gin.Context) {
	var jsonData struct {
		Email   string
		GroupID string
	}

	//accessToken := c.Request.Header.Get("Access-Token") //temporal
	//accessToken := c.GetHeader("Authorization")
	//tokenString := strings.TrimPrefix(accessToken, "Bearer ")

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var group models.Groups
	searchGroup := initializers.DB.Find(&group, "keycloak_id = ?", jsonData.GroupID)

	if searchGroup.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	var userF models.Users
	userKeycloak, err := h.clientKC.GetUserInf(c.Request.Context(), jsonData.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	userFind := initializers.DB.Find(&userF, "keycloak_id = ?", userKeycloak.ID)

	if userFind.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Missing",
		})
		return
	}

	groupRelation := models.GroupsRelation{Idgroup: group.ID, Iduser: userF.ID, Accepted: false}

	createRelation := initializers.DB.Create(&groupRelation)

	if createRelation.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Relation not created",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group invitation sended",
	})
}

func (h *HandlerAPI) acceptInvitation(c *gin.Context) {
	var jsonData struct {
		GroupID int
	}
	accessToken := c.Request.Header.Get("Access-Token") //temporal

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//search user
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
	//group
	var group models.Groups
	var GroupsRelation models.GroupsRelation

	groupFound := initializers.DB.First(&group, "id = ?", jsonData.GroupID)

	if groupFound.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	GroupRelationChange := initializers.DB.Model(&GroupsRelation).
		Where("idgroup = ? AND iduser = ?", jsonData.GroupID, userF.ID).
		Update("accepted", true)

	if GroupRelationChange.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group invitation accepted successfully",
	})
}

func (h *HandlerAPI) declineGroup(c *gin.Context) {
	var jsonData struct {
		GroupID int
	}
	accessToken := c.Request.Header.Get("Access-Token") //temporal

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//search user
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
	//group
	var group models.Groups
	var GroupsRelation models.GroupsRelation

	groupFound := initializers.DB.First(&group, "id = ?", jsonData.GroupID)

	if groupFound.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	GroupRelationChange := initializers.DB.Model(&GroupsRelation).
		Where("idgroup = ? AND iduser = ?", jsonData.GroupID, userF.ID).Delete(&GroupsRelation)

	if GroupRelationChange.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group invitation declined successfully",
	})
}

func (h *HandlerAPI) showAllGroups(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token") //temporal

	userKeycloak, err := h.clientKC.UserInfo(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	userKeycloakGroups, err := h.clientKC.GetGroups(c.Request.Context(), accessToken, userKeycloak.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Groups": userKeycloakGroups,
	})
}
