package initializers

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"fmt"
)

var DB *gorm.DB

func ConnectToDB() {
	//load .env file
	loadEnv := godotenv.Load()
	if loadEnv != nil {
		log.Fatal("Error in load .env")
	}

	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)// hardcode values, is interesting the fmt, is equal that ${} in js

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) //connect with the db

	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
}
