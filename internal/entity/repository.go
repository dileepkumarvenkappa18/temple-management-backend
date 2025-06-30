package entity

import (
	"errors"
	"fmt"
	"time"

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


func (r *Repository) GetAllEntities() ([]Entity, error) {
	var entities []Entity
	err := r.DB.
		Preload("Address").
		Preload("Documents").
		Find(&entities).Error
	return entities, err
}

func (r *Repository) GetEntityByID(id int) (Entity, error) {
	var e Entity
	err := r.DB.
		Preload("Address").
		Preload("Documents").
		First(&e, id).Error
	return e, err
}

func (r *Repository) UpdateEntity(e Entity) error {
	return r.DB.Model(&Entity{}).Where("id = ?", e.ID).Updates(e).Error
}

func (r *Repository) ToggleEntityStatus(id int, isActive bool) error {
	return r.DB.Model(&Entity{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *Repository) DeleteEntity(id int) error {
	return r.DB.Delete(&Entity{}, id).Error
}

// ========== ADDRESS ==========

func (r *Repository) AddEntityAddress(addr EntityAddress) error {
	return r.DB.Create(&addr).Error
}

func (r *Repository) GetEntityAddress(entityID int) (EntityAddress, error) {
	var addr EntityAddress
	result := r.DB.Where("entity_id = ?", entityID).First(&addr)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return addr, fmt.Errorf("no address found for entity_id %d", entityID)
		}
		return addr, result.Error
	}

	return addr, nil
}

// ========== DOCUMENTS ==========

func (r *Repository) AddEntityDocument(doc EntityDocument) error {
	doc.UploadedAt = time.Now()
	return r.DB.Create(&doc).Error
}

func (r *Repository) GetEntityDocuments(entityID int) ([]EntityDocument, error) {
	var docs []EntityDocument
	err := r.DB.Where("entity_id = ?", entityID).Find(&docs).Error
	return docs, err
}