package entity

import (
	"database/sql"
	"time"
	"math"
	"fmt"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// ========== ENTITY CORE ==========

// Create a new temple entity
func (r *Repository) CreateEntity(e *Entity) error {
	return r.DB.Create(e).Error
}

// Get tenant ID for a user from tenant_user_assignments table
func (r *Repository) GetTenantIDForUser(userID uint) (uint, error) {
	var tenantID uint
	
	err := r.DB.Table("tenant_user_assignments").
		Select("tenant_id").
		Where("user_id = ? AND status = ?", userID, "active").
		Limit(1).
		Scan(&tenantID).
		Error
		
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	
	return tenantID, nil
}

// Get user's role ID
func (r *Repository) GetUserRoleID(userID uint) (uint, error) {
	var roleID uint
	
	err := r.DB.Table("users").
		Select("role_id").
		Where("id = ?", userID).
		Limit(1).
		Scan(&roleID).
		Error
		
	if err != nil {
		return 0, err
	}
	
	return roleID, nil
}

// Create an approval request for the temple
func (r *Repository) CreateApprovalRequest(req *auth.ApprovalRequest) error {
	return r.DB.Create(req).Error
}

// Fetch all temple entities (ordered by most recent)
func (r *Repository) GetAllEntities() ([]Entity, error) {
	var entities []Entity
	err := r.DB.Order("created_at DESC").Find(&entities).Error
	return entities, err
}

// Get entities with creator role information
func (r *Repository) GetEntitiesWithRoleInfo() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := r.DB.Table("entities e").
		Select(`e.*, 
				ur.role_name as creator_role_name,
				CASE WHEN ur.id = 1 THEN true ELSE false END as is_auto_approved`).
		Joins("LEFT JOIN user_roles ur ON ur.id = e.creator_role_id").
		Order("e.created_at DESC").
		Find(&results).Error
		
	return results, err
}

// Fetch entities created by a specific user
func (r *Repository) GetEntitiesByCreator(creatorID uint) ([]Entity, error) {
	var entities []Entity
	err := r.DB.Where("created_by = ?", creatorID).Order("created_at DESC").Find(&entities).Error
	return entities, err
}

// Get approval statistics by role
func (r *Repository) GetApprovalStatsByRole() (map[string]interface{}, error) {
	type RoleStats struct {
		RoleID        uint   `json:"role_id"`
		RoleName      string `json:"role_name"`
		TotalTemples  int64  `json:"total_temples"`
		AutoApproved  int64  `json:"auto_approved"`
		PendingCount  int64  `json:"pending_count"`
		ApprovedCount int64  `json:"approved_count"`
		RejectedCount int64  `json:"rejected_count"`
	}
	
	var stats []RoleStats
	
	err := r.DB.Table("entities e").
		Select(`e.creator_role_id as role_id,
				ur.role_name,
				COUNT(*) as total_temples,
				COUNT(CASE WHEN e.status = 'approved' AND e.creator_role_id = 1 THEN 1 END) as auto_approved,
				COUNT(CASE WHEN e.status = 'pending' THEN 1 END) as pending_count,
				COUNT(CASE WHEN e.status = 'approved' THEN 1 END) as approved_count,
				COUNT(CASE WHEN e.status = 'rejected' THEN 1 END) as rejected_count`).
		Joins("LEFT JOIN user_roles ur ON ur.id = e.creator_role_id").
		Group("e.creator_role_id, ur.role_name").
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"role_statistics": stats,
	}, nil
}

// Fetch a single temple entity by ID
func (r *Repository) GetEntityByID(id int) (Entity, error) {
	var entity Entity
	err := r.DB.First(&entity, id).Error
	return entity, err
}

// Update an existing temple entity
func (r *Repository) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	
	updates := map[string]interface{}{
		"name":                    e.Name,
		"main_deity":              e.MainDeity,
		"temple_type":             e.TempleType,
		"established_year":        e.EstablishedYear,
		"email":                   e.Email,
		"phone":                   e.Phone,
		"description":             e.Description,
		"street_address":          e.StreetAddress,
		"landmark":                e.Landmark,
		"city":                    e.City,
		"district":                e.District,
		"state":                   e.State,
		"pincode":                 e.Pincode,
		"map_link":                e.MapLink,
		"registration_cert_url":   e.RegistrationCertURL,
		"registration_cert_info":  e.RegistrationCertInfo,
		"trust_deed_url":          e.TrustDeedURL,
		"trust_deed_info":         e.TrustDeedInfo,
		"property_docs_url":       e.PropertyDocsURL,
		"property_docs_info":      e.PropertyDocsInfo,
		"additional_docs_urls":    e.AdditionalDocsURLs,
		"additional_docs_info":    e.AdditionalDocsInfo,
		"updated_at":              e.UpdatedAt,
	}
	
	return r.DB.Model(&Entity{}).Where("id = ?", e.ID).Updates(updates).Error
}
// UpdateEntityStatus updates only the IsActive field of an entity
func (r *Repository) UpdateEntityStatus(id int, isActive bool) error {
    updates := map[string]interface{}{
        "isactive":   isActive,
        "updated_at": time.Now(),
    }
    
    result := r.DB.Model(&Entity{}).Where("id = ?", id).Updates(updates)
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("no entity found with id %d", id)
    }
    
    return nil
}

// GetActiveEntities retrieves only active entities
func (r *Repository) GetActiveEntities() ([]Entity, error) {
    var entities []Entity
    err := r.DB.Where("isactive = ?", true).Order("created_at DESC").Find(&entities).Error
    return entities, err
}

// GetActiveEntitiesByCreator retrieves only active entities created by a specific user
func (r *Repository) GetActiveEntitiesByCreator(creatorID uint) ([]Entity, error) {
    var entities []Entity
    err := r.DB.Where("created_by = ? AND isactive = ?", creatorID, true).
        Order("created_at DESC").
        Find(&entities).Error
    return entities, err
}
// Delete a temple entity by ID
func (r *Repository) DeleteEntity(id int) error {
	return r.DB.Delete(&Entity{}, id).Error
}

// ========== DEVOTEE MANAGEMENT ==========

type DevoteeDTO struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

func (r *Repository) GetDevoteesByEntityID(entityID uint) ([]DevoteeDTO, error) {
	var devotees []DevoteeDTO

	err := r.DB.
		Table("user_entity_memberships AS uem").
		Select("u.id AS user_id, u.full_name, u.email, u.phone, uem.status").
		Joins("JOIN users u ON u.id = uem.user_id").
		Joins("JOIN user_roles ur ON u.role_id = ur.id").
		Where("uem.entity_id = ? AND ur.role_name = ?", entityID, "devotee").
		Scan(&devotees).Error

	return devotees, err
}

type DevoteeStats struct {
	TotalDevotees  int64 `json:"total_devotees"`
	ActiveDevotees int64 `json:"active_devotees"`
	NewThisMonth   int64 `json:"new_this_month"`
}

func (r *Repository) GetDevoteeStats(entityID uint) (DevoteeStats, error) {
	var stats DevoteeStats

	err := r.DB.Table("user_entity_memberships").
		Joins("JOIN users ON users.id = user_entity_memberships.user_id").
		Joins("JOIN user_roles ON user_roles.id = users.role_id").
		Where("user_entity_memberships.entity_id = ? AND user_roles.role_name = ?", entityID, "devotee").
		Count(&stats.TotalDevotees).Error
	if err != nil {
		return stats, err
	}

	err = r.DB.Table("user_entity_memberships").
		Joins("JOIN users ON users.id = user_entity_memberships.user_id").
		Joins("JOIN user_roles ON user_roles.id = users.role_id").
		Where("user_entity_memberships.entity_id = ? AND user_roles.role_name = ? AND user_entity_memberships.status = ?", entityID, "devotee", "active").
		Count(&stats.ActiveDevotees).Error
	if err != nil {
		return stats, err
	}

	err = r.DB.Table("user_entity_memberships").
		Joins("JOIN users ON users.id = user_entity_memberships.user_id").
		Joins("JOIN user_roles ON user_roles.id = users.role_id").
		Where("user_entity_memberships.entity_id = ? AND user_roles.role_name = ? AND user_entity_memberships.created_at >= DATE_TRUNC('month', NOW())", entityID, "devotee").
		Count(&stats.NewThisMonth).Error
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func (r *Repository) CountDevotees(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("user_entity_memberships AS uem").
		Joins("JOIN user_roles ur ON ur.id = (SELECT role_id FROM users WHERE users.id = uem.user_id)").
		Where("uem.entity_id = ? AND ur.role_name = ?", entityID, "devotee").
		Count(&count).Error
	return count, err
}

func (r *Repository) CountDevoteesThisMonth(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("user_entity_memberships AS uem").
		Joins("JOIN user_roles ur ON ur.id = (SELECT role_id FROM users WHERE users.id = uem.user_id)").
		Where("uem.entity_id = ? AND ur.role_name = ? AND uem.created_at >= DATE_TRUNC('month', NOW())", entityID, "devotee").
		Count(&count).Error
	return count, err
}

func (r *Repository) CountSevaBookingsToday(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("seva_bookings").
		Where("entity_id = ? AND DATE(booking_time) = CURRENT_DATE", entityID).
		Count(&count).Error
	return count, err
}

func (r *Repository) CountSevaBookingsThisMonth(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("seva_bookings").
		Where("entity_id = ? AND booking_time >= DATE_TRUNC('month', NOW())", entityID).
		Count(&count).Error
	return count, err
}

func (r *Repository) GetMonthDonationsWithChange(entityID uint) (float64, float64, error) {
	var currentMonth, previousMonth float64

	err := r.DB.
		Table("donations").
		Select("COALESCE(SUM(amount), 0)").
		Where("entity_id = ? AND created_at >= DATE_TRUNC('month', NOW())", entityID).
		Scan(&currentMonth).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.DB.
		Table("donations").
		Select("COALESCE(SUM(amount), 0)").
		Where("entity_id = ? AND created_at >= DATE_TRUNC('month', NOW()) - INTERVAL '1 month' AND created_at < DATE_TRUNC('month', NOW())", entityID).
		Scan(&previousMonth).Error
	if err != nil {
		return 0, 0, err
	}

	var percentChange float64
	if previousMonth > 0 {
		percentChange = ((currentMonth - previousMonth) / previousMonth) * 100
	} else if currentMonth > 0 {
		percentChange = 100
	} else {
		percentChange = 0
	}

	return currentMonth, math.Round(percentChange*100) / 100, nil
}

func (r *Repository) CountUpcomingEvents(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("events").
		Where("entity_id = ? AND event_date >= CURRENT_DATE", entityID).
		Count(&count).Error
	return count, err
}

func (r *Repository) CountUpcomingEventsThisWeek(entityID uint) (int64, error) {
	var count int64
	err := r.DB.
		Table("events").
		Where(`entity_id = ? AND event_date >= DATE_TRUNC('week', CURRENT_DATE) AND event_date < DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '7 days'`, entityID).
		Count(&count).Error
	return count, err
}