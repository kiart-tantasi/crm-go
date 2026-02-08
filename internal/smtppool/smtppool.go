package smtppool

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/kiart-tantasi/crm-go/internal/env"

	"github.com/kiart-tantasi/crm-go/internal/emailsends"
)

type Pool struct {
	emailSendService emailsends.EmailSendServiceI
	host             string
	port             int
	helo             string
	user             string
	pass             string
	size             int
	taskCh           chan SMTPTask
	clients          chan *smtp.Client
	sendAttempts     int
}

type SMTPTask struct {
	from      string
	to        []string
	body      []byte
	emailID   int
	contactID int
}

func New(emailSendService emailsends.EmailSendServiceI) (*Pool, error) {
	host := env.GetEnv("SMTP_HOST", "localhost")
	helo := env.GetEnv("SMTP_HELO", "")
	user := env.GetEnv("SMTP_USER", "")
	pass := env.GetEnv("SMTP_PASS", "")
	portStr := env.GetEnv("SMTP_PORT", "25")
	sizeStr := env.GetEnv("SMTP_POOL_SIZE", "20")
	sendAttemptsStr := env.GetEnv("SMTP_SEND_ATTEMPTS", "3")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP_PORT: %v", err)
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP_POOL_SIZE: %v", err)
	}

	sendAttempts, err := strconv.Atoi(sendAttemptsStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP_SEND_ATTEMPTS: %v", err)
	}

	p := &Pool{
		emailSendService: emailSendService,
		host:             host,
		port:             port,
		helo:             helo,
		user:             user,
		pass:             pass,
		size:             size,
		taskCh:           make(chan SMTPTask, size*10),
		clients:          make(chan *smtp.Client, size),
		sendAttempts:     sendAttempts,
	}

	for i := 0; i < size; i++ {
		c, err := p.newClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create initial smtp client: %w", err)
		}
		p.clients <- c
	}

	log.Printf("[SMTPPool] Pool created with %d clients", size)
	for i := 0; i < size; i++ {
		go p.Worker(i)
	}

	return p, nil
}

func (p *Pool) EnqueueEmail(from string, to []string, body []byte, emailID, contactID int) {
	p.taskCh <- SMTPTask{
		from:      from,
		to:        to,
		body:      body,
		emailID:   emailID,
		contactID: contactID,
	}
}

func (p *Pool) Worker(id int) {
	for task := range p.taskCh {
		// Get client from pool
		client := <-p.clients

		// Run NOOP command to check client health
		// Keep creating a new client if failed
		createAttempt := 0
		for {
			createAttempt++
			err := client.Noop()
			if err == nil {
				break
			}
			log.Printf("[SMTPPool] [Worker %d] [Create attempt %d] Failed to run NOOP on a smtp client with error: %v", id, createAttempt, err)
			client.Close()
			// We can ignore error here because we will keep creating new client in this for-loop
			client, _ = p.newClient()
		}

		// Send email
		sendAttempt := 0
		for sendAttempt < p.sendAttempts {
			sendAttempt++
			err := p.sendEmail(client, task)
			if err == nil {
				// TODO: Log in debugging mode only
				log.Printf("[DEBUG] [SMTPPool] [Worker %d] [Send attempt %d] Email was successfully sent", id, sendAttempt)
				break
			}
			log.Printf("[SMTPPool] [Worker %d] [Send attempt %d] Failed to send email with error: %v", id, sendAttempt, err)
		}

		// No deadline for now
		ctx := context.Background()

		// Upsert status into table email_sends
		// TODO: add column "status" with value "sent" or "failed"
		// TODO: check error and upsert status accordingly
		if err := p.emailSendService.Upsert(ctx, task.emailID, task.contactID); err != nil {
			log.Printf("failed to upsert email sends for contact id %d: %v\n", task.contactID, err)
		}

		// Return client to pool
		p.clients <- client
	}
}

func (p *Pool) newClient() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	c, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	// HELO
	if p.helo != "" {
		if err := c.Hello(p.helo); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, nil
}

func (p *Pool) sendEmail(c *smtp.Client, task SMTPTask) error {
	if err := c.Mail(task.from); err != nil {
		return err
	}
	for _, addr := range task.to {
		if err := c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(task.body)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}
