package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kiart-tantasi/crm-go/internal/api"
	"github.com/kiart-tantasi/crm-go/internal/email"
	"github.com/kiart-tantasi/crm-go/internal/emails"
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
	bodyEmail := flag.String("body-email", "", "Email body (Go email)")
	apiURL := flag.String("api-url", "", "API URL to fetch data for email")
	// Others
	debugMode := flag.Bool("debug", false, "Enable debug mode")
	// Parse flags from arguments
	flag.Parse()

	// Validate required flags
	// - from email address
	// - to email address
	// - subject
	// - body email
	if *fromAddr == "" || *toAddr == "" || *subject == "" || *bodyEmail == "" {
		fmt.Println("Usage of send-email:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. Fetch data from API
	data := map[string]any{}
	if *apiURL != "" {
		client := api.NewClient()
		var err error
		data, err = client.FetchDataAndMap(*apiURL)
		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}
	}

	// 2. Render email
	renderedBody, err := emails.Render(*bodyEmail, data)
	if err != nil {
		log.Fatalf("Error rendering email: %v", err)
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
