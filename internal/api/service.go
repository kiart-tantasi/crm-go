package api

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/contactlists"
	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/emails"
	"github.com/kiart-tantasi/crm-go/internal/health"
	"github.com/kiart-tantasi/crm-go/internal/users"
)

func SetupHandlers(r *gin.Engine, emailService *emails.Service, contactService *contacts.Service, userService *users.Service, contactListService *contactlists.Service) {
	// Health
	r.GET("/healthz", health.HealthHandler)

	// Emails
	emailHandler := emails.NewEmailHandler(emailService)
	emailsGroup := r.Group("/emails")
	emailsGroup.GET("", emailHandler.ListHandler)
	emailsGroup.GET("/:id", emailHandler.GetHandler)
	emailsGroup.POST("", emailHandler.PostHandler)

	// Contacts
	contactHandler := contacts.NewContactHandler(contactService)
	contactsGroup := r.Group("/contacts")
	contactsGroup.GET("", contactHandler.ListHandler)
	contactsGroup.GET("/:id", contactHandler.GetHandler)
	contactsGroup.POST("", contactHandler.PostHandler)

	// Contact Lists
	contactListHandler := contactlists.NewContactListHandler(contactListService)
	contactListsGroup := r.Group("/contact-lists")
	contactListsGroup.GET("", contactListHandler.ListHandler)
	contactListsGroup.GET("/:id", contactListHandler.GetHandler)
	contactListsGroup.POST("", contactListHandler.PostHandler)
	contactListsGroup.POST("/:id/contacts", contactListHandler.AddContactsHandler)
	contactListsGroup.DELETE("/:id/contacts", contactListHandler.RemoveContactsHandler)

	// Users
	userHandler := users.NewUserHandler(userService)
	usersGroup := r.Group("/users")
	usersGroup.GET("", userHandler.ListHandler)
	usersGroup.GET("/:id", userHandler.GetHandler)
	usersGroup.POST("", userHandler.PostHandler)
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
