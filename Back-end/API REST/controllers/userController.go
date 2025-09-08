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



/*
	func register(c *gin.Context) {
		var jsonData struct {
			Name     string
			Lastname string
			Email    string
			Password string
		} //create body request

		if c.ShouldBindJSON(&jsonData) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body",
			})
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(jsonData.Password), 10)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.Users{Name: jsonData.Name, Lastname: jsonData.Lastname, Email: jsonData.Email, Password: string(hash)}
		userExist := initializers.DB.First(&user, "email = ?", jsonData.Email)

		if userExist.Error == nil {

c.JSON(409, gin.H{"error": "This email already in use"})

			return
		}

		createU := initializers.DB.Create(&user) //create user

		if createU.Error != nil {
			c.JSON(400, gin.H{"error": createU.Error})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Message": "User created"})
	}

	func login(c *gin.Context) {
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
		var user models.Users
		userFound := initializers.DB.First(&user, "email = ?", jsonData.Email)

		if userFound.Error != nil {
			c.JSON(404, gin.H{
				"error": "User not found",
			})
			return
		}

		comparePassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(jsonData.Password))
		if comparePassword != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password or email not match",
			})
			return
		}
		//generate code
		var code2step string
		for i := 1; i <= 5; i++ {
			code2step += fmt.Sprintf("%d", rand.Intn(10))
		}
		fmt.Println(code2step)
		//generate jwt for TwoStep code
		tokenCodeTwo_STEP := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": []interface{}{user.Email, code2step},
			"exp": time.Now().Add(time.Minute * 2).Unix(),
		})

		tokenString, err := tokenCodeTwo_STEP.SignedString([]byte(os.Getenv("SECRET_TOKEN_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}

func twoStep(c *gin.Context) {
	var jsonData struct {
		Code           string
		TokenEncrypter string
	}

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	partsToken := strings.Split(jsonData.TokenEncrypter, ".")

	if len(partsToken) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token malformed",
		})
	}

	payloadBase64 := partsToken[1]

	if l := len(payloadBase64) % 4; l > 0 {
		payloadBase64 += strings.Repeat("=", 4-l)
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(payloadBase64)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	var payload map[string]interface{}

	err = json.Unmarshal(payloadBytes, &payload)

	if err != nil {
		fmt.Println("Error parsing to JSON:", err)
		return
	}

	//eval if the token was expired
	if expVal, ok := payload["exp"].(float64); ok {
		expTime := time.Unix(int64(expVal), 0)
		if time.Now().After(expTime) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing expiration in token",
		})
		return
	}

	var CodeDecrypter string

	if sub, ok := payload["sub"].([]interface{}); ok && len(sub) >= 2 {
		if code, ok := sub[1].(string); ok {
			CodeDecrypter = code
		}
	} else {
		fmt.Println("'sub' not is a string")
	}

	var user models.Users
	if strings.TrimSpace(CodeDecrypter) == strings.TrimSpace(jsonData.Code) {
		tokenCodeTwo_STEP := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": []interface{}{user.Email, user.ID},
			"exp": time.Now().Add(time.Minute * 2).Unix(),
		})

		tokenString, err := tokenCodeTwo_STEP.SignedString([]byte(os.Getenv("SECRET_TOKEN_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true) //create cookie

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Code not match or expired",
		})
		return
	}

	//la idea: decodificar el token, tomar el email de usuario, verificar que el codigo sea igual que el encriptado y devolver token de usuario con id, email
	//userFound := initializers.DB.First(&user, "email = ?", jsonData.Email)
}
*/
