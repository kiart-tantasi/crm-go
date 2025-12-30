package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	template := flag.String("template", "", "Template for email body")
	// Others
	debugMode := flag.Bool("debug", false, "Enable debug mode")
	// Parse flags from arguments
	flag.Parse()

	// Validate required flags
	// - from email address
	// - to email address
	// - subject
	// - body email
	if *fromAddr == "" || *toAddr == "" || *subject == "" || *template == "" {
		fmt.Println("Usage of send-email:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Render email
	renderedBody, err := emails.Render(*template)
	if err != nil {
		log.Fatalf("Error rendering email: %v", err)
	}
	if *debugMode {
		log.Printf("[DEBUG] Rendered body:\n%s", renderedBody)
	}

	// Send email to SMTP server
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
