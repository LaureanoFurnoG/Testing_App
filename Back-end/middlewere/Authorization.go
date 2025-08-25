package middlewere

import (
	//"Go-API-T/initializers"
	//"Go-API-T/models"
	//"fmt"
	//"log"
	"net/http"
	//"os"
	//"time"
	keycloak "Go-API-T/Keycloak"
	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	client *keycloak.ClientKeycloak
}

func NewMiddleware(client *keycloak.ClientKeycloak) *Middleware {
	return &Middleware{
		client: client,
	}
}
func (m Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Access-Token")
		refreshToken := c.Request.Header.Get("Refresh-Token")

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Access token not found",
			})
			c.Abort()
			return
		}
		_, err := m.client.UserInfo(c.Request.Context(), accessToken)
		if err != nil {
			if refreshToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Message": "Access token expired and no refresh token provided",
				})
				c.Abort()
				return
			}

			// try refresh token
			newToken, refreshErr := m.client.RefreshToken(
				c.Request.Context(),
				refreshToken,
			)
			if refreshErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Message": "Failed to refresh token",
				})
				c.Abort()
				return
			}

			c.Request.Header.Set("Access-Token", newToken.AccessToken)

			c.Next()
			return
		}

		c.Next()
	}
}

/*
func RequireAuth(c *gin.Context) {
	//get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	//decode and validate code
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return os.Getenv("SECRET_TOKEN_KEY"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//find the user with token sub
		var user models.Users
		sub := claims["sub"].([]interface{})

		initializers.DB.First(&user, uint(sub[1].(float64)))

		if user.ID == 0{
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//attach to req
		c.Set("user", user)
		//continue
		c.Next()
	} else {
		fmt.Println(err)
	}
}
*/
