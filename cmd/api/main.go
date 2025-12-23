package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/api"
	"github.com/kiart-tantasi/crm-go/internal/emails"
	"github.com/kiart-tantasi/crm-go/internal/middlewares"
)

func main() {
	// Gin engine
	r := gin.Default()

	// Shared repositories and services
	// TODO: get db credentials from env
	db, err := sql.Open("mysql", "admin:admin@tcp(localhost:3309)/crm-go")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	emailRepo := emails.NewRepository(db)
	emailService := emails.NewService(emailRepo)

	// Middlewares
	middlewares.SetupMiddlewares(r)

	// Handlers
	api.SetupHandlers(r, emailService)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
