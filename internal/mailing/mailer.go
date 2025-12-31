package mailing

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewMailer(host string, port int, username, password string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

type EmailParams struct {
	FromName string
	FromAddr string
	ToName   string
	ToAddr   string
	Subject  string
	Body     string
}

func (m *Mailer) Send(params EmailParams) error {
	// SMTP address
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	// Email headers
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", params.FromName, params.FromAddr)
	header["To"] = fmt.Sprintf("%s <%s>", params.ToName, params.ToAddr)
	header["Subject"] = params.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	var message strings.Builder
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	// Email body
	message.WriteString(params.Body)

	// SMTP auth
	var auth smtp.Auth
	if m.Username != "" || m.Password != "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	// Send email with basic smtp client
	// TODO: create smtp client pool
	err := smtp.SendMail(addr, auth, params.FromAddr, []string{params.ToAddr}, []byte(message.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
