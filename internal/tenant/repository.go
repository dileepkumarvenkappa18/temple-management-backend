package tenant

import (
    "errors"
    "gorm.io/gorm"
    "log"
)

// Repository handles database operations
type Repository struct {
    db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
    // Debug check for table existence
    var tableExists bool
    db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tenant_user_assignments')").Scan(&tableExists)
    log.Printf("tenant_user_assignments table exists: %v", tableExists)
    
    // Debug count of records
    var count int64
    db.Table("tenant_user_assignments").Count(&count)
    log.Printf("Found %d records in tenant_user_assignments table", count)
    
    return &Repository{db: db}
}

// GetUserByEmail fetches a user by email
func (r *Repository) GetUserByEmail(email string) (*User, error) {
    var user User
    err := r.db.Where("email = ?", email).First(&user).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // User not found, but not an error
    }
    return &user, err
}

// GetTenantUsers fetches all users assigned to a tenant
func (r *Repository) GetTenantUsers(tenantID uint, role string) ([]UserResponse, error) {
    log.Printf("REPOSITORY: Fetching users for tenant ID: %d", tenantID)
    var userResponses []UserResponse
    
    // Build the query
    query := r.db.Table("users u").
        Select("u.id, u.full_name as name, u.email, u.phone, u.status, u.created_at").
        Joins("JOIN tenant_user_assignments tua ON u.id = tua.user_id").
        Where("tua.tenant_id = ? AND tua.status = 'active'", tenantID)
    
    // Execute the query
    err := query.Scan(&userResponses).Error
    if err != nil {
        log.Printf("Error fetching tenant users: %v", err)
        return nil, err
    }
    
    // Initialize empty array if nil
    if userResponses == nil {
        log.Printf("No users found, returning empty array")
        userResponses = []UserResponse{}
    } else {
        log.Printf("Found %d users for tenant ID %d", len(userResponses), tenantID)
    }
    
    return userResponses, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(user *User) error {
    log.Printf("Creating new user: %s (%s)", user.FullName, user.Email)
    return r.db.Create(user).Error
}

// UpdateTenantUserAssignment updates an existing tenant-user assignment or creates a new one
func (r *Repository) UpdateTenantUserAssignment(userID, tenantID, createdBy uint) error {
    log.Printf("🚨 REPOSITORY - Received tenant ID: %d for user ID: %d", tenantID, userID)
    
    // First check if record exists
    var count int64
    r.db.Model(&TenantUserAssignment{}).Where("user_id = ? AND tenant_id = ?", userID, tenantID).Count(&count)
    
    if count == 0 {
        // Create new assignment using GORM with explicit values
        assignment := TenantUserAssignment{
            UserID:    userID,
            TenantID:  tenantID,  // Explicitly set the tenant ID
            CreatedBy: createdBy,
            Status:    "active",
        }
        
        log.Printf("Creating new assignment with tenant_id=%d", tenantID)
        if err := r.db.Create(&assignment).Error; err != nil {
            log.Printf("Error creating assignment: %v", err)
            return err
        }
        
        // Verify the assignment was created with correct tenant_id
        var result struct {
            TenantID uint
        }
        err := r.db.Table("tenant_user_assignments").
            Select("tenant_id").
            Where("user_id = ? AND created_by = ?", userID, createdBy).
            Order("created_at DESC").
            Limit(1).
            Scan(&result).Error
            
        if err != nil {
            log.Printf("Warning: Could not verify assignment: %v", err)
        } else {
            log.Printf("Verified assignment - tenant_id in database: %d", result.TenantID)
            if result.TenantID != tenantID {
                log.Printf("WARNING: Expected tenant_id %d but found %d", tenantID, result.TenantID)
            }
        }
        
        return nil
    } else {
        // Update existing record
        log.Printf("Record exists, updating status for tenant_id=%d", tenantID)
        return r.db.Model(&TenantUserAssignment{}).
            Where("user_id = ? AND tenant_id = ?", userID, tenantID).
            Update("status", "active").
            Update("updated_at", gorm.Expr("NOW()")).
            Error
    }
}

// GetRoleIDByName gets role ID by name
func (r *Repository) GetRoleIDByName(roleName string) (uint, error) {
    // Convert frontend roles to database role names
    var dbRoleName string
    switch roleName {
    case "StandardUser":
        dbRoleName = "standarduser"
    case "MonitoringUser":
        dbRoleName = "monitoringuser"
    default:
        dbRoleName = roleName
    }
    
    // Mapping of role names to IDs based on your database
    roleIDs := map[string]uint{
        "superadmin":     1,
        "templeadmin":    2,
        "devotee":        3,
        "volunteer":      4,
        "standarduser":   5,
        "monitoringuser": 6,
    }
    
    log.Printf("Looking up role ID for '%s' (converted to '%s')", roleName, dbRoleName)
    
    // Check if role exists in our mapping
    if roleID, exists := roleIDs[dbRoleName]; exists {
        log.Printf("Found role ID %d for '%s'", roleID, dbRoleName)
        return roleID, nil
    }
    
    // If not in our mapping, try database lookup
    var role struct {
        ID uint
    }
    
    err := r.db.Table("user_roles").
        Select("id").
        Where("role_name = ?", dbRoleName).
        First(&role).Error
    
    if err != nil {
        log.Printf("Role '%s' not found in DB", dbRoleName)
        return 0, errors.New("invalid role name")
    }
    
    log.Printf("Found role ID %d for '%s' from database", role.ID, dbRoleName)
    return role.ID, nil
}