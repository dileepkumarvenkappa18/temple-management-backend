package auth

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(userID uint) (User, error) // 👈 ADD THIS
	FindRoleByName(name string) (*UserRole, error)
	FindEntityIDByUserID(userID uint) (*uint, error)
	CreateApprovalRequest(userID uint, requestType string) error
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

// ✅ FIXED: Get the approved EntityID even if entity_id is NULL (temporary login allowed)
func (r *repository) FindEntityIDByUserID(userID uint) (*uint, error) {
	var req ApprovalRequest
	err := r.db.
		Where("user_id = ? AND status = ?", userID, "approved").
		Order("id DESC").
		First(&req).Error

	if err != nil {
		return nil, err
	}

	// Return pointer to entity ID (may still be nil if not created yet)
	return req.EntityID, nil
}

// ✅ Used during templeadmin registration to create approval request
func (r *repository) CreateApprovalRequest(userID uint, requestType string) error {
	req := ApprovalRequest{
		UserID:      userID,
		RequestType: requestType,
		Status:      "pending",
	}
	return r.db.Create(&req).Error
}

func (r *repository) FindByID(userID uint) (User, error) {
	var user User
	err := r.db.Preload("Role").First(&user, userID).Error
	return user, err
}
