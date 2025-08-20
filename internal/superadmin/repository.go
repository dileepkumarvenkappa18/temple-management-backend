package superadmin

import (
	"context"
	"time"
	"errors"
	

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"gorm.io/gorm"
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

	// Build base query excluding devotee and volunteer roles
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

	// Build main query with temple details for templeadmin users
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
			td.updated_at as temple_updated_at
		`).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Joins("LEFT JOIN tenant_details td ON users.id = td.user_id AND user_roles.role_name = 'templeadmin'").
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

	// Execute query with pagination
	rows, err := query.Limit(limit).Offset(offset).Order("users.created_at DESC").Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var user UserResponse
		var templeID *uint
		var templeName, templePlace, templeAddress, templePhoneNo, templeDescription *string
		var templeCreatedAt, templeUpdatedAt *time.Time

		err := rows.Scan(
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
		)
		if err != nil {
			return nil, 0, err
		}

		// If temple details exist, populate them
		if templeID != nil && templeName != nil {
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
			td.updated_at as temple_updated_at
		`).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Joins("LEFT JOIN tenant_details td ON users.id = td.user_id AND user_roles.role_name = 'templeadmin'").
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
	)

	if err != nil {
		return nil, err
	}

	// If temple details exist, populate them
	if templeID != nil && templeName != nil {
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





