package auth

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(userID uint) (User, error) // 👈 ADD THIS
	FindRoleByName(name string) (*UserRole, error)
	FindEntityIDByUserID(userID uint) (*uint, error)
	CreateApprovalRequest(userID uint, requestType string) error
	UpdateEntityID(userID uint, entityID uint) error

}


type repository struct{ db *gorm.DB }

func (r *repository) UpdateEntityID(userID uint, entityID uint) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("entity_id", entityID).Error
}

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
	var entityID uint

	// First: Check for templeadmin approved request
	var req ApprovalRequest
	err := r.db.
		Where("user_id = ? AND status = ?", userID, "approved").
		Order("id DESC").
		First(&req).Error

	if err == nil && req.EntityID != nil {
		return req.EntityID, nil
	}

	// Second: Check for devotee/volunteer membership
	type membership struct {
		EntityID uint
	}

	var m membership
	err = r.db.
		Table("user_entity_memberships").
		Select("entity_id").
		Where("user_id = ?", userID).
		Order("joined_at DESC").
		First(&m).Error

	if err == nil {
		entityID = m.EntityID
		return &entityID, nil
	}

	return nil, gorm.ErrRecordNotFound
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
