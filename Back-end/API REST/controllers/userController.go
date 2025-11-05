package controllers

import (
	keycloak "Go-API-T/Keycloak"
	"Go-API-T/middlewere"
	"Go-API-T/services"
	"crypto/rand"
	"math/big"

	"Go-API-T/initializers"
	"Go-API-T/models"
	"fmt"

	//"strings"
	//"time"
	//"example/Go-API-T/services"
	//"encoding/base64"
	//"encoding/json"
	//"context"
	//"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v5"
	//"golang.org/x/crypto/bcrypt"
	//"strconv"

)

var otpStore = make(map[string]string)


// function to save all endpoints in a "router".
func UserRoutes(rg *gin.RouterGroup, handler *HandlerAPI, mw *middlewere.Middleware) {
	user := rg.Group("/user") //prefix that all routes(endpoints)

	user.POST("/register", handler.register)
	user.POST("/login", handler.login)
	user.POST("/verify", handler.Verify2Step)

	//user.POST("/TwoStep", twoStep)
}

// define a
type HandlerAPI struct {
	clientKC *keycloak.ClientKeycloak
}

func NewHandlerAPI(client *keycloak.ClientKeycloak) *HandlerAPI {
	return &HandlerAPI{
		clientKC: client,
	}
}

func generate2Step() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%s", n)
}

// send code

func (h *HandlerAPI) register(c *gin.Context) {
	var jsonData struct {
		Username string
		Name     string
		Lastname string
		Email    string
		Password string
	} //create body request

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	params := keycloak.CreateUserParams{
		Username: jsonData.Username,
		Name:     jsonData.Name,
		Lastname: jsonData.Lastname,
		Email:    jsonData.Email,
		Password: jsonData.Password,
	}

	userID, err := h.clientKC.CreateUser(c.Request.Context(), params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	user := models.Users{KeycloakID: userID}

	createU := initializers.DB.Create(&user)

	if createU.Error != nil {
		c.JSON(400, gin.H{"error": createU.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "User created"})
}

type SignInResponse struct {
	AccessToken  string
	RefreshToken string
}

func (h *HandlerAPI) login(c *gin.Context) {
	var jsonData struct {
		Email    string
		Password string
	}

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	jwt, err := h.clientKC.Login(c.Request.Context(), jsonData.Email, jsonData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	otp := generate2Step()
	otpStore[jsonData.Email] = otp

	if err := services.SendEmail(jsonData.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP"})
		return
	}

	otpStore[jsonData.Email+"_AccessToken"] = jwt.AccessToken
	otpStore[jsonData.Email+"_RefreshToken"] = jwt.RefresToken

	c.JSON(http.StatusOK, gin.H{
		"message": "Code sended",
	})
}

func (h *HandlerAPI) Verify2Step(c *gin.Context) {
	var body struct {
		Email string
		OTP   string
	}

	if c.ShouldBindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	expectedOtp, exists := otpStore[body.Email]

	if !exists || expectedOtp != body.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
		return
	}

	AccessToken := otpStore[body.Email+"_AccessToken"]
	RefreshToken := otpStore[body.Email+"_RefreshToken"]

	delete(otpStore, body.Email)
	delete(otpStore, body.Email+"_AccessToken")
	delete(otpStore, body.Email+"_RefreshToken")

	signInResp := SignInResponse{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successful",
		"data":    signInResp,
	})
}

type MyKeycloakClient struct {
	*keycloak.ClientKeycloak
}
