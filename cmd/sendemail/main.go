package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kiart-tantasi/crm-go/internal/api"
	"github.com/kiart-tantasi/crm-go/internal/email"
	"github.com/kiart-tantasi/crm-go/internal/templates"
)

func main() {
	// [Flags]
	// SMTP
	smtpHost := flag.String("smtp-host", "localhost", "SMTP Host")
	smtpPort := flag.Int("smtp-port", 1025, "SMTP Port")
	smtpUser := flag.String("smtp-user", "", "SMTP Username")
	smtpPass := flag.String("smtp-pass", os.Getenv("SMTP_PASS"), "SMTP Password (defaults to SMTP_PASS env var)")
	// From/To
	fromName := flag.String("from-name", "", "Sender name")
	fromAddr := flag.String("from-addr", "", "Sender email address")
	toName := flag.String("to-name", "", "Recipient name")
	toAddr := flag.String("to-addr", "", "Recipient email address")
	// Payload
	subject := flag.String("subject", "", "Email subject")
	bodyTemplate := flag.String("body-template", "", "Email body template (Go template)")
	apiURL := flag.String("api-url", "", "API URL to fetch data for template")
	// Others
	debugMode := flag.Bool("debug", false, "Enable debug mode")
	// Parse flags from arguments
	flag.Parse()

	// Validate required flags
	// - from email address
	// - to email address
	// - subject
	// - body template
	if *fromAddr == "" || *toAddr == "" || *subject == "" || *bodyTemplate == "" {
		fmt.Println("Usage of send-email:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. Fetch data from API
	data := map[string]any{}
	if *apiURL != "" {
		client := api.NewClient()
		var err error
		data, err = client.FetchData(*apiURL)
		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}
	}

	// 2. Render template
	renderedBody, err := templates.Render(*bodyTemplate, data)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}
	if *debugMode {
		log.Printf("[DEBUG] Rendered body:\n%s", renderedBody)
	}

	// 3. Send email
	mailer := email.NewMailer(*smtpHost, *smtpPort, *smtpUser, *smtpPass)
	err = mailer.Send(email.EmailParams{
		FromName: *fromName,
		FromAddr: *fromAddr,
		ToName:   *toName,
		ToAddr:   *toAddr,
		Subject:  *subject,
		Body:     renderedBody,
	})

	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	log.Println("Email sent successfully!")
}
