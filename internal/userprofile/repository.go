package userprofile

import (
	"errors"

	"gorm.io/gorm"
)

// Repository defines all methods related to Devotee profiles and memberships.
type Repository interface {
	// DevoteeProfile methods
	Create(profile *DevoteeProfile) error
	GetByUserID(userID uint) (*DevoteeProfile, error)
	Update(profile *DevoteeProfile) error

	// Membership methods
	CreateMembership(m *UserEntityMembership) error
	GetMembership(userID, entityID uint) (*UserEntityMembership, error)
	ListMembershipsByUser(userID uint) ([]UserEntityMembership, error)
	ListUserIDsByEntity(entityID uint) ([]uint, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==============================
// 🔹 DevoteeProfile Operations
// ==============================

// Create inserts a new DevoteeProfile into the database
func (r *repository) Create(profile *DevoteeProfile) error {
	return r.db.Create(profile).Error
}

// GetByUserID fetches a DevoteeProfile by the associated user ID
func (r *repository) GetByUserID(userID uint) (*DevoteeProfile, error) {
	var profile DevoteeProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update modifies an existing DevoteeProfile
func (r *repository) Update(profile *DevoteeProfile) error {
	return r.db.Save(profile).Error
}

// ==============================
// 🔹 Membership Operations
// ==============================

// CreateMembership inserts a new user-entity membership record
func (r *repository) CreateMembership(m *UserEntityMembership) error {
	return r.db.Create(m).Error
}

// GetMembership checks if a user has already joined a specific entity
func (r *repository) GetMembership(userID, entityID uint) (*UserEntityMembership, error) {
	var membership UserEntityMembership
	err := r.db.Where("user_id = ? AND entity_id = ?", userID, entityID).First(&membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &membership, err
}

// ListMembershipsByUser returns all temple memberships for a user
func (r *repository) ListMembershipsByUser(userID uint) ([]UserEntityMembership, error) {
	var memberships []UserEntityMembership
	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	return memberships, err
}

// ListUserIDsByEntity returns all user IDs that joined a given temple (entity)
func (r *repository) ListUserIDsByEntity(entityID uint) ([]uint, error) {
	var memberships []UserEntityMembership
	err := r.db.Select("user_id").Where("entity_id = ?", entityID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint, len(memberships))
	for i, m := range memberships {
		userIDs[i] = m.UserID
	}
	return userIDs, nil
}
