package emails

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: implement these handlers

// ListHandlers handles GET /emails
func ListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{
		{
			"id":      1,
			"name":    "Welcome Email",
			"subject": "Welcome to our platform!",
		},
	})
}

// GetHandler handles GET /emails/:id
func GetHandler(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"name":    "Dummy Email",
		"subject": "Dummy Subject",
		"body":    "Hello World",
	})
}

// CreateHandler handles POST /emails
func CreateHandler(c *gin.Context) {
	var input Email
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Email created successfully",
		"email":   input,
	})
}
