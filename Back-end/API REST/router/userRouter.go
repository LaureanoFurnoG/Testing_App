package router

import (
	"Go-API-T/Keycloak"
	"Go-API-T/middlewere"
	"Go-API-T/controllers"     //export controllers package, all packages
	"github.com/gin-gonic/gin" //gin handles the routes
)

// function that create and config the principal router of the app
func SetupRouter() *gin.Engine {
	r := gin.Default() //create a new instance of the router of gin with middlewares defaults

	api := r.Group("/api") //add all routes in /api
	clientKC := keycloak.NewClientKeycloak()
	mw := middlewere.NewMiddleware(clientKC)
	handler := controllers.NewHandlerAPI(clientKC)
	controllers.UserRoutes(api, handler, mw) //call the function in UserController to register the users routes
	controllers.GroupsController(api, handler, mw)
	controllers.TestsRoutes(api, handler, mw)
	
	return r // return the router
}
