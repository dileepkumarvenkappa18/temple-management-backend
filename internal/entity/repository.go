package entity

import (
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"gorm.io/gorm"
)

// Repository defines the interface for temple entity operations.
type Repository interface {
	CreateEntity(e *Entity) error
	CreateApprovalRequest(req *auth.ApprovalRequest) error

	// Temple listing
	GetAllEntities() ([]Entity, error)                  // Super Admin
	GetEntitiesByCreator(userID uint) ([]Entity, error) // Temple Admin
	GetAllEntitiesByUser(userID int) ([]Entity, error)  // Alias for compatibility

	GetEntityByID(id int) (Entity, error)
	UpdateEntity(e Entity) error
	DeleteEntity(id int) error

	// Dashboard stats
	GetTotalDonationsByEntity(entityID int) (float64, error)
	GetRecentDonationsCount(entityID int, since time.Time) (int64, error)
	GetRecentDonors(entityID int, limit int) ([]string, error)
}

// GORM implementation of the Repository interface
type repository struct {
	DB *gorm.DB
}

// Constructor
func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

// CreateEntity inserts a new temple (entity) into the database
func (r *repository) CreateEntity(e *Entity) error {
	return r.DB.Create(e).Error
}

// CreateApprovalRequest stores a new approval request
func (r *repository) CreateApprovalRequest(req *auth.ApprovalRequest) error {
	return r.DB.Create(req).Error
}

// GetAllEntitiesByUser forwards to GetEntitiesByCreator for compatibility
func (r *repository) GetAllEntitiesByUser(userID int) ([]Entity, error) {
	return r.GetEntitiesByCreator(uint(userID))
}

// GetEntitiesByCreator fetches temples created by a specific user
func (r *repository) GetEntitiesByCreator(userID uint) ([]Entity, error) {
	var entities []Entity
	err := r.DB.
		Where("created_by = ?", userID).
		Order("created_at DESC").
		Find(&entities).Error
	return entities, err
}

// GetAllEntities returns all temples (only for superadmin)
func (r *repository) GetAllEntities() ([]Entity, error) {
	var entities []Entity
	err := r.DB.
		Order("created_at DESC").
		Find(&entities).Error
	return entities, err
}

// GetEntityByID retrieves a temple by ID
func (r *repository) GetEntityByID(id int) (Entity, error) {
	var entity Entity
	err := r.DB.First(&entity, id).Error
	return entity, err
}

// UpdateEntity modifies an existing temple
func (r *repository) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	return r.DB.Model(&Entity{}).
		Where("id = ?", e.ID).
		Updates(e).Error
}

// DeleteEntity removes a temple from the database
func (r *repository) DeleteEntity(id int) error {
	return r.DB.Delete(&Entity{}, id).Error
}

// GetTotalDonationsByEntity calculates the total amount donated to a temple
func (r *repository) GetTotalDonationsByEntity(entityID int) (float64, error) {
	var total float64
	err := r.DB.
		Table("donations").
		Where("entity_id = ?", entityID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetRecentDonationsCount returns count of donations since a certain time
func (r *repository) GetRecentDonationsCount(entityID int, since time.Time) (int64, error) {
	var count int64
	err := r.DB.
		Table("donations").
		Where("entity_id = ? AND created_at >= ?", entityID, since).
		Count(&count).Error
	return count, err
}

// GetRecentDonors returns a list of recent donor names
func (r *repository) GetRecentDonors(entityID int, limit int) ([]string, error) {
	var names []string
	err := r.DB.
		Table("donations").
		Where("entity_id = ?", entityID).
		Order("created_at DESC").
		Limit(limit).
		Pluck("donor_name", &names).Error
	return names, err
}