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

// Super Admin → View all temples, Temple Admin → View only their created temples
func (h *Handler) GetAllEntities(c *gin.Context) {
	// Get authenticated user
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}

	var entities []Entity
	var err error

	// Check user role - if super admin, get all entities; if temple admin, get only their created entities
	// TODO: Replace "Role.RoleName" with the correct field from your auth.User struct
	if user.Role.RoleName == "superadmin" {
		entities, err = h.Service.GetAllEntities()
	} else {
		// For temple admin or other roles, get only entities created by them
		entities, err = h.Service.GetEntitiesByCreator(user.ID)
	}

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

// Temple Admin → Get devotees by entity
func (h *Handler) GetDevoteesByEntity(c *gin.Context) {
	entityIDParam := c.Param("id")
	entityIDUint, err := strconv.ParseUint(entityIDParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// ✅ Get authenticated user from context
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or missing user"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}

	// ✅ Check if the entity ID matches the user's entity (or your permission rules)
	if user.EntityID == nil || *user.EntityID != uint(entityIDUint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or missing entity"})
		return
	}

	// ✅ Fetch devotees for the given entity
	devotees, err := h.Service.GetDevotees(uint(entityIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch devotees"})
		return
	}

	c.JSON(http.StatusOK, devotees)
}

// Temple Admin → Get devotee statistics for entity
func (h *Handler) GetDevoteeStats(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	stats, err := h.Service.GetDevoteeStats(uint(entityID))
	if err != nil {
		log.Printf("Error fetching devotee stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devotee stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Temple Admin → Update devotee membership status
// PATCH /entities/:entityID/devotees/:userID/status
func (h *Handler) UpdateDevoteeMembershipStatus(c *gin.Context) {
	entityID, err := strconv.ParseUint(c.Param("entityID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.MembershipService.UpdateMembershipStatus(uint(userID), uint(entityID), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership status updated successfully"})
}

// Temple Admin → Dashboard Summary
// GET /entities/dashboard-summary
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	// Extract the authenticated user
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid entity"})
		return
	}

	// Call service
	summary, err := h.Service.GetDashboardSummary(*user.EntityID)
	if err != nil {
		log.Printf("Dashboard Summary Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
