package smtppool

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/kiart-tantasi/crm-go/internal/env"
)

type Pool struct {
	host   string
	port   int
	user   string
	pass   string
	size   int
	taskCh chan mailTask
}

type mailTask struct {
	from string
	to   []string
	body []byte
}

// TODO: return error when failinng to init
func New() (*Pool, error) {
	host := env.GetEnv("SMTP_HOST", "localhost")
	portStr := env.GetEnv("SMTP_PORT", "25")
	user := env.GetEnv("SMTP_USER", "")
	pass := env.GetEnv("SMTP_PASS", "")
	sizeStr := env.GetEnv("SMTP_POOL_SIZE", "10")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP_PORT: %v", err)
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP_POOL_SIZE: %v", err)
	}

	p := &Pool{
		host:   host,
		port:   port,
		user:   user,
		pass:   pass,
		size:   size,
		taskCh: make(chan mailTask, size*10),
	}

	for i := 0; i < size; i++ {
		go p.worker(i)
	}

	return p, nil
}

func (p *Pool) SendMail(from string, to []string, body []byte) {
	p.taskCh <- mailTask{
		from: from,
		to:   to,
		body: body,
	}
}

// TODO: experiment on client pool vs non client pool
func (p *Pool) worker(id int) {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	var auth smtp.Auth
	if p.user != "" || p.pass != "" {
		auth = smtp.PlainAuth("", p.user, p.pass, p.host)
	}

	for task := range p.taskCh {
		err := smtp.SendMail(addr, auth, task.from, task.to, task.body)
		if err != nil {
			log.Printf("[SMTPPool] Worker %d: failed to send email from %s to %v: %v", id, task.from, task.to, err)
		} else {
			log.Printf("[SMTPPool] Worker %d: email sent successfully from %s to %v", id, task.from, task.to)
		}
	}
}
