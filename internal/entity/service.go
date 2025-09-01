package entity

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
)

type MembershipService interface {
	UpdateMembershipStatus(userID, entityID uint, status string) error
}

type Service struct {
	Repo              *Repository
	MembershipService MembershipService
	AuditService      auditlog.Service
}

func NewService(r *Repository, ms MembershipService, as auditlog.Service) *Service {
	return &Service{
		Repo:              r,
		MembershipService: ms,
		AuditService:      as,
	}
}

var (
	ErrMissingFields = errors.New("temple name, deity, phone, and email are required")
	ErrEntityNotFound = errors.New("entity not found")
	ErrNoAccessibleEntity = errors.New("no accessible entity found for user")
)

// ========== ENTITY CORE ==========
func (s *Service) CreateEntity(e *Entity, userID uint, ip string) error {
    // Validate required fields
    if strings.TrimSpace(e.Name) == "" ||
        e.MainDeity == nil || strings.TrimSpace(*e.MainDeity) == "" ||
        strings.TrimSpace(e.Phone) == "" ||
        strings.TrimSpace(e.Email) == "" {
        
        // 🔍 LOG FAILED TEMPLE CREATION ATTEMPT
        auditDetails := map[string]interface{}{
            "temple_name": strings.TrimSpace(e.Name),
            "email":       strings.TrimSpace(e.Email),
            "error":       "Missing required fields",
        }
        s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")
        
        return ErrMissingFields
    }

    now := time.Now()

    // Set metadata
    e.Status = "pending"
    // DO NOT override the CreatedBy field here - it's already set by the handler
    // The CreatedBy field should already be set correctly
    e.CreatedAt = now
    e.UpdatedAt = now

    // Sanitize inputs
    e.Name = strings.TrimSpace(e.Name)
    e.Email = strings.TrimSpace(e.Email)
    e.Phone = strings.TrimSpace(e.Phone)
    e.TempleType = strings.TrimSpace(e.TempleType)
    e.Description = strings.TrimSpace(e.Description)
    e.StreetAddress = strings.TrimSpace(e.StreetAddress)
    e.City = strings.TrimSpace(e.City)
    e.State = strings.TrimSpace(e.State)
    e.District = strings.TrimSpace(e.District)
    e.Pincode = strings.TrimSpace(e.Pincode)

    // Trim main deity if present
    if e.MainDeity != nil {
        trimmed := strings.TrimSpace(*e.MainDeity)
        e.MainDeity = &trimmed
    }

    // Trim document URLs
    e.RegistrationCertURL = strings.TrimSpace(e.RegistrationCertURL)
    e.TrustDeedURL = strings.TrimSpace(e.TrustDeedURL)
    e.PropertyDocsURL = strings.TrimSpace(e.PropertyDocsURL)
    e.AdditionalDocsURLs = strings.TrimSpace(e.AdditionalDocsURLs)

    // Save entity
    if err := s.Repo.CreateEntity(e); err != nil {
        // 🔍 LOG FAILED TEMPLE CREATION (DB ERROR)
        auditDetails := map[string]interface{}{
            "temple_name": e.Name,
            "email":       e.Email,
            "error":       err.Error(),
        }
        s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")
        
        return err
    }

    // Create approval request
    req := &auth.ApprovalRequest{
        UserID:      userID,
        EntityID:    &e.ID,
        RequestType: "temple_approval",
        Status:      "pending",
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    if err := s.Repo.CreateApprovalRequest(req); err != nil {
        // 🔍 LOG FAILED APPROVAL REQUEST CREATION
        auditDetails := map[string]interface{}{
            "temple_name": e.Name,
            "temple_id":   e.ID,
            "email":       e.Email,
            "error":       err.Error(),
        }
        s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_APPROVAL_REQUEST_FAILED", auditDetails, ip, "failure")
        
        return err
    }

    // 🔍 LOG SUCCESSFUL TEMPLE CREATION
    auditDetails := map[string]interface{}{
        "temple_name":   e.Name,
        "temple_id":     e.ID,
        "temple_type":   e.TempleType,
        "email":         e.Email,
        "phone":         e.Phone,
        "city":          e.City,
        "state":         e.State,
        "main_deity":    e.MainDeity,
        "status":        e.Status,
    }
    s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_CREATED", auditDetails, ip, "success")

    return nil
}





// Super Admin → Get all temples
func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

// Temple Admin → Get entities created by specific user
func (s *Service) GetEntitiesByCreator(creatorID uint) ([]Entity, error) {
	return s.Repo.GetEntitiesByCreator(creatorID)
}

// NEW: Get all entities created by users within the same tenant/entity context
func (s *Service) GetEntitiesCreatedByTenantUsers(entityID uint) ([]Entity, error) {
	log.Printf("Getting entities created by users within entity %d context", entityID)
	
	// Get all users who have access to this entity (through membership, direct assignment, etc.)
	var userIDs []uint
	
	// Get users through memberships
	err := s.Repo.DB.Table("user_entity_memberships").
		Where("entity_id = ?", entityID).
		Distinct("user_id").
		Pluck("user_id", &userIDs).Error
	
	if err != nil {
		log.Printf("Error getting user IDs from memberships for entity %d: %v", entityID, err)
		return []Entity{}, err
	}
	
	// Get users who have this entity as their direct entity
	var directUserIDs []uint
	err = s.Repo.DB.Table("users").
		Where("entity_id = ?", entityID).
		Pluck("id", &directUserIDs).Error
	
	if err != nil {
		log.Printf("Error getting direct user IDs for entity %d: %v", entityID, err)
	} else {
		// Merge direct user IDs with membership user IDs
		userIDMap := make(map[uint]bool)
		for _, id := range userIDs {
			userIDMap[id] = true
		}
		for _, id := range directUserIDs {
			if !userIDMap[id] {
				userIDs = append(userIDs, id)
				userIDMap[id] = true
			}
		}
	}
	
	log.Printf("Found %d users with access to entity %d", len(userIDs), entityID)
	
	if len(userIDs) == 0 {
		return []Entity{}, nil
	}
	
	// Get all entities created by these users
	var entities []Entity
	err = s.Repo.DB.Where("created_by IN ?", userIDs).Find(&entities).Error
	if err != nil {
		log.Printf("Error fetching entities created by tenant users for entity %d: %v", entityID, err)
		return []Entity{}, err
	}
	
	log.Printf("Found %d entities created by users within entity %d context", len(entities), entityID)
	return entities, nil
}

// Anyone → View a temple by ID
func (s *Service) GetEntityByID(id int) (Entity, error) {
	if id <= 0 {
		return Entity{}, ErrEntityNotFound
	}
	
	entity, err := s.Repo.GetEntityByID(id)
	if err != nil {
		log.Printf("Entity not found with ID %d: %v", id, err)
		return Entity{}, ErrEntityNotFound
	}
	
	return entity, nil
}

// Enhanced method to get entities for a specific user with cross-tenant visibility - FIXED
func (s *Service) GetEntitiesForUser(userID uint, role string, directEntityID *uint, assignedEntityID *uint) ([]Entity, error) {
	log.Printf("GetEntitiesForUser - UserID: %d, Role: %s, DirectEntityID: %v, AssignedEntityID: %v", 
		userID, role, directEntityID, assignedEntityID)

	switch role {
	case "superadmin":
		return s.GetAllEntities()
		
	case "templeadmin":
		var entities []Entity
		entityIDMap := make(map[uint]bool) // To avoid duplicates
		
		// Priority 1: DirectEntityID
		if directEntityID != nil {
			entity, err := s.GetEntityByID(int(*directEntityID))
			if err == nil {
				entities = append(entities, entity)
				entityIDMap[entity.ID] = true
			} else {
				log.Printf("DirectEntityID %d not found for templeadmin %d: %v", *directEntityID, userID, err)
			}
		}
		
		// Priority 2: Entities created by them
		createdEntities, err := s.GetEntitiesByCreator(userID)
		if err != nil {
			log.Printf("Failed to get entities by creator for user %d: %v", userID, err)
		} else {
			for _, created := range createdEntities {
				if !entityIDMap[created.ID] {
					entities = append(entities, created)
					entityIDMap[created.ID] = true
				}
			}
		}
		
		// Priority 3: NEW - Get entities created by other users within their accessible entities
		accessibleEntityIDs := s.getAccessibleEntityIDs(userID, role, directEntityID, assignedEntityID)
		for _, accessibleEntityID := range accessibleEntityIDs {
			tenantEntities, err := s.GetEntitiesCreatedByTenantUsers(accessibleEntityID)
			if err != nil {
				log.Printf("Failed to get tenant entities for entity %d: %v", accessibleEntityID, err)
				continue
			}
			
			for _, tenantEntity := range tenantEntities {
				if !entityIDMap[tenantEntity.ID] {
					entities = append(entities, tenantEntity)
					entityIDMap[tenantEntity.ID] = true
				}
			}
		}
		
		if len(entities) == 0 {
			log.Printf("No entities found for templeadmin %d", userID)
		}
		return entities, nil
		
	case "standarduser", "monitoringuser":
		var entities []Entity
		var hasAnyErrors bool = false
		entityIDMap := make(map[uint]bool) // To avoid duplicates
		
		// Priority 1: AssignedEntityID
		if assignedEntityID != nil {
			log.Printf("Checking assigned entity ID %d for %s %d", *assignedEntityID, role, userID)
			
			if !s.CheckEntityExists(*assignedEntityID) {
				log.Printf("CRITICAL: Assigned entity %d does not exist in database for %s %d", *assignedEntityID, role, userID)
				hasAnyErrors = true
			} else {
				log.Printf("Assigned entity %d exists in database for %s %d", *assignedEntityID, role, userID)
				
				entity, err := s.GetEntityByID(int(*assignedEntityID))
				if err == nil {
					log.Printf("Successfully fetched assigned entity %d for %s %d", *assignedEntityID, role, userID)
					entities = append(entities, entity)
					entityIDMap[entity.ID] = true
				} else {
					log.Printf("Error fetching assigned entity %d for %s %d: %v", *assignedEntityID, role, userID, err)
					hasAnyErrors = true
				}
			}
		}
		
		// Priority 2: DirectEntityID
		if directEntityID != nil {
			entity, err := s.GetEntityByID(int(*directEntityID))
			if err == nil {
				if !entityIDMap[entity.ID] {
					log.Printf("Found direct entity %d for %s %d", *directEntityID, role, userID)
					entities = append(entities, entity)
					entityIDMap[entity.ID] = true
				}
			} else {
				log.Printf("DirectEntityID %d not found for %s %d: %v", *directEntityID, role, userID, err)
				hasAnyErrors = true
			}
		}
		
		// Priority 3: Entities created by the user
		createdEntities, err := s.GetEntitiesByCreator(userID)
		if err != nil {
			log.Printf("Failed to get entities by creator for user %d: %v", userID, err)
			hasAnyErrors = true
		} else {
			log.Printf("Found %d entities created by user %d", len(createdEntities), userID)
			for _, created := range createdEntities {
				if !entityIDMap[created.ID] {
					log.Printf("Found created entity %d for %s %d", created.ID, role, userID)
					entities = append(entities, created)
					entityIDMap[created.ID] = true
				}
			}
		}
		
		// Priority 4: Entities through memberships
		membershipEntities, err := s.GetEntitiesWithUserAccess(userID)
		if err != nil {
			log.Printf("Failed to get entities with user access for user %d: %v", userID, err)
			hasAnyErrors = true
		} else {
			log.Printf("Found %d entities through memberships for user %d", len(membershipEntities), userID)
			for _, membership := range membershipEntities {
				if !entityIDMap[membership.ID] {
					log.Printf("Found membership entity %d for %s %d", membership.ID, role, userID)
					entities = append(entities, membership)
					entityIDMap[membership.ID] = true
				}
			}
		}
		
		// Priority 5: NEW - Get entities created by other users within their accessible entities
		accessibleEntityIDs := s.getAccessibleEntityIDs(userID, role, directEntityID, assignedEntityID)
		for _, accessibleEntityID := range accessibleEntityIDs {
			tenantEntities, err := s.GetEntitiesCreatedByTenantUsers(accessibleEntityID)
			if err != nil {
				log.Printf("Failed to get tenant entities for entity %d: %v", accessibleEntityID, err)
				continue
			}
			
			for _, tenantEntity := range tenantEntities {
				if !entityIDMap[tenantEntity.ID] {
					log.Printf("Found tenant entity %d created by other users for %s %d", tenantEntity.ID, role, userID)
					entities = append(entities, tenantEntity)
					entityIDMap[tenantEntity.ID] = true
				}
			}
		}
		
		// Return results - only return error if no entities found AND we have assigned/direct entity IDs that should exist
		if len(entities) == 0 {
			log.Printf("No accessible entities found for %s %d", role, userID)
			log.Printf("Debug: assignedEntityID=%v, directEntityID=%v, hasAnyErrors=%v", assignedEntityID, directEntityID, hasAnyErrors)
			
			// If we have assigned or direct entity IDs but couldn't find them, return EntityNotFound
			if (assignedEntityID != nil || directEntityID != nil) && hasAnyErrors {
				log.Printf("Returning ErrEntityNotFound because assigned/direct entity was not found")
				return []Entity{}, ErrEntityNotFound
			}
			
			// Otherwise, return NoAccessibleEntity
			log.Printf("Returning ErrNoAccessibleEntity because no entities were found")
			return []Entity{}, ErrNoAccessibleEntity
		}
		
		log.Printf("Found %d total entities for %s %d", len(entities), role, userID)
		return entities, nil
		
	default:
		return []Entity{}, errors.New("unsupported role: " + role)
	}
}

// NEW: Helper method to get accessible entity IDs for a user
func (s *Service) getAccessibleEntityIDs(userID uint, role string, directEntityID *uint, assignedEntityID *uint) []uint {
	var entityIDs []uint
	entityIDMap := make(map[uint]bool)
	
	// Add assigned entity ID
	if assignedEntityID != nil {
		entityIDs = append(entityIDs, *assignedEntityID)
		entityIDMap[*assignedEntityID] = true
	}
	
	// Add direct entity ID
	if directEntityID != nil && !entityIDMap[*directEntityID] {
		entityIDs = append(entityIDs, *directEntityID)
		entityIDMap[*directEntityID] = true
	}
	
	// Add entities through memberships
	var membershipEntityIDs []uint
	err := s.Repo.DB.Table("user_entity_memberships").
		Where("user_id = ?", userID).
		Pluck("entity_id", &membershipEntityIDs).Error
	
	if err == nil {
		for _, id := range membershipEntityIDs {
			if !entityIDMap[id] {
				entityIDs = append(entityIDs, id)
				entityIDMap[id] = true
			}
		}
	}
	
	return entityIDs
}

// GetEntitiesWithUserAccess - Find entities user has access to through memberships
func (s *Service) GetEntitiesWithUserAccess(userID uint) ([]Entity, error) {
	log.Printf("Looking for entity memberships for user %d", userID)
	
	// First try to find entities through user_entity_memberships
	var entityIDs []uint
	err := s.Repo.DB.Table("user_entity_memberships").
		Where("user_id = ?", userID).
		Distinct("entity_id").
		Pluck("entity_id", &entityIDs).Error
		
	if err != nil {
		log.Printf("Error querying user_entity_memberships for user %d: %v", userID, err)
		return []Entity{}, err
	}
	
	log.Printf("Found entity IDs through memberships for user %d: %v", userID, entityIDs)
	
	if len(entityIDs) == 0 {
		// No memberships found, return empty
		return []Entity{}, nil
	}
	
	// Get all entities user has membership to
	var entities []Entity
	err = s.Repo.DB.Where("id IN ?", entityIDs).Find(&entities).Error
	if err != nil {
		log.Printf("Error fetching entities by IDs %v: %v", entityIDs, err)
		return []Entity{}, err
	}
	
	log.Printf("Successfully fetched %d entities for user %d", len(entities), userID)
	return entities, nil
}

// CheckEntityExists - Utility method to check if entity exists
func (s *Service) CheckEntityExists(entityID uint) bool {
	var count int64
	err := s.Repo.DB.Model(&Entity{}).Where("id = ?", entityID).Count(&count).Error
	if err != nil {
		log.Printf("Error checking entity existence for ID %d: %v", entityID, err)
		return false
	}
	return count > 0
}

// GetUserAccessibleEntityIDs - Get all entity IDs a user can access - UPDATED
func (s *Service) GetUserAccessibleEntityIDs(userID uint, role string, directEntityID *uint) ([]uint, error) {
	var entityIDs []uint
	entityIDMap := make(map[uint]bool) // To avoid duplicates
	
	switch role {
	case "superadmin":
		// Get all entity IDs
		err := s.Repo.DB.Model(&Entity{}).Pluck("id", &entityIDs).Error
		return entityIDs, err
		
	case "templeadmin":
		if directEntityID != nil {
			entityIDs = append(entityIDs, *directEntityID)
			entityIDMap[*directEntityID] = true
		}
		// Also add entities created by them
		var createdEntityIDs []uint
		err := s.Repo.DB.Model(&Entity{}).Where("created_by = ?", userID).Pluck("id", &createdEntityIDs).Error
		if err == nil {
			for _, id := range createdEntityIDs {
				if !entityIDMap[id] {
					entityIDs = append(entityIDs, id)
					entityIDMap[id] = true
				}
			}
		}
		return entityIDs, err
		
	case "standarduser", "monitoringuser":
		if directEntityID != nil {
			entityIDs = append(entityIDs, *directEntityID)
			entityIDMap[*directEntityID] = true
		}
		
		// Add entities created by them
		var createdEntityIDs []uint
		err := s.Repo.DB.Model(&Entity{}).Where("created_by = ?", userID).Pluck("id", &createdEntityIDs).Error
		if err == nil {
			for _, id := range createdEntityIDs {
				if !entityIDMap[id] {
					entityIDs = append(entityIDs, id)
					entityIDMap[id] = true
				}
			}
		}
		
		// Add entities from memberships
		var membershipEntityIDs []uint
		err2 := s.Repo.DB.Table("user_entity_memberships").
			Where("user_id = ?", userID).
			Pluck("entity_id", &membershipEntityIDs).Error
		if err2 == nil {
			for _, id := range membershipEntityIDs {
				if !entityIDMap[id] {
					entityIDs = append(entityIDs, id)
					entityIDMap[id] = true
				}
			}
		}
		
		// Return the first non-nil error, or nil if both succeeded
		if err != nil {
			return entityIDs, err
		}
		return entityIDs, err2
		
	default:
		return []uint{}, errors.New("unsupported role")
	}
}

// CheckUserHasEntityAccess - Check if user has access to specific entity - UPDATED FOR CROSS-TENANT ACCESS
func (s *Service) CheckUserHasEntityAccess(userID uint, role string, entityID uint, directEntityID *uint, assignedEntityID *uint) bool {
	log.Printf("CheckUserHasEntityAccess - UserID: %d, Role: %s, EntityID: %d, DirectEntityID: %v, AssignedEntityID: %v", 
		userID, role, entityID, directEntityID, assignedEntityID)

	switch role {
	case "superadmin":
		return true
		
	case "templeadmin":
		// Check if it's their direct entity
		if directEntityID != nil && *directEntityID == entityID {
			log.Printf("Access granted: templeadmin %d has direct access to entity %d", userID, entityID)
			return true
		}
		
		// Check if they created this entity
		var count int64
		err := s.Repo.DB.Model(&Entity{}).Where("id = ? AND created_by = ?", entityID, userID).Count(&count).Error
		if err == nil && count > 0 {
			log.Printf("Access granted: templeadmin %d created entity %d", userID, entityID)
			return true
		}
		
		// NEW: Check if entity was created by users within their tenant context
		if directEntityID != nil {
			hasAccess := s.checkEntityCreatedByTenantUser(entityID, *directEntityID)
			if hasAccess {
				log.Printf("Access granted: entity %d was created by user within templeadmin %d's tenant context", entityID, userID)
				return true
			}
		}
		
		log.Printf("Access denied: templeadmin %d has no access to entity %d", userID, entityID)
		return false
		
	case "standarduser", "monitoringuser":
		// Check assigned entity (from tenant header)
		if assignedEntityID != nil && *assignedEntityID == entityID {
			log.Printf("Access granted: %s %d has assigned access to entity %d", role, userID, entityID)
			return true
		}
		
		// Check direct entity
		if directEntityID != nil && *directEntityID == entityID {
			log.Printf("Access granted: %s %d has direct access to entity %d", role, userID, entityID)
			return true
		}
		
		// Check if they created this entity
		var count int64
		err := s.Repo.DB.Model(&Entity{}).Where("id = ? AND created_by = ?", entityID, userID).Count(&count).Error
		if err == nil && count > 0 {
			log.Printf("Access granted: %s %d created entity %d", role, userID, entityID)
			return true
		}
		
		// Check membership
		err = s.Repo.DB.Table("user_entity_memberships").
			Where("user_id = ? AND entity_id = ?", userID, entityID).Count(&count).Error
		if err == nil && count > 0 {
			log.Printf("Access granted: %s %d has membership to entity %d", role, userID, entityID)
			return true
		}
		
		// NEW: Check if entity was created by users within their tenant context
		accessibleEntityIDs := s.getAccessibleEntityIDs(userID, role, directEntityID, assignedEntityID)
		for _, accessibleEntityID := range accessibleEntityIDs {
			hasAccess := s.checkEntityCreatedByTenantUser(entityID, accessibleEntityID)
			if hasAccess {
				log.Printf("Access granted: entity %d was created by user within %s %d's tenant context", entityID, role, userID)
				return true
			}
		}
		
		log.Printf("Access denied: %s %d has no access to entity %d", role, userID, entityID)
		return false
		
	default:
		log.Printf("Access denied: unsupported role %s for user %d", role, userID)
		return false
	}
}

// NEW: Check if an entity was created by a user who has access to the given tenant entity
func (s *Service) checkEntityCreatedByTenantUser(entityID, tenantEntityID uint) bool {
	// Get the creator of the entity
	var creatorID uint
	err := s.Repo.DB.Model(&Entity{}).Where("id = ?", entityID).Pluck("created_by", &creatorID).Error
	if err != nil {
		log.Printf("Error getting creator for entity %d: %v", entityID, err)
		return false
	}
	
	// Check if the creator has access to the tenant entity
	// 1. Check if creator has direct access (entity_id in users table)
	var count int64
	err = s.Repo.DB.Table("users").Where("id = ? AND entity_id = ?", creatorID, tenantEntityID).Count(&count).Error
	if err == nil && count > 0 {
		log.Printf("Entity %d creator %d has direct access to tenant entity %d", entityID, creatorID, tenantEntityID)
		return true
	}
	
	// 2. Check if creator has membership to the tenant entity
	err = s.Repo.DB.Table("user_entity_memberships").
		Where("user_id = ? AND entity_id = ?", creatorID, tenantEntityID).Count(&count).Error
	if err == nil && count > 0 {
		log.Printf("Entity %d creator %d has membership to tenant entity %d", entityID, creatorID, tenantEntityID)
		return true
	}
	
	log.Printf("Entity %d creator %d has no access to tenant entity %d", entityID, creatorID, tenantEntityID)
	return false
}

// Temple Admin → Update own temple
func (s *Service) UpdateEntity(e Entity, userID uint, ip string) error {
	// Check if entity exists first
	if !s.CheckEntityExists(e.ID) {
		auditDetails := map[string]interface{}{
			"temple_id": e.ID,
			"error":     "Temple not found",
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATE_FAILED", auditDetails, ip, "failure")
		return ErrEntityNotFound
	}

	// Get existing entity for comparison
	existingEntity, err := s.Repo.GetEntityByID(int(e.ID))
	if err != nil {
		// LOG FAILED TEMPLE UPDATE ATTEMPT (NOT FOUND)
		auditDetails := map[string]interface{}{
			"temple_id": e.ID,
			"error":     "Temple not found",
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATE_FAILED", auditDetails, ip, "failure")
		
		return err
	}

	e.UpdatedAt = time.Now()
	
	if err := s.Repo.UpdateEntity(e); err != nil {
		// LOG FAILED TEMPLE UPDATE (DB ERROR)
		auditDetails := map[string]interface{}{
			"temple_id":   e.ID,
			"temple_name": e.Name,
			"error":       err.Error(),
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATE_FAILED", auditDetails, ip, "failure")
		
		return err
	}

	// LOG SUCCESSFUL TEMPLE UPDATE
	auditDetails := map[string]interface{}{
		"temple_id":         e.ID,
		"temple_name":       e.Name,
		"previous_name":     existingEntity.Name,
		"temple_type":       e.TempleType,
		"email":             e.Email,
		"phone":             e.Phone,
		"city":              e.City,
		"state":             e.State,
		"main_deity":        e.MainDeity,
		"description":       e.Description,
		"updated_fields":    getUpdatedFields(existingEntity, e),
	}
	s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATED", auditDetails, ip, "success")

	return nil
}

// Super Admin → Delete temple
func (s *Service) DeleteEntity(id int, userID uint, ip string) error {
	// Get existing entity for audit log
	existingEntity, err := s.Repo.GetEntityByID(id)
	if err != nil {
		// LOG FAILED TEMPLE DELETION ATTEMPT (NOT FOUND)
		auditDetails := map[string]interface{}{
			"temple_id": id,
			"error":     "Temple not found",
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETE_FAILED", auditDetails, ip, "failure")
		
		return err
	}

	if err := s.Repo.DeleteEntity(id); err != nil {
		// LOG FAILED TEMPLE DELETION (DB ERROR)
		auditDetails := map[string]interface{}{
			"temple_id":   id,
			"temple_name": existingEntity.Name,
			"error":       err.Error(),
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETE_FAILED", auditDetails, ip, "failure")
		
		return err
	}

	// LOG SUCCESSFUL TEMPLE DELETION
	auditDetails := map[string]interface{}{
		"temple_id":   id,
		"temple_name": existingEntity.Name,
		"temple_type": existingEntity.TempleType,
		"email":       existingEntity.Email,
		"city":        existingEntity.City,
		"state":       existingEntity.State,
	}
	entityID := uint(id)
	s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETED", auditDetails, ip, "success")

	return nil
}

// ========== DEVOTEE MANAGEMENT ==========

// Temple Admin → Get devotees for specific entity
func (s *Service) GetDevotees(entityID uint) ([]DevoteeDTO, error) {
	return s.Repo.GetDevoteesByEntityID(entityID)
}

// Temple Admin → Get devotee statistics for entity
func (s *Service) GetDevoteeStats(entityID uint) (DevoteeStats, error) {
	return s.Repo.GetDevoteeStats(entityID)
}

// DashboardSummary is the structured JSON response
type DashboardSummary struct {
	RegisteredDevotees struct {
		Total     int64 `json:"total"`
		ThisMonth int64 `json:"this_month"`
	} `json:"registered_devotees"`

	SevaBookings struct {
		Today     int64 `json:"today"`
		ThisMonth int64 `json:"this_month"`
	} `json:"seva_bookings"`

	MonthDonations struct {
		Amount        float64 `json:"amount"`
		PercentChange float64 `json:"percent_change"`
	} `json:"month_donations"`

	UpcomingEvents struct {
		Total     int64 `json:"total"`
		ThisWeek  int64 `json:"this_week"`
	} `json:"upcoming_events"`
}

// Temple Admin → Dashboard Summary
func (s *Service) GetDashboardSummary(entityID uint) (DashboardSummary, error) {
	var summary DashboardSummary

	log.Printf("Generating dashboard summary for entity %d", entityID)

	// Check if entity exists first
	if !s.CheckEntityExists(entityID) {
		log.Printf("Entity %d not found for dashboard summary", entityID)
		return summary, ErrEntityNotFound
	}

	// ============= DEVOTEES =============
	totalDevotees, err := s.Repo.CountDevotees(entityID)
	if err != nil {
		log.Printf("Error counting devotees for entity %d: %v", entityID, err)
		return summary, err
	}
	thisMonthDevotees, err := s.Repo.CountDevoteesThisMonth(entityID)
	if err != nil {
		log.Printf("Error counting devotees this month for entity %d: %v", entityID, err)
		return summary, err
	}
	summary.RegisteredDevotees.Total = totalDevotees
	summary.RegisteredDevotees.ThisMonth = thisMonthDevotees

	// ============= SEVA BOOKINGS =============
	todaySevas, err := s.Repo.CountSevaBookingsToday(entityID)
	if err != nil {
		log.Printf("Error counting seva bookings today for entity %d: %v", entityID, err)
		return summary, err
	}
	monthSevas, err := s.Repo.CountSevaBookingsThisMonth(entityID)
	if err != nil {
		log.Printf("Error counting seva bookings this month for entity %d: %v", entityID, err)
		return summary, err
	}
	summary.SevaBookings.Today = todaySevas
	summary.SevaBookings.ThisMonth = monthSevas

	// ============= DONATIONS =============
	monthDonationAmount, percentChange, err := s.Repo.GetMonthDonationsWithChange(entityID)
	if err != nil {
		log.Printf("Error getting donation data for entity %d: %v", entityID, err)
		return summary, err
	}
	summary.MonthDonations.Amount = monthDonationAmount
	summary.MonthDonations.PercentChange = percentChange

	// ============= UPCOMING EVENTS =============
	totalUpcoming, err := s.Repo.CountUpcomingEvents(entityID)
	if err != nil {
		log.Printf("Error counting upcoming events for entity %d: %v", entityID, err)
		return summary, err
	}
	thisWeekUpcoming, err := s.Repo.CountUpcomingEventsThisWeek(entityID)
	if err != nil {
		log.Printf("Error counting upcoming events this week for entity %d: %v", entityID, err)
		return summary, err
	}
	summary.UpcomingEvents.Total = totalUpcoming
	summary.UpcomingEvents.ThisWeek = thisWeekUpcoming

	log.Printf("Dashboard summary generated successfully for entity %d", entityID)
	return summary, nil
}

// Helper function to track what fields were updated
func getUpdatedFields(old, new Entity) []string {
	var updatedFields []string

	if old.Name != new.Name {
		updatedFields = append(updatedFields, "name")
	}
	if (old.MainDeity == nil && new.MainDeity != nil) || 
	   (old.MainDeity != nil && new.MainDeity == nil) ||
	   (old.MainDeity != nil && new.MainDeity != nil && *old.MainDeity != *new.MainDeity) {
		updatedFields = append(updatedFields, "main_deity")
	}
	if old.TempleType != new.TempleType {
		updatedFields = append(updatedFields, "temple_type")
	}
	if old.Email != new.Email {
		updatedFields = append(updatedFields, "email")
	}
	if old.Phone != new.Phone {
		updatedFields = append(updatedFields, "phone")
	}
	if old.Description != new.Description {
		updatedFields = append(updatedFields, "description")
	}
	if old.StreetAddress != new.StreetAddress {
		updatedFields = append(updatedFields, "street_address")
	}
	if old.City != new.City {
		updatedFields = append(updatedFields, "city")
	}
	if old.State != new.State {
		updatedFields = append(updatedFields, "state")
	}
	if old.District != new.District {
		updatedFields = append(updatedFields, "district")
	}
	if old.Pincode != new.Pincode {
		updatedFields = append(updatedFields, "pincode")
	}

	return updatedFields
}