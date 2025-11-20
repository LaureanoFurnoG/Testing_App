package initializers

import (
	"log"
	"github.com/joho/godotenv" //is the .env import
)

func InitEnv() {
	err := godotenv.Load(".env") // load .env
	if err != nil {
		log.Fatal("Error in load .env", err)
	}
}
