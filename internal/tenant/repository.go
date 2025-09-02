package tenant

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Repository handles CRUD operations for TenantUser.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes the repository and migrates the TenantUser table.
func NewRepository(db *gorm.DB) *Repository {
	// Auto-migrate the TenantUser table
	if err := db.AutoMigrate(&TenantUser{}); err != nil {
		panic("Failed to migrate TenantUser table: " + err.Error())
	}
	return &Repository{db: db}
}

// GetUsers fetches users for a tenant with optional role or name filters
func (r *Repository) GetUsers(tenantID uint, role string, name string) ([]TenantUser, error) {
	var users []TenantUser
	query := r.db.Model(&TenantUser{}).Where("tenant_id = ?", tenantID)

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByEmail fetches a single user by email and tenant
func (r *Repository) GetByEmail(email string, tenantID uint) (*TenantUser, error) {
	var user TenantUser
	err := r.db.Where("email = ? AND tenant_id = ?", email, tenantID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // return nil if user not found
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create inserts a new TenantUser
func (r *Repository) Create(user *TenantUser) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.Create(user).Error
}

// Update modifies an existing TenantUser safely within the tenant
func (r *Repository) Update(user *TenantUser) error {
	user.UpdatedAt = time.Now()
	result := r.db.Model(&TenantUser{}).
		Where("id = ? AND tenant_id = ?", user.ID, user.TenantID).
		Updates(map[string]interface{}{
			"name":      user.Name,
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"password":  user.Password,
			"updated_at": user.UpdatedAt,
		})
	if result.RowsAffected == 0 {
		return errors.New("user not found or tenant mismatch")
	}
	return result.Error
}

// Delete performs a soft delete on a user within the tenant
func (r *Repository) Delete(userID, tenantID uint) error {
	result := r.db.Where("id = ? AND tenant_id = ?", userID, tenantID).Delete(&TenantUser{})
	if result.RowsAffected == 0 {
		return errors.New("user not found or tenant mismatch")
	}
	return result.Error
}
