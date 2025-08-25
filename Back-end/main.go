package main

import (
	//"errors"
	//"fmt"
	//"github.com/gin-gonic/gin"
	"log"
	//"net/http"
	"Go-API-T/initializers"
	"Go-API-T/router" // import routers
	"os"

)

func init() {
	initializers.InitEnv()
	initializers.ConnectToDB()
}
func main() {
	
	r := router.SetupRouter() //call all routers
	
	err := r.SetTrustedProxies([]string{"127.0.0.1"}) //create proxies
	if err != nil{
		log.Fatalf("Error setting trusted proxies")
	}

	if serverUp := r.Run(os.Getenv("SERVER_PORT")); serverUp != nil {
		log.Fatalf("Error starting server: %v", serverUp)
	}
}
