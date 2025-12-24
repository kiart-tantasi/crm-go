package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	service *Service
}

func NewContactHandler(service *Service) *ContactHandler {
	return &ContactHandler{service: service}
}

// GET /contacts
func (h *ContactHandler) ListHandler(c *gin.Context) {
	list, err := h.service.List(c.Request.Context(), c.Query("limit"), c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contacts listed successfully",
		"data":    map[string][]Contact{"contacts": list},
	})
}

// GET /contacts/:id
func (h *ContactHandler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	contact, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if contact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact retrieved successfully",
		"data":    contact,
	})
}

// POST /contacts
func (h *ContactHandler) PostHandler(c *gin.Context) {
	var input Contact
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Upsert(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact upserted successfully",
	})
}
