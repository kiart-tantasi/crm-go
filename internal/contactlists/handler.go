package contactlists

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactListHandler struct {
	service *Service
}

func NewContactListHandler(service *Service) *ContactListHandler {
	return &ContactListHandler{service: service}
}

// GET /contact-lists
func (h *ContactListHandler) ListHandler(c *gin.Context) {
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
		"message": "Contact lists listed successfully",
		"data":    map[string][]ContactList{"contact_lists": list},
	})
}

// GET /contact-lists/:id
func (h *ContactListHandler) GetHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	cl, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contact list not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact list retrieved successfully",
		"data":    cl,
	})
}

// POST /contact-lists
func (h *ContactListHandler) PostHandler(c *gin.Context) {
	var input ContactList
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Insert(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact list inserted successfully",
	})
}

// POST /contact-lists/:id/contacts
func (h *ContactListHandler) AddContactsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact list id"})
		return
	}

	var input BatchAddContactsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddContacts(c.Request.Context(), id, input.Contacts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact list contacts added successfully",
	})
}

// DELETE /contact-lists/:id/contacts
func (h *ContactListHandler) RemoveContactsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact list id"})
		return
	}

	var input BatchRemoveContactsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveContacts(c.Request.Context(), id, input.ContactIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact list contacts removed successfully",
	})
}
