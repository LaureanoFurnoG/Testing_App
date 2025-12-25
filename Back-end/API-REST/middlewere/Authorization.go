package middlewere

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	keycloak "Go-API-T/Keycloak"
	"Go-API-T/initializers"
	"Go-API-T/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		accessHeader := c.GetHeader("Authorization")
		if accessHeader == "" || !strings.HasPrefix(accessHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access token missing"})
			return
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(accessHeader, "Bearer "))
		accessSub, accessExpired, err := ExtractSubAndExp(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid access token"})
			return
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {

			if !accessExpired {
				c.Set("user_id", accessSub)
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Access expired and no refresh token",
			})
			return
		}

		refreshSub, err := ExtractSubFromJWT(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
			return
		}

		if !accessExpired {
			if accessSub != refreshSub {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Access token does not belong to refresh token user",
				})
				return
			}

			c.Set("user_id", accessSub)
			c.Next()
			return
		}

		newToken, err := m.client.RefreshToken(c.Request.Context(), refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Failed to refresh token",
			})
			return
		}

		c.SetCookie(
			"refresh_token",
			newToken.RefreshToken,
			3600*24*30,
			"/",
			"localhost",
			false,
			true,
		)
		c.Header("Access-Control-Expose-Headers", "Authorization")
		c.Header("Authorization", "Bearer "+newToken.AccessToken)

		c.Set("user_id", refreshSub)

		c.Next()
	}
}

func (m Middleware) BelongsGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupId := c.Param("groupId")
		profileCookie, cookieErr := c.Cookie("profile")
		if cookieErr != nil || profileCookie == "" {
			fmt.Println(cookieErr)
			fmt.Println("cook", profileCookie)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "You don't have profile cookie",
			})
			c.Abort()
			return
		}
		var user models.Users
		UserID := initializers.DB.First(&user, "keycloak_id = ?", profileCookie)
		if UserID.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User Missing",
			})
			c.Abort()
			return
		}

		var GroupsRelation models.GroupsRelation
		groupRleatioNFound := initializers.DB.First(&GroupsRelation, "iduser = ? AND idgroup = ?", user.ID, groupId)
		if groupRleatioNFound.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Group relation Missing",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ExtractSubAndExp(tokenStr string) (string, bool, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", false, err
	}

	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", false, fmt.Errorf("sub missing")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", false, fmt.Errorf("exp missing")
	}

	expired := time.Now().Unix() > int64(expFloat)

	return sub, expired, nil
}

func ExtractSubFromJWT(tokenStr string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub not found")
	}

	return sub, nil
}
