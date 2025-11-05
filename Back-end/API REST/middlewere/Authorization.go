package middlewere

import (
	//"Go-API-T/initializers"
	//"Go-API-T/models"
	//"log"
	"net/http"
	"strings"

	//"os"
	//"time"
	keycloak "Go-API-T/Keycloak"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		accessToken := c.GetHeader("Authorization")
		refreshToken := c.Request.Header.Get("Refresh-Token")

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Access token not found",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(accessToken, "Bearer ")

		_, err := m.client.UserInfo(c.Request.Context(), tokenString)

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

			c.Writer.Header().Set("Authorization", "Bearer "+newToken.AccessToken)

			c.Request.Header.Set("Access-Token", newToken.AccessToken) //temporal

			c.Header("Refresh-Token", newToken.RefreshToken)

			tokenString = newToken.AccessToken

		}
		claims := jwt.MapClaims{}
		parser := jwt.NewParser()
		_, _, _ = parser.ParseUnverified(tokenString, claims)

		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}

		c.Next()
	}
}
