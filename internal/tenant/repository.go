package tenant

import (
    "errors"
    "log"
    "gorm.io/gorm"
)

// Repository handles database operations for tenant users
type Repository struct {
    db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}
// In repository.go, replace GetTenantProfileByUserID with this improved version:

func (r *Repository) GetTenantProfileByUserID(userID uint) (*TenantProfileResponse, error) {
	log.Printf("🔍 REPO: Fetching tenant profile for user ID: %d", userID)

	var row tenantProfileRow

	// Strategy 1: Try via tenant_user_assignments (for standarduser, monitoringuser)
	err := r.db.Table("users").
		Select(`
			tenant_details.id AS tenant_id,
			tenant_details.temple_name,
			tenant_details.temple_place,
			tenant_details.temple_address,
			tenant_details.temple_phone_no,
			tenant_details.temple_description,
			tenant_details.logo_url,
			tenant_details.intro_video_url,
			users.id AS user_id,
			users.full_name AS user_full_name,
			users.email AS user_email,
			users.phone AS user_phone,
			user_roles.role_name AS user_role
		`).
		Joins("INNER JOIN tenant_user_assignments ON users.id = tenant_user_assignments.user_id").
		Joins("INNER JOIN tenant_details ON tenant_user_assignments.tenant_id = tenant_details.id").
		Joins("LEFT JOIN user_roles ON users.role_id = user_roles.id").
		Where("users.id = ?", userID).
		Where("tenant_user_assignments.status = ?", "active").
		Scan(&row).Error

	// ✅ Check if Strategy 1 succeeded
	if err == nil && row.TenantID > 0 {
		log.Printf("✅ REPO: Found tenant via assignment: %d", row.TenantID)
		return r.buildProfileResponse(&row), nil
	}

	// Log why Strategy 1 failed
	if err != nil {
		log.Printf("⚠️ REPO: Assignment query error: %v", err)
	} else {
		log.Printf("⚠️ REPO: No active assignment found for user %d", userID)
	}

	// Strategy 2: Try entity_id (for templeadmin)
	log.Printf("🔍 REPO: Trying entity_id lookup for user %d", userID)
	
	// Reset the row variable
	row = tenantProfileRow{}
	
	err = r.db.Table("users").
		Select(`
			tenant_details.id AS tenant_id,
			tenant_details.temple_name,
			tenant_details.temple_place,
			tenant_details.temple_address,
			tenant_details.temple_phone_no,
			tenant_details.temple_description,
			tenant_details.logo_url,
			tenant_details.intro_video_url,
			users.id AS user_id,
			users.full_name AS user_full_name,
			users.email AS user_email,
			users.phone AS user_phone,
			user_roles.role_name AS user_role
		`).
		Joins("INNER JOIN tenant_details ON users.entity_id = tenant_details.id").
		Joins("LEFT JOIN user_roles ON users.role_id = user_roles.id").
		Where("users.id = ?", userID).
		Where("users.entity_id IS NOT NULL").
		Where("users.entity_id > 0").
		Scan(&row).Error
	
	// ✅ Check if Strategy 2 succeeded
	if err == nil && row.TenantID > 0 {
		log.Printf("✅ REPO: Found tenant via entity_id: %d", row.TenantID)
		return r.buildProfileResponse(&row), nil
	}

	// Both strategies failed
	if err != nil {
		log.Printf("❌ REPO: Entity_id query error: %v", err)
		return nil, err
	}

	log.Printf("⚠️ REPO: No tenant found for user ID: %d (tried both assignments and entity_id)", userID)
	return nil, errors.New("no tenant assignment found")
}

// ✅ Helper method to build profile response (DRY principle)
func (r *Repository) buildProfileResponse(row *tenantProfileRow) *TenantProfileResponse {
	profile := &TenantProfileResponse{
		TenantID:          row.TenantID,
		TempleName:        row.TempleName,
		TemplePlace:       row.TemplePlace,
		TempleAddress:     row.TempleAddress,
		TemplePhoneNo:     row.TemplePhoneNo,
		TempleDescription: row.TempleDescription,
		LogoURL:           row.LogoURL,
		IntroVideoURL:     row.IntroVideoURL,
	}

	profile.User.ID = row.UserID
	profile.User.FullName = row.UserFullName
	profile.User.Email = row.UserEmail
	profile.User.Phone = row.UserPhone
	profile.User.Role = row.UserRole

	log.Printf("✅ REPO: Tenant profile built - Tenant ID: %d, User: %s", profile.TenantID, profile.User.FullName)
	return profile
}

// UpdateTenantProfile updates tenant and user information
func (r *Repository) UpdateTenantProfile(tenantID, userID uint, input UpdateTenantProfileRequest) error {
	log.Printf("🔄 REPO: Updating tenant profile - Tenant ID: %d, User ID: %d", tenantID, userID)

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update user info
	userUpdates := map[string]interface{}{}
	if input.FullName != "" {
		userUpdates["full_name"] = input.FullName
	}
	if input.Phone != "" {
		userUpdates["phone"] = input.Phone
	}

	if len(userUpdates) > 0 {
		if err := tx.Table("users").Where("id = ?", userID).Updates(userUpdates).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update tenant info
	tenantUpdates := map[string]interface{}{}
	if input.TempleName != "" {
		tenantUpdates["temple_name"] = input.TempleName
	}
	if input.TemplePlace != "" {
		tenantUpdates["temple_place"] = input.TemplePlace
	}
	if input.TempleAddress != "" {
		tenantUpdates["temple_address"] = input.TempleAddress
	}
	if input.TemplePhoneNo != "" {
		tenantUpdates["temple_phone_no"] = input.TemplePhoneNo
	}
	if input.TempleDescription != "" {
		tenantUpdates["temple_description"] = input.TempleDescription
	}
	if input.LogoURL != "" {
		tenantUpdates["logo_url"] = input.LogoURL
	}
	if input.IntroVideoURL != "" {
		tenantUpdates["intro_video_url"] = input.IntroVideoURL
	}

	if len(tenantUpdates) > 0 {
		if err := tx.Table("tenant_details").
			Where("id = ?", tenantID).
			Updates(tenantUpdates).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(userID uint) (*User, error) {
    var user User
    err := r.db.Where("id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(email string) (*User, error) {
    var user User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

// CheckUserBelongsToTenant checks if a user belongs to a tenant
func (r *Repository) CheckUserBelongsToTenant(userID, tenantID uint) (bool, error) {
    var count int64
    err := r.db.Table("tenant_user_assignments").
        Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, "active").
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// UpdateUserStatus updates user status in both users and tenant_user_assignments tables
func (r *Repository) UpdateUserStatus(userID, tenantID uint, status string) error {
    tx := r.db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    // Update in users table
    err := tx.Table("users").Where("id = ?", userID).Update("status", status).Error
    if err != nil {
        tx.Rollback()
        return err
    }

    // Update in tenant_user_assignments table
    err = tx.Table("tenant_user_assignments").
        Where("user_id = ? AND tenant_id = ?", userID, tenantID).
        Update("status", status).Error
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

// UpdateUserDetails updates user details
func (r *Repository) UpdateUserDetails(userID uint, input UserInput) error {
    updates := make(map[string]interface{})
    
    if input.Name != "" {
        updates["full_name"] = input.Name
    }
    if input.Email != "" {
        updates["email"] = input.Email
    }
    if input.Phone != "" {
        updates["phone"] = input.Phone
    }
    
    if len(updates) == 0 {
        return nil
    }
    
    return r.db.Table("users").Where("id = ?", userID).Updates(updates).Error
}

func (r *Repository) GetTenantUsers(tenantID uint, role string) ([]UserResponse, error) {
    var users []UserResponse

    query := r.db.Table("users").
        Select(`
            users.id,
            users.full_name as name,
            users.email,
            users.phone,
            users.status,
            users.created_at
        `).
        Joins("INNER JOIN tenant_user_assignments tua ON users.id = tua.user_id").
        Where("tua.tenant_id = ?", tenantID).
        Where("tua.status = ?", "active")

    // ❌ role filter REMOVED from DB
    // Role is handled by JWT / frontend mapping

    if err := query.Scan(&users).Error; err != nil {
        return nil, err
    }

    return users, nil
}


// GetRoleIDByName retrieves role ID by role name
func (r *Repository) GetRoleIDByName(roleName string) (uint, error) {
    var roleID uint
    err := r.db.Table("roles").
        Select("id").
        Where("role_name = ?", roleName).
        Scan(&roleID).Error
    if err != nil {
        return 0, err
    }
    return roleID, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(user *User) error {
    return r.db.Create(user).Error
}

// UpdateTenantUserAssignment updates or creates tenant user assignment
func (r *Repository) UpdateTenantUserAssignment(userID, tenantID, creatorID uint) error {
    // Check if assignment exists
    var existing TenantUserAssignment
    err := r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
        First(&existing).Error
    
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }
    
    if existing.ID > 0 {
        // Update existing assignment
        return r.db.Model(&existing).
            Updates(map[string]interface{}{
                "status": "active",
                "created_by": creatorID,
            }).Error
    }
    
    // Create new assignment
    assignment := TenantUserAssignment{
        UserID:    userID,
        TenantID:  tenantID,
        CreatedBy: creatorID,
        Status:    "active",
    }
    
    return r.db.Create(&assignment).Error
}