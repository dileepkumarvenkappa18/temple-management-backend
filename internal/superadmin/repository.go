package superadmin

import (
	"context"
	"time"

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