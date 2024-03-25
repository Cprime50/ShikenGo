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

// var (
// 	_   = loadEnv()
// 	ENV = mustHaveEnv("ENV")
// 	// HTTP_PORT            = mustHaveEnv("HTTP_PORT")
// 	GRPC_PORT = mustHaveEnv("GRPC_PORT")
// 	// CLIENT_URL           = mustHaveEnv("CLIENT_URL")
// 	// SERVER_HTTP          = mustHaveEnv("SERVER_HTTP")
// 	// COOKIE_DOMAIN        = mustHaveEnv("COOKIE_DOMAIN")
// 	CERT_PATH = mustHaveEnv("CERT_PATH")
// 	KEY_PATH  = mustHaveEnv("KEY_PATH")
// 	// GOOGLE_CLIENT_ID     = mustHaveEnv("GOOGLE_CLIENT_ID")
// 	// GOOGLE_CLIENT_SECRET = mustHaveEnv("GOOGLE_CLIENT_SECRET")
// 	// GITHUB_CLIENT_ID     = mustHaveEnv("GITHUB_CLIENT_ID")
// 	// GITHUB_CLIENT_SECRET = mustHaveEnv("GITHUB_CLIENT_SECRET")
// )
