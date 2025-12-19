package router

import (
	"Go-API-T/Keycloak"
	"Go-API-T/controllers" //export controllers package, all packages
	"Go-API-T/middlewere"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" //gin handles the routes
)

// function that create and config the principal router of the app
func SetupRouter() *gin.Engine {
	r := gin.Default() //create a new instance of the router of gin with middlewares defaults

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "refresh-token"},
    	ExposeHeaders:    []string{"Authorization", "Refresh-Token", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	api := r.Group("/api") //add all routes in /api
	clientKC := keycloak.NewClientKeycloak()
	mw := middlewere.NewMiddleware(clientKC)
	handler := controllers.NewHandlerAPI(clientKC)
	controllers.UserRoutes(api, handler, mw) //call the function in UserController to register the users routes
	controllers.GroupsController(api, handler, mw)
	controllers.TestsRoutes(api, handler, mw)
	
	return r // return the router
}
