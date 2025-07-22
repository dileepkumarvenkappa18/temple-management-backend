package entity

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	// "github.com/sharath018/temple-management-backend/utils"
)

type Handler struct {
	Service *Service
}

// NewHandler creates a new Entity handler
func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// Temple Admin → Create Temple (Triggers approval request)
func (h *Handler) CreateEntity(c *gin.Context) {
	var input Entity

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Bind Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate required dropdown fields
	if input.TempleType == "" || input.State == "" || input.EstablishedYear == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Temple Type, State, and Established Year are required"})
		return
	}

	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(auth.User).ID
	input.CreatedBy = userID

	if input.Status == "" {
		input.Status = "pending"
	}

	if err := h.Service.CreateEntity(&input, userID); err != nil {
		log.Printf("Service Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Temple registration request submitted successfully"})
}

// Super Admin or Tenant → View all temples
func (h *Handler) GetAllEntities(c *gin.Context) {
	entities, err := h.Service.GetAllEntities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch temples", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities)
}

// Anyone → View a specific temple by ID
func (h *Handler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	entity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

// Temple Admin → Update existing temple
func (h *Handler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	var input Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Update Bind Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	input.ID = uint(id)
	input.UpdatedAt = time.Now()

	if err := h.Service.UpdateEntity(input); err != nil {
		log.Printf("Update Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update temple", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temple updated successfully"})
}

// Super Admin → Delete a temple
func (h *Handler) DeleteEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	if err := h.Service.DeleteEntity(id); err != nil {
		log.Printf("Delete Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete temple", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temple deleted successfully"})
}

// func (h *Handler) UploadTempleDocument(c *gin.Context) {
// 	file, header, err := c.Request.FormFile("document")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
// 		return
// 	}
// 	defer file.Close()

// 	url, err := utils.UploadToCloudinary(file, header, "temples")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":     "Uploaded successfully",
// 		"documentURL": url,
// 	})
// }




