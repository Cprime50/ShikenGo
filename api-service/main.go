package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Cprime50/api-service/middleware"
	routes "github.com/Cprime50/api-service/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Serve static html file to test firebase auth in the browser
	http.Handle("/", http.FileServer(http.Dir(".")))

	go func() {
		// Serve Gin server
		r := gin.Default()
		r.Use(cors.Default())
		client, err := middleware.InitAuth()
		if err != nil {
			log.Println(err)
			return
		}

		// Sets your email as admin on firebase
		//middleware.SetDefaultFirebaseAdmin(context.Background(), client)

		routes.RegisterProfileRoutes(r, client)
		routes.RegisterAdminRoutes(r, client)

		// Set port
		port := os.Getenv("PORT")
		if port == "" {
			port = "localhost:8080" // Default port
		}

		// Start server
		log.Printf("Gin server is running on port %s", port)
		if err := r.Run(port); err != nil {
			log.Fatalf("Failed to start Gin server: %v", err)
		}
	}()

	// Start static html server
	log.Println("Static file server is running on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Failed to start static file server: %v", err)
	}
}
