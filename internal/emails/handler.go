package emails

import (
	"net/http"
	"strconv"

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
		"message": "Emails listed successfully",
		"data":    map[string][]Email{"emails": list},
	})
}

// GET /emails/:id
func (h *EmailHandler) GetHandler(c *gin.Context) {
	// Validate param
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	email, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if email == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email retrieved successfully",
		"data":    email,
	})
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

// POST /emails/:id/contact-lists
func (h *EmailHandler) AddContactListsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	var input BatchAddContactListsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddContactLists(c.Request.Context(), id, input.ContactListIDs, input.AddedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email contact lists added successfully",
	})
}

// DELETE /emails/:id/contact-lists
func (h *EmailHandler) RemoveContactListsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	var input BatchRemoveContactListsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveContactLists(c.Request.Context(), id, input.ContactListIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email contact lists removed successfully",
	})
}

// POST /emails/:id/send
// TODO: add limit param with default of 1000
// TODO: add order by contact id and min-max to prevent issues caused by replication lag
func (h *EmailHandler) SendHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	if err := h.service.Send(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent process initiated",
	})
}
