package api

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kiart-tantasi/crm-go/internal/health"
	"github.com/kiart-tantasi/crm-go/internal/templates"
)

func SetupHandlers(r *gin.Engine) {
	// Health
	r.GET("/healthz", health.HealthHandler)

	// Templates
	templatesGroup := r.Group("/templates")
	{
		templatesGroup.GET("", templates.ListHandler)
		templatesGroup.GET("/:id", templates.GetHandler)
		templatesGroup.POST("", templates.CreateHandler)
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
