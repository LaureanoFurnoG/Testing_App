package controllers

import (
	keycloak "Go-API-T/Keycloak"
	"Go-API-T/initializers"
	"Go-API-T/middlewere"
	"Go-API-T/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GroupsController(rg *gin.RouterGroup, handler *HandlerAPI, mw *middlewere.Middleware) {
	group := rg.Group("/group")

	group.POST("/createGroup", mw.RequireAuth(), handler.createGroup)
	group.POST("/inviteGroup/:groupId", mw.RequireAuth(), mw.BelongsGroup(), handler.inviteGroup)
	group.PATCH("/acceptInvitation", mw.RequireAuth(), handler.acceptInvitation)
	group.DELETE("/declineGroup", mw.RequireAuth(), handler.declineGroup)

	group.DELETE("/deleteGroup/:groupId", mw.RequireAuth(), mw.BelongsGroup(), handler.deleteGroup)

	group.GET("/showAllGroups", mw.RequireAuth(), handler.showAllGroups)

}

func (h *HandlerAPI) createGroup(c *gin.Context) {
	var jsonData struct {
		Name string
	}

	accessHeader := c.GetHeader("Authorization")
	if accessHeader == "" || !strings.HasPrefix(accessHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access token not found or invalid format"})
		c.Abort()
		return
	}
	if jsonData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing the group name",
		})
		return
	}

	accessToken := strings.TrimPrefix(accessHeader, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)

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
		GroupID int
	}

	convInt, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse groupId",
		})
		return
	}
	jsonData.GroupID = convInt

	//accessToken  := c.GetHeader("Authorization")
	//accessToken := c.GetHeader("Authorization")
	//tokenString := strings.TrimPrefix(accessToken, "Bearer ")

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var group models.Groups
	searchGroup := initializers.DB.Find(&group, "id = ?", jsonData.GroupID)

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

	accessHeader := c.GetHeader("Authorization")
	if accessHeader == "" || !strings.HasPrefix(accessHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access token not found or invalid format"})
		c.Abort()
		return
	}

	accessToken := strings.TrimPrefix(accessHeader, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)

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

	accessHeader := c.GetHeader("Authorization")
	if accessHeader == "" || !strings.HasPrefix(accessHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access token not found or invalid format"})
		c.Abort()
		return
	}

	accessToken := strings.TrimPrefix(accessHeader, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)

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
		Where("idgroup = ? AND iduser = ? AND accepted = ?", jsonData.GroupID, userF.ID, false).Delete(&GroupsRelation)

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

type GroupUser struct {
	IDKeycloak string                 `json:"id_keycloak"`
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	SubGroups  map[string]interface{} `json:"sub_groups"`
}

func (h *HandlerAPI) showAllGroups(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authenticated",
		})
		return
	}

	userKeycloakGroups, err := h.clientKC.GetGroups(c.Request.Context(), "", userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}

	var groups []GroupUser
	for i := 0; i < len(userKeycloakGroups); i++ {
		var groupDB models.Groups

		_ = initializers.DB.
			Where("keycloak_id = ?", userKeycloakGroups[i].ID).
			First(&groupDB)

		nameGroup := strings.Split(*userKeycloakGroups[i].Name, "-")[0]
		group := GroupUser{
			IDKeycloak: *userKeycloakGroups[i].ID,
			ID:         groupDB.ID,
			Name:       nameGroup,
			SubGroups:  make(map[string]interface{}),
		}
		groups = append(groups, group)
	}
	c.JSON(http.StatusOK, gin.H{
		"Groups": groups,
	})
}
