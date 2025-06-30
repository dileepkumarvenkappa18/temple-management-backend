package auth

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindRoleByName(name string) (*UserRole, error)
	FindEntityIDByUserID(userID uint) (*uint, error)
	CreateApprovalRequest(userID uint, requestType string) error // ✅ NEW
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var u User
	err := r.db.Preload("Role").Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *repository) FindRoleByName(name string) (*UserRole, error) {
	var role UserRole
	err := r.db.Where("role_name = ?", name).First(&role).Error
	return &role, err
}

// ✅ NEW: Get the approved EntityID for a user (for tenant/devotee/volunteer)
func (r *repository) FindEntityIDByUserID(userID uint) (*uint, error) {
	var req ApprovalRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, "approved").First(&req).Error
	if err != nil {
		return nil, err
	}
	return req.EntityID, nil
}

// ✅ NEW: Used during templeadmin registration to create approval request
func (r *repository) CreateApprovalRequest(userID uint, requestType string) error {
	req := ApprovalRequest{
		UserID:      userID,
		RequestType: requestType, // should be "entity"
		Status:      "pending",
	}
	return r.db.Create(&req).Error
}
