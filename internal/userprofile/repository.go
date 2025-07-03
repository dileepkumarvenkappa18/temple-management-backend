package userprofile

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	// DevoteeProfile methods
	Create(profile *DevoteeProfile) error
	GetByUserID(userID uint) (*DevoteeProfile, error)
	Update(profile *DevoteeProfile) error

	// Membership methods
	CreateMembership(m *UserEntityMembership) error
	GetMembership(userID, entityID uint) (*UserEntityMembership, error)
	ListMembershipsByUser(userID uint) ([]UserEntityMembership, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ========== DevoteeProfile ==========

func (r *repository) Create(profile *DevoteeProfile) error {
	return r.db.Create(profile).Error
}

func (r *repository) GetByUserID(userID uint) (*DevoteeProfile, error) {
	var profile DevoteeProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *repository) Update(profile *DevoteeProfile) error {
	return r.db.Save(profile).Error
}

// ========== Membership ==========

func (r *repository) CreateMembership(m *UserEntityMembership) error {
	return r.db.Create(m).Error
}

func (r *repository) GetMembership(userID, entityID uint) (*UserEntityMembership, error) {
	var membership UserEntityMembership
	err := r.db.Where("user_id = ? AND entity_id = ?", userID, entityID).First(&membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &membership, err
}

func (r *repository) ListMembershipsByUser(userID uint) ([]UserEntityMembership, error) {
	var memberships []UserEntityMembership
	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	return memberships, err
}
