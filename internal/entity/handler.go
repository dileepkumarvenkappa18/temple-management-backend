package entity

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// ========== ENTITY CORE ==========

// Temple Admin → Create Temple (Triggers Approval Request)
func (h *Handler) CreateEntity(c *gin.Context) {
	var input Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := user.(auth.User).ID

	if err := h.Service.CreateEntity(&input, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Temple creation request submitted for approval"})
}

// Super Admin / Temple Admin → View All Temples
func (h *Handler) GetAllEntities(c *gin.Context) {
	entities, err := h.Service.GetAllEntities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entities"})
		return
	}
	c.JSON(http.StatusOK, entities)
}

// Anyone → View Specific Temple
func (h *Handler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	entity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

// Temple Admin → Update Own Temple (Optional)
func (h *Handler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	var e Entity
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	e.ID = uint(id)
	e.UpdatedAt = time.Now()

	if err := h.Service.UpdateEntity(e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entity: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entity updated successfully"})
}

// Super Admin → Delete a Temple (Optional)
func (h *Handler) DeleteEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	if err := h.Service.DeleteEntity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entity: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted successfully"})
}
