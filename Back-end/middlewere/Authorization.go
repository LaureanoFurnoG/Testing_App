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
		// parsear claims para sacar el email y guardarlo en contexto
		claims := jwt.MapClaims{}
		parser := jwt.NewParser()
		_, _, _ = parser.ParseUnverified(tokenString, claims)

		if email, ok := claims["email"].(string); ok {
			c.Request.Header.Set("email", email)
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
