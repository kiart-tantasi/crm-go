package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/api"
	"github.com/kiart-tantasi/crm-go/internal/middlewares"
)

func main() {
	// Gin engine
	r := gin.Default()

	// Middlewares
	middlewares.SetupMiddlewares(r)

	// Handlers
	api.SetupHandlers(r)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
