package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/api/handlers"
	"github.com/kiart-tantasi/crm-go/internal/middleware"
	"github.com/kiart-tantasi/crm-go/internal/templates"
)

func main() {
	r := gin.Default()

	// Add the throttling middleware (200ms)
	r.Use(middleware.MinDelay(200 * time.Millisecond))

	// Health check
	r.GET("/healthz", handlers.HealthCheck)

	// Template routes
	templateGroup := r.Group("/templates")
	{
		templateGroup.GET("", templates.ListHandler)
		templateGroup.GET("/:id", templates.GetHandler)
		templateGroup.POST("", templates.CreateHandler)
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
