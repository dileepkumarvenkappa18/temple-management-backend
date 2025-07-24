package userprofile

import (
	"errors"
	"strings"

	"github.com/sharath018/temple-management-backend/internal/entity"
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

	// Temple Search
	SearchTemples(query string, state string, templeType string) ([]entity.Entity, error)
	ListPreApprovedTemples(limit int) ([]entity.Entity, error)
	GetTempleByID(entityID uint) (*entity.Entity, error)
	GetFullTempleByID(entityID uint) (*entity.Entity, error)
	FetchRecentTemples() ([]entity.Entity, error)

}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==============================
// 🔹 DevoteeProfile Operations
// ==============================

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

// ==============================
// 🔹 Membership Operations
// ==============================

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

// ==============================
// 🔹 Temple Search Operations
// ==============================

func (r *repository) SearchTemples(query string, state string, templeType string) ([]entity.Entity, error) {
	var temples []entity.Entity

	db := r.db.Model(&entity.Entity{}).
		Where("LOWER(status) = ?", "approved")

	if query != "" {
		q := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(city) LIKE ? OR LOWER(state) LIKE ?", q, q, q)
	}

	if state != "" {
		db = db.Where("LOWER(state) = ?", strings.ToLower(state))
	}

	if templeType != "" {
		db = db.Where("LOWER(temple_type) = ?", strings.ToLower(templeType))
	}

	err := db.Find(&temples).Error
	return temples, err
}

func (r *repository) ListPreApprovedTemples(limit int) ([]entity.Entity, error) {
	var temples []entity.Entity
	err := r.db.Model(&entity.Entity{}).
		Where("LOWER(status) = ?", "approved").
		Limit(limit).
		Find(&temples).Error
	return temples, err
}

func (r *repository) GetTempleByID(entityID uint) (*entity.Entity, error) {
	var temple entity.Entity
	err := r.db.Where("id = ?", entityID).First(&temple).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &temple, err
}

func (r *repository) GetFullTempleByID(entityID uint) (*entity.Entity, error) {
	var temple entity.Entity
	err := r.db.Where("id = ?", entityID).First(&temple).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &temple, err
}

func (r *repository) FetchRecentTemples() ([]entity.Entity, error) {
	var temples []entity.Entity
	err := r.db.Model(&entity.Entity{}).
		Where("LOWER(status) = ?", "approved").
		Order("created_at DESC").
		Limit(6).
		Find(&temples).Error
	return temples, err
}

