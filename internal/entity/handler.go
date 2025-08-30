package entity

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
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

	// GET IP ADDRESS FOR AUDIT LOGGING
	ip := middleware.GetIPFromContext(c)

	if err := h.Service.CreateEntity(&input, userID, ip); err != nil {
		log.Printf("Service Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Temple registration request submitted successfully"})
}

// Super Admin → View all temples, Temple Admin → View only their created temples, Standard User → View assigned temples
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

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Use the enhanced service method with better error handling
	entities, err := h.Service.GetEntitiesForUser(
		user.ID, 
		user.Role.RoleName, 
		accessContext.DirectEntityID, 
		accessContext.AssignedEntityID,
	)

	if err != nil {
		log.Printf("GetAllEntities error for user %d: %v", user.ID, err)
		
		// Handle specific errors
		if err == ErrNoAccessibleEntity {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No entity assigned to this user",
				"debug_info": map[string]interface{}{
					"user_id": user.ID,
					"role": user.Role.RoleName,
					"user_entity_id": user.EntityID,
					"direct_entity_id": accessContext.DirectEntityID,
					"assigned_entity_id": accessContext.AssignedEntityID,
				},
			})
			return
		}
		
		if err == ErrEntityNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Assigned entity not found",
				"debug_info": map[string]interface{}{
					"user_id": user.ID,
					"role": user.Role.RoleName,
					"user_entity_id": user.EntityID,
					"direct_entity_id": accessContext.DirectEntityID,
					"assigned_entity_id": accessContext.AssignedEntityID,
				},
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch temples", 
			"details": err.Error(),
			"debug_info": map[string]interface{}{
				"user_id": user.ID,
				"role": user.Role.RoleName,
			},
		})
		return
	}

	log.Printf("Successfully fetched %d entities for user %d (%s)", len(entities), user.ID, user.Role.RoleName)
	c.JSON(http.StatusOK, entities)
}

// Anyone → View a specific temple by ID - UPDATED ACCESS CHECK
func (h *Handler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Check permissions
	user, _ := c.Get("user")
	userObj := user.(auth.User)

	// Use the new enhanced access check method
	entityIDUint := uint(id)
	hasAccess := h.Service.CheckUserHasEntityAccess(
		userObj.ID,
		userObj.Role.RoleName,
		entityIDUint,
		accessContext.DirectEntityID,
		accessContext.AssignedEntityID,
	)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this entity",
			"debug_info": map[string]interface{}{
				"entity_id": entityIDUint,
				"user_id": userObj.ID,
				"role": userObj.Role.RoleName,
				"direct_entity_id": accessContext.DirectEntityID,
				"assigned_entity_id": accessContext.AssignedEntityID,
			},
		})
		return
	}

	entity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

// Temple Admin/Standard User → Update existing temple - UPDATED ACCESS CHECK
func (h *Handler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Check write permissions
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have write permissions"})
		return
	}

	var input Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Update Bind Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Get authenticated user for audit logging
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := user.(auth.User)
	userID := userObj.ID

	// Use the new enhanced access check method
	entityIDUint := uint(id)
	hasAccess := h.Service.CheckUserHasEntityAccess(
		userObj.ID,
		userObj.Role.RoleName,
		entityIDUint,
		accessContext.DirectEntityID,
		accessContext.AssignedEntityID,
	)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to update this entity",
			"debug_info": map[string]interface{}{
				"entity_id": entityIDUint,
				"user_id": userObj.ID,
				"role": userObj.Role.RoleName,
				"direct_entity_id": accessContext.DirectEntityID,
				"assigned_entity_id": accessContext.AssignedEntityID,
			},
		})
		return
	}

	input.ID = uint(id)
	input.UpdatedAt = time.Now()

	// GET IP ADDRESS FOR AUDIT LOGGING
	ip := middleware.GetIPFromContext(c)

	if err := h.Service.UpdateEntity(input, userID, ip); err != nil {
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

	// Get authenticated user for audit logging
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(auth.User).ID

	// GET IP ADDRESS FOR AUDIT LOGGING
	ip := middleware.GetIPFromContext(c)

	if err := h.Service.DeleteEntity(id, userID, ip); err != nil {
		log.Printf("Delete Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete temple", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temple deleted successfully"})
}

// Temple Admin/Standard User → Get devotees by entity - UPDATED ACCESS CHECK
func (h *Handler) GetDevoteesByEntity(c *gin.Context) {
	entityIDParam := c.Param("id")
	entityIDUint, err := strconv.ParseUint(entityIDParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Get user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := user.(auth.User)

	// Use the new enhanced access check method
	entityID := uint(entityIDUint)
	hasAccess := h.Service.CheckUserHasEntityAccess(
		userObj.ID,
		userObj.Role.RoleName,
		entityID,
		accessContext.DirectEntityID,
		accessContext.AssignedEntityID,
	)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to devotees for this entity",
			"debug_info": map[string]interface{}{
				"entity_id": entityID,
				"user_id": userObj.ID,
				"role": userObj.Role.RoleName,
				"direct_entity_id": accessContext.DirectEntityID,
				"assigned_entity_id": accessContext.AssignedEntityID,
			},
		})
		return
	}

	// Fetch devotees for the given entity
	devotees, err := h.Service.GetDevotees(uint(entityIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch devotees"})
		return
	}

	c.JSON(http.StatusOK, devotees)
}

// Temple Admin/Standard User → Get devotee statistics for entity - UPDATED ACCESS CHECK
func (h *Handler) GetDevoteeStats(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityIDUint, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Get user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := user.(auth.User)

	// Use the new enhanced access check method
	entityID := uint(entityIDUint)
	hasAccess := h.Service.CheckUserHasEntityAccess(
		userObj.ID,
		userObj.Role.RoleName,
		entityID,
		accessContext.DirectEntityID,
		accessContext.AssignedEntityID,
	)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to devotee stats for this entity",
			"debug_info": map[string]interface{}{
				"entity_id": entityID,
				"user_id": userObj.ID,
				"role": userObj.Role.RoleName,
				"direct_entity_id": accessContext.DirectEntityID,
				"assigned_entity_id": accessContext.AssignedEntityID,
			},
		})
		return
	}

	stats, err := h.Service.GetDevoteeStats(entityID)
	if err != nil {
		log.Printf("Error fetching devotee stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devotee stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Temple Admin/Standard User → Update devotee membership status - UPDATED ACCESS CHECK
// PATCH /entities/:entityID/devotees/:userID/status
func (h *Handler) UpdateDevoteeMembershipStatus(c *gin.Context) {
	entityIDUint, err := strconv.ParseUint(c.Param("entityID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Check write permissions
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have write permissions"})
		return
	}

	// Get user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := user.(auth.User)

	// Use the new enhanced access check method
	entityID := uint(entityIDUint)
	hasAccess := h.Service.CheckUserHasEntityAccess(
		userObj.ID,
		userObj.Role.RoleName,
		entityID,
		accessContext.DirectEntityID,
		accessContext.AssignedEntityID,
	)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to manage devotees for this entity",
			"debug_info": map[string]interface{}{
				"entity_id": entityID,
				"user_id": userObj.ID,
				"role": userObj.Role.RoleName,
				"direct_entity_id": accessContext.DirectEntityID,
				"assigned_entity_id": accessContext.AssignedEntityID,
			},
		})
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

	err = h.Service.MembershipService.UpdateMembershipStatus(uint(userID), uint(entityIDUint), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership status updated successfully"})
}

// Temple Admin/Standard User → Dashboard Summary
// GET /entities/dashboard-summary
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	// Get access context
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// Get the accessible entity ID
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No accessible entity"})
		return
	}

	// Call service
	summary, err := h.Service.GetDashboardSummary(*entityID)
	if err != nil {
		log.Printf("Dashboard Summary Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}