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

func (r *Repository) CreateEntity(e *Entity) error {
	return r.DB.Create(e).Error
}

func (r *Repository) CreateApprovalRequest(req *auth.ApprovalRequest) error {
	return r.DB.Create(req).Error
}

func (r *Repository) GetAllEntities() ([]Entity, error) {
	var entities []Entity
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r *Repository) GetEntityByID(id int) (Entity, error) {
	var e Entity
	err := r.DB.First(&e, id).Error
	return e, err
}

func (r *Repository) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	return r.DB.Model(&Entity{}).Where("id = ?", e.ID).Updates(e).Error
}

func (r *Repository) DeleteEntity(id int) error {
	return r.DB.Delete(&Entity{}, id).Error
}
