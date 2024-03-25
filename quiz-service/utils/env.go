package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() interface{} {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println(".env file loaded successfully")
	return nil
}

func MustHaveEnv(key string) string {
	//TODO figure out how to make this env load once first before the vars below begin laoding
	//Main.go loads this util package before loading it's own functions so placing it there doesn help
	//loadEnv()
	env := os.Getenv("ENV")
	value := os.Getenv(key)
	if value == "" && env != "test" {
		panic("Missing environment variable: " + key)
	}
	return value
}
