package superadmin

import (
	"context"
	"time"
	"errors"
	"fmt"
	

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"gorm.io/gorm"
	"database/sql"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// =========================== TENANT ===========================

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (auth.User, error) {
	var user auth.User
	err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		First(&user).Error
	return user, err
}

func (r *Repository) GetPendingTenants(ctx context.Context) ([]auth.User, error) {
	var tenants []auth.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ? AND users.status = ?", "templeadmin", "pending").
		Find(&tenants).Error
	return tenants, err
}

func (r *Repository) GetTenantsWithFilters(ctx context.Context, status string, limit, page int) ([]TenantWithDetails, int64, error) {
	var tenants []TenantWithDetails
	var total int64

	offset := (page - 1) * limit

	// Build the base query for counting
	countQuery := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ?", "templeadmin")

	if status != "" {
		countQuery = countQuery.Where("LOWER(users.status) = LOWER(?)", status)
	}

	// Get total count
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build the main query with LEFT JOIN to include temple details
	query := r.db.WithContext(ctx).
		Table("users").
		Select(`
			users.id,
			users.full_name,
			users.email,
			users.phone,
			users.role_id,
			users.status,
			users.created_at,
			users.updated_at,
			td.id as temple_id,
			td.temple_name,
			td.temple_place,
			td.temple_address,
			td.temple_phone_no,
			td.temple_description,
			td.created_at as temple_created_at,
			td.updated_at as temple_updated_at
		`).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Joins("LEFT JOIN tenant_details td ON users.id = td.user_id").
		Where("user_roles.role_name = ?", "templeadmin")

	if status != "" {
		query = query.Where("LOWER(users.status) = LOWER(?)", status)
	}

	// Execute query with pagination
	rows, err := query.Limit(limit).Offset(offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results into our custom struct
	for rows.Next() {
		var tenant TenantWithDetails
		var templeID *uint
		var templeName, templePlace, templeAddress, templePhoneNo, templeDescription *string
		var templeCreatedAt, templeUpdatedAt *time.Time

		err := rows.Scan(
			&tenant.ID,
			&tenant.FullName,
			&tenant.Email,
			&tenant.Phone,
			&tenant.RoleID,
			&tenant.Status,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
			&templeID,
			&templeName,
			&templePlace,
			&templeAddress,
			&templePhoneNo,
			&templeDescription,
			&templeCreatedAt,
			&templeUpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// If temple details exist, populate them
		if templeID != nil && templeName != nil {
			tenant.TempleDetails = &TenantTempleDetails{
				ID:                *templeID,
				TempleName:        *templeName,
				TemplePlace:       *templePlace,
				TempleAddress:     *templeAddress,
				TemplePhoneNo:     *templePhoneNo,
				TempleDescription: *templeDescription,
				CreatedAt:         *templeCreatedAt,
				UpdatedAt:         *templeUpdatedAt,
			}
		}

		tenants = append(tenants, tenant)
	}

	return tenants, total, nil
}

func (r *Repository) ApproveTenant(ctx context.Context, userID uint, adminID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Update("status", "active").Error
}

func (r *Repository) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Update("status", "rejected").Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
		Updates(map[string]interface{}{
			"status":      "rejected",
			"approved_by": adminID,
			"rejected_at": time.Now(),
			"admin_notes": reason,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repository) MarkTenantApprovalApproved(ctx context.Context, userID uint, adminID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": adminID,
			"approved_at": time.Now(),
			"updated_at":  time.Now(),
		}).Error
}

// =========================== ENTITY ===========================

func (r *Repository) GetPendingEntities(ctx context.Context) ([]entity.Entity, error) {
	var temples []entity.Entity
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Find(&temples).Error
	return temples, err
}

func (r *Repository) GetEntitiesWithFilters(ctx context.Context, status string, limit, page int) ([]entity.Entity, int64, error) {
	var temples []entity.Entity
	var total int64

	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&entity.Entity{})

	if status != "" {
		query = query.Where("LOWER(status) = LOWER(?)", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&temples).Error; err != nil {
		return nil, 0, err
	}

	return temples, total, nil
}

func (r *Repository) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("id = ?", entityID).
		Updates(map[string]interface{}{
			"status":     "approved",
			"updated_at": time.Now(),
		}).Error
}

func (r *Repository) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string, rejectedAt time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("id = ?", entityID).
		Update("status", "rejected").Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&auth.ApprovalRequest{}).
		Where("entity_id = ? AND request_type = ?", entityID, "entity").
		Updates(map[string]interface{}{
			"status":      "rejected",
			"approved_by": adminID,
			"admin_notes": reason,
			"rejected_at": rejectedAt,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repository) GetEntityByID(ctx context.Context, entityID uint) (entity.Entity, error) {
	var ent entity.Entity
	err := r.db.WithContext(ctx).
		Where("id = ?", entityID).
		First(&ent).Error
	return ent, err
}

func (r *Repository) LinkEntityToUser(ctx context.Context, userID uint, entityID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Update("entity_id", entityID).Error
}

func (r *Repository) MarkApprovalApproved(ctx context.Context, userID uint, adminID uint, entityID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "entity").
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": adminID,
			"approved_at": time.Now(),
			"entity_id":   entityID,
			"updated_at":  time.Now(),
		}).Error
}

// Count tenants (TempleAdmins) by status (active, pending, rejected)
func (r *Repository) CountTenantsByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ? AND LOWER(users.status) = LOWER(?)", "templeadmin", status).
		Count(&count).Error
	return count, err
}

// Count temples (Entities) by status
func (r *Repository) CountEntitiesByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("LOWER(status) = LOWER(?)", status).
		Count(&count).Error
	return count, err
}

// Count total users with role 'devotee'
func (r *Repository) CountUsersByRole(ctx context.Context, roleName string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("LOWER(user_roles.role_name) = LOWER(?)", roleName).
		Count(&count).Error
	return count, err
}

// =========================== USER MANAGEMENT ===========================

// Create user (admin-created users bypass email validation and approval process)
func (r *Repository) CreateUser(ctx context.Context, user *auth.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Create tenant details for templeadmin users
func (r *Repository) CreateTenantDetails(ctx context.Context, details *auth.TenantDetails) error {
	return r.db.WithContext(ctx).Create(details).Error
}

// Get users with pagination and filters (excluding devotee and volunteer roles)
func (r *Repository) GetUsers(ctx context.Context, limit, page int, search, roleFilter, statusFilter string) ([]UserResponse, int64, error) {
	var users []UserResponse
	var total int64

	offset := (page - 1) * limit

	// Build base query for total count
	baseQuery := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name NOT IN (?)", []string{"devotee", "volunteer"})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where(
			"users.full_name ILIKE ? OR users.email ILIKE ? OR users.phone ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Apply role filter
	if roleFilter != "" {
		baseQuery = baseQuery.Where("LOWER(user_roles.role_name) = LOWER(?)", roleFilter)
	}

	// Apply status filter
	if statusFilter != "" {
		baseQuery = baseQuery.Where("LOWER(users.status) = LOWER(?)", statusFilter)
	}

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build main query with all necessary joins and select fields
	query := r.db.WithContext(ctx).
		Table("users").
		Select(`
            users.id,
            users.full_name,
            users.email,
            users.phone,
            users.status,
            users.created_at,
            users.updated_at,
            user_roles.id as role_id,
            user_roles.role_name,
            td.id as temple_id,
            td.temple_name,
            td.temple_place,
            td.temple_address,
            td.temple_phone_no,
            td.temple_description,
            td.created_at as temple_created_at,
            td.updated_at as temple_updated_at,
            tenant_user.full_name as tenant_name,
            tua.created_at as assignment_created_at,
            tua.updated_at as assignment_updated_at
        `).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Joins("LEFT JOIN tenant_details td ON users.id = td.user_id AND user_roles.role_name = 'templeadmin'").
		Joins("LEFT JOIN tenant_user_assignments tua ON users.id = tua.user_id AND user_roles.role_name IN ('standarduser', 'monitoringuser') AND tua.status = 'active'").
		Joins("LEFT JOIN users tenant_user ON tua.tenant_id = tenant_user.id").
		Where("user_roles.role_name NOT IN (?)", []string{"devotee", "volunteer"})

	// Apply same filters to main query
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"users.full_name ILIKE ? OR users.email ILIKE ? OR users.phone ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if roleFilter != "" {
		query = query.Where("LOWER(user_roles.role_name) = LOWER(?)", roleFilter)
	}

	if statusFilter != "" {
		query = query.Where("LOWER(users.status) = LOWER(?)", statusFilter)
	}

	// Execute query with pagination and ordering
	rows, err := query.Limit(limit).Offset(offset).Order("users.created_at DESC").Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var user UserResponse
		// Temple details
		var templeID *uint
		var templeName, templePlace, templeAddress, templePhoneNo, templeDescription *string
		var templeCreatedAt, templeUpdatedAt *time.Time
		// Tenant assignment details
		var tenantName *string
		var assignmentCreatedAt, assignmentUpdatedAt *time.Time

		err := rows.Scan(
			&user.ID, &user.FullName, &user.Email, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt,
			&user.Role.ID, &user.Role.RoleName,
			&templeID, &templeName, &templePlace, &templeAddress, &templePhoneNo, &templeDescription, &templeCreatedAt, &templeUpdatedAt,
			&tenantName, &assignmentCreatedAt, &assignmentUpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Conditionally populate TempleDetails
		if templeID != nil {
			user.TempleDetails = &TenantTempleDetails{
				ID:                *templeID,
				TempleName:        *templeName,
				TemplePlace:       *templePlace,
				TempleAddress:     *templeAddress,
				TemplePhoneNo:     *templePhoneNo,
				TempleDescription: *templeDescription,
				CreatedAt:         *templeCreatedAt,
				UpdatedAt:         *templeUpdatedAt,
			}
		}

		// Conditionally populate TenantAssignmentDetails and set assignment flags
		if tenantName != nil {
			user.TenantAssignmentDetails = &TenantAssignmentDetails{
				TenantName: *tenantName,
				AssignedOn: *assignmentCreatedAt,
				UpdatedOn:  *assignmentUpdatedAt,
			}
			// Set assignment flags for UserResponse
			user.Assigned = true
			user.AssignedDate = assignmentCreatedAt
			user.ReassignmentDate = assignmentUpdatedAt
		} else {
			user.Assigned = false
		}

		users = append(users, user)
	}

	return users, total, nil
}

// Get user by ID with temple details
func (r *Repository) GetUserWithDetails(ctx context.Context, userID uint) (*UserResponse, error) {
	var user UserResponse
	var templeID *uint
	var templeName, templePlace, templeAddress, templePhoneNo, templeDescription *string
	var templeCreatedAt, templeUpdatedAt *time.Time
	var tenantName *string
	var assignmentCreatedAt, assignmentUpdatedAt *time.Time

	query := r.db.WithContext(ctx).
		Table("users").
		Select(`
            users.id,
            users.full_name,
            users.email,
            users.phone,
            users.status,
            users.created_at,
            users.updated_at,
            user_roles.id as role_id,
            user_roles.role_name,
            td.id as temple_id,
            td.temple_name,
            td.temple_place,
            td.temple_address,
            td.temple_phone_no,
            td.temple_description,
            td.created_at as temple_created_at,
            td.updated_at as temple_updated_at,
            tenant_user.full_name as tenant_name,
            tua.created_at as assignment_created_at,
            tua.updated_at as assignment_updated_at
        `).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Joins("LEFT JOIN tenant_details td ON users.id = td.user_id AND user_roles.role_name = 'templeadmin'").
		Joins("LEFT JOIN tenant_user_assignments tua ON users.id = tua.user_id AND user_roles.role_name IN ('standarduser', 'monitoringuser') AND tua.status = 'active'").
		Joins("LEFT JOIN users tenant_user ON tua.tenant_id = tenant_user.id").
		Where("users.id = ?", userID)

	row := query.Row()
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.ID,
		&user.Role.RoleName,
		&templeID,
		&templeName,
		&templePlace,
		&templeAddress,
		&templePhoneNo,
		&templeDescription,
		&templeCreatedAt,
		&templeUpdatedAt,
		&tenantName,
		&assignmentCreatedAt,
		&assignmentUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if templeID != nil {
		user.TempleDetails = &TenantTempleDetails{
			ID:                *templeID,
			TempleName:        *templeName,
			TemplePlace:       *templePlace,
			TempleAddress:     *templeAddress,
			TemplePhoneNo:     *templePhoneNo,
			TempleDescription: *templeDescription,
			CreatedAt:         *templeCreatedAt,
			UpdatedAt:         *templeUpdatedAt,
		}
	}

	// Conditionally populate TenantAssignmentDetails and set assignment flags
	if tenantName != nil {
		user.TenantAssignmentDetails = &TenantAssignmentDetails{
			TenantName: *tenantName,
			AssignedOn: *assignmentCreatedAt,
			UpdatedOn:  *assignmentUpdatedAt,
		}
		// Set assignment flags for UserResponse
		user.Assigned = true
		user.AssignedDate = assignmentCreatedAt
		user.ReassignmentDate = assignmentUpdatedAt
	} else {
		user.Assigned = false
	}

	return &user, nil
}


// Update user
func (r *Repository) UpdateUser(ctx context.Context, userID uint, user *auth.User) error {
	return r.db.WithContext(ctx).Model(&auth.User{}).Where("id = ?", userID).Updates(user).Error
}

// Update tenant details
func (r *Repository) UpdateTenantDetails(ctx context.Context, userID uint, details *auth.TenantDetails) error {
	return r.db.WithContext(ctx).Model(&auth.TenantDetails{}).Where("user_id = ?", userID).Updates(details).Error
}

// Delete user (soft delete)
func (r *Repository) DeleteUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Delete(&auth.User{}, userID).Error
}

// Update user status
func (r *Repository) UpdateUserStatus(ctx context.Context, userID uint, status string) error {
	return r.db.WithContext(ctx).Model(&auth.User{}).Where("id = ?", userID).Update("status", status).Error
}

// Get all user roles with complete information
func (r *Repository) GetUserRoles(ctx context.Context) ([]UserRole, error) {
	var roles []UserRole
	err := r.db.WithContext(ctx).
		Model(&auth.UserRole{}).
		Select("id, role_name, description, can_register_publicly").
		Find(&roles).Error
	return roles, err
}

// Find role by name
func (r *Repository) FindRoleByName(ctx context.Context, roleName string) (*auth.UserRole, error) {
	var role auth.UserRole
	err := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error
	return &role, err
}

// Check if user exists by email
func (r *Repository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&auth.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}


// =========================== USER ROLES ===========================

// Get all user roles (filtered by active status)
func (r *Repository) GetAllUserRoles(ctx context.Context) ([]auth.UserRole, error) {
    var roles []auth.UserRole
    err := r.db.WithContext(ctx).
        Where("status = ?", "active").
        Find(&roles).Error
    return roles, err
}

// GetUserRoleByID fetches a single role by its ID
func (r *Repository) GetUserRoleByID(ctx context.Context, roleID uint) (*auth.UserRole, error) {
    var role auth.UserRole
    err := r.db.WithContext(ctx).Where("id = ?", roleID).First(&role).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // Return nil if not found, not an error
        }
        return nil, err
    }
    return &role, nil
}

// Create a new user role
func (r *Repository) CreateUserRole(ctx context.Context, role *auth.UserRole) error {
    return r.db.WithContext(ctx).Create(role).Error
}

// CheckIfRoleNameExists checks if a role with the given name already exists
func (r *Repository) CheckIfRoleNameExists(ctx context.Context, roleName string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&auth.UserRole{}).
        Where("role_name = ?", roleName).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// UpdateUserRole saves the provided role object to the database
func (r *Repository) UpdateUserRole(ctx context.Context, role *auth.UserRole) error {
    return r.db.WithContext(ctx).Save(role).Error
}

// DeactivateUserRole updates the status of a role to 'inactive'
func (r *Repository) DeactivateUserRole(ctx context.Context, roleID uint) error {
    return r.db.WithContext(ctx).
        Model(&auth.UserRole{}).
        Where("id = ?", roleID).
        Update("status", "inactive").Error
}

// =========================== PASSWORD RESET ===========================

// GetUserByEmail retrieves a user by their email address
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
    var user auth.User
    result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// UpdateUserPassword updates a user's password
func (r *Repository) UpdateUserPassword(ctx context.Context, userID uint, newPasswordHash string) error {
    result := r.db.WithContext(ctx).Model(&auth.User{}).Where("id = ?", userID).
        Update("password_hash", newPasswordHash)
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return errors.New("user not found")
    }
    
    return nil
}


func (r *Repository) GetAssignableTenants(ctx context.Context, limit, page int) ([]AssignableTenant, int64, error) {
    var tenants []AssignableTenant
    var total int64

    // Calculate the offset based on the requested page and limit
    offset := (page - 1) * limit

    // First, count the total number of records that match the WHERE clause.
    // This is done without applying limit or offset.
    countQuery := r.db.WithContext(ctx).
        Table("users").
        Joins("JOIN user_roles ON users.role_id = user_roles.id").
        Where("user_roles.role_name = ? AND users.status = ?", "templeadmin", "active").
        Count(&total)

    if countQuery.Error != nil {
        return nil, 0, countQuery.Error
    }

    // Now, fetch the paginated data.
    // The same query is used, but with Select, Joins, Limit, and Offset.
    err := r.db.WithContext(ctx).
        Table("users").
        Select("users.id as user_id, users.full_name as tenant_name, users.email, COALESCE(entities.name, tenant_details.temple_name) AS temple_name, COALESCE(entities.street_address, tenant_details.temple_address) AS temple_address, COALESCE(entities.phone, tenant_details.temple_phone_no) AS temple_phone, COALESCE(entities.description, tenant_details.temple_description) AS temple_description").
        Joins("JOIN user_roles ON users.role_id = user_roles.id").
        Joins("LEFT JOIN entities ON users.id = entities.created_by").
        Joins("LEFT JOIN tenant_details ON users.id = tenant_details.user_id").
        Where("user_roles.role_name = ? AND users.status = ?", "templeadmin", "active").
        Limit(limit).
        Offset(offset).
        Scan(&tenants).Error

    if err != nil {
        return nil, 0, err
    }

    return tenants, total, nil
}

func (r *Repository) GetTenantsForSelection(ctx context.Context) ([]TenantSelectionResponse, error) {
    var tenants []TenantSelectionResponse

    // Modified query to explicitly join with tenant_details table and select fields directly
    query := `
        SELECT 
            u.id,
            u.full_name as name,
            u.email,
            td.temple_name,  -- Explicitly select temple_name
            td.temple_place, -- Explicitly select temple_place
            u.status,
            COALESCE(entity_count.count, 0) as temples_count
        FROM users u
        JOIN user_roles ur ON u.role_id = ur.id
        LEFT JOIN tenant_details td ON u.id = td.user_id
        LEFT JOIN (
            SELECT created_by, COUNT(*) as count 
            FROM entities 
            WHERE status = 'approved' 
            GROUP BY created_by
        ) entity_count ON u.id = entity_count.created_by
        WHERE ur.role_name = 'templeadmin' 
        AND u.status = 'active'
        ORDER BY u.full_name ASC
    `

    rows, err := r.db.WithContext(ctx).Raw(query).Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var tenant TenantSelectionResponse
        var templeName sql.NullString
        var location sql.NullString
        
        err := rows.Scan(
            &tenant.ID,
            &tenant.Name,
            &tenant.Email,
            &templeName,
            &location,
            &tenant.Status,
            &tenant.TemplesCount,
        )
        if err != nil {
            return nil, err
        }
        
        // Directly assign these fields to match the expected frontend field names
        tenant.TempleName = templeName.String
        tenant.Location = location.String
        
        tenants = append(tenants, tenant)
    }

    return tenants, nil
}

// Get assigned tenants for StandardUser / MonitoringUser
func (r *Repository) GetAssignedTenantsForUser(ctx context.Context, userID uint) ([]TenantSelectionResponse, error) {
	var tenants []TenantSelectionResponse

	query := `
		SELECT 
			tenant_user.id,
			tenant_user.full_name as name,
			tenant_user.email,
			COALESCE(td.temple_name, td.temple_place, '') as location,
			tenant_user.status,
			COALESCE(entity_count.count, 0) as temples_count
		FROM tenant_user_assignments tua
		JOIN users tenant_user ON tua.tenant_id = tenant_user.id
		JOIN user_roles ON tenant_user.role_id = user_roles.id
		LEFT JOIN tenant_details td ON tenant_user.id = td.user_id
		LEFT JOIN (
			SELECT created_by, COUNT(*) as count 
			FROM entities 
			WHERE status = 'approved' 
			GROUP BY created_by
		) entity_count ON tenant_user.id = entity_count.created_by
		WHERE tua.user_id = ? 
		AND tua.status = 'active'
		AND user_roles.role_name = 'templeadmin'
		AND tenant_user.status = 'active'
		ORDER BY tenant_user.full_name ASC
	`

	rows, err := r.db.WithContext(ctx).Raw(query, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tenant TenantSelectionResponse
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Email,
			&tenant.Location,
			&tenant.Status,
			&tenant.TemplesCount,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}



// New method to get tenants with temple details
// New method to get tenants with temple details
func (r *Repository) GetTenantsWithTempleDetails(ctx context.Context, role, status string) ([]TenantResponse, error) {
    var responses []TenantResponse
    
    // Base query with dynamic WHERE clause
    query := `
        SELECT 
            u.id, 
            u.full_name as "fullName",
            u.email,
            ur.role_name as "role",
            u.status,
            e.id as temple_id, 
            COALESCE(e.name, td.temple_name) as temple_name, 
            COALESCE(e.city, td.temple_place) as temple_city, 
            COALESCE(e.state, '') as temple_state
        FROM 
            users u
        JOIN 
            user_roles ur ON u.role_id = ur.id
        LEFT JOIN 
            tenant_details td ON u.id = td.user_id
        LEFT JOIN 
            entities e ON u.id = e.created_by
        WHERE 1=1
    `
    
    // Build dynamic query params
    params := []interface{}{}
    
    if role != "" {
        query += " AND ur.role_name = ?"
        params = append(params, role)
    }
    
    if status != "" {
        query += " AND u.status = ?"
        params = append(params, status)
    }
    
    rows, err := r.db.WithContext(ctx).Raw(query, params...).Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var tr TenantResponse
        var templeID sql.NullInt64
        var templeName, templeCity, templeState sql.NullString
        
        err := rows.Scan(
            &tr.ID,
            &tr.FullName,
            &tr.Email,
            &tr.Role,
            &tr.Status,
            &templeID,
            &templeName,
            &templeCity,
            &templeState,
        )
        
        if err != nil {
            return nil, err
        }
        
        if templeID.Valid && templeName.Valid {
            tr.Temple = &TempleDetails{
                ID:    uint(templeID.Int64),
                Name:  templeName.String,
                City:  templeCity.String,
                State: templeState.String,
            }
        }
        
        responses = append(responses, tr)
    }
    
    return responses, nil
}

// BulkCreateUsers inserts multiple users safely with better error handling
func (r *Repository) BulkCreateUsers(ctx context.Context, users []auth.User) error {
	if len(users) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var createdCount int
		
		for i, user := range users {
			// Check if email already exists
			var existingUser auth.User
			err := tx.Where("email = ?", user.Email).First(&existingUser).Error
			if err == nil {
				// User exists, skip
				fmt.Printf("User with email %s already exists, skipping\n", user.Email)
				continue
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				// Database error
				return fmt.Errorf("error checking existing user %s: %v", user.Email, err)
			}

			// Create the user
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("error creating user %d (%s): %v", i+1, user.Email, err)
			}
			
			createdCount++
			fmt.Printf("Successfully created user: %s\n", user.Email)
		}
		
		fmt.Printf("Transaction completed. Created %d users out of %d\n", createdCount, len(users))
		return nil
	})
}