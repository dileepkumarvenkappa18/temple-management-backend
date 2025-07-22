package entity

import (
	"time"

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

// Create an approval request for the temple (linked to auth module)
func (r *Repository) CreateApprovalRequest(req *auth.ApprovalRequest) error {
	return r.DB.Create(req).Error
}

// Fetch all temple entities (ordered by most recent)
func (r *Repository) GetAllEntities() ([]Entity, error) {
	var entities []Entity
	err := r.DB.Order("created_at DESC").Find(&entities).Error
	return entities, err
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
	return r.DB.Model(&Entity{}).Where("id = ?", e.ID).Updates(e).Error
}

// Delete a temple entity by ID
func (r *Repository) DeleteEntity(id int) error {
	return r.DB.Delete(&Entity{}, id).Error
}

