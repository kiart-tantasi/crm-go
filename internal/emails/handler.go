package emails

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	service *Service
}

func NewEmailHandler(service *Service) *EmailHandler {
	return &EmailHandler{service: service}
}

// GET /emails
func (h *EmailHandler) ListHandler(c *gin.Context) {
	list, err := h.service.List(c.Request.Context(), c.Query("limit"), c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Emails listed successfully",
		"data":    map[string][]Email{"emails": list},
	})
}

// GET /emails/:id
func (h *EmailHandler) GetHandler(c *gin.Context) {
	// Validate param
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	h.service.GetByID(c.Request.Context(), c.Param("id"))
}

// POST /emails
func (h *EmailHandler) PostHandler(c *gin.Context) {
	// Validate request body
	var input Email
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Upsert(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email upserted successfully",
	})
}
