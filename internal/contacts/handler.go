package contacts

import (
	"net/http"
	"strconv"

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
	limit := 100
	if limitParam := c.Query("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
	}
	offset := 0
	if offsetParam := c.Query("offset"); offsetParam != "" {
		var err error
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
	}
	list, err := h.service.List(c.Request.Context(), limit, offset)
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
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
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
