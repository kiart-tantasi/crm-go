package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/contactlists"
	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/db"
	"github.com/kiart-tantasi/crm-go/internal/emails"
	"github.com/kiart-tantasi/crm-go/internal/env"
	"github.com/kiart-tantasi/crm-go/internal/handlers"
	"github.com/kiart-tantasi/crm-go/internal/middlewares"
	"github.com/kiart-tantasi/crm-go/internal/users"
)

func main() {
	// Load env files
	if err := env.LoadEnvFiles(); err != nil {
		log.Printf("Warning: %v", err)
	}

	// Gin engine
	r := gin.Default()

	// Init DB connection
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Shared repositories and services
	emailRepo := emails.NewRepository(db)
	emailService := emails.NewService(emailRepo)
	contactRepo := contacts.NewRepository(db)
	contactService := contacts.NewService(contactRepo)
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	contactListRepo := contactlists.NewRepository(db)
	contactListService := contactlists.NewService(contactListRepo)

	// Middlewares
	middlewares.SetupMiddlewares(r)

	// Handlers
	handlers.SetupHandlers(r, emailService, contactService, userService, contactListService)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
