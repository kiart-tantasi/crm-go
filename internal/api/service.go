package api

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/emails"
	"github.com/kiart-tantasi/crm-go/internal/health"
)

func SetupHandlers(r *gin.Engine, emailService *emails.Service, contactService *contacts.Service) {
	// Health
	r.GET("/healthz", health.HealthHandler)

	// Emails
	emailHandler := emails.NewEmailHandler(emailService)
	emailsGroup := r.Group("/emails")
	{
		emailsGroup.GET("", emailHandler.ListHandler)
		emailsGroup.GET("/:id", emailHandler.GetHandler)
		emailsGroup.POST("", emailHandler.PostHandler)
	}

	// Contacts
	contactHandler := contacts.NewContactHandler(contactService)
	contactsGroup := r.Group("/contacts")
	{
		contactsGroup.GET("", contactHandler.ListHandler)
		contactsGroup.GET("/:id", contactHandler.GetHandler)
		contactsGroup.POST("", contactHandler.PostHandler)
	}
}

func (c *Client) FetchDataAndMap(url string) (map[string]any, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return data, nil
}
