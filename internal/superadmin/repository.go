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

// ✅ Mark entity approval in approval_requests
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

// ✅ Tenant approval
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

func (r *Repository) CreateEntity(ctx context.Context, ent *entity.Entity) error {
	return r.db.WithContext(ctx).Create(ent).Error
}

func (r *Repository) GetPendingTenants(ctx context.Context) ([]auth.User, error) {
	var tenants []auth.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ? AND users.status = ?", "templeadmin", "pending").
		Find(&tenants).Error
	return tenants, err
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

func (r *Repository) GetPendingEntities(ctx context.Context) ([]entity.Entity, error) {
	var temples []entity.Entity
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Find(&temples).Error
	return temples, err
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
	// Update entity table
	if err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("id = ?", entityID).
		Update("status", "rejected").Error; err != nil {
		return err
	}

	// Update approval_requests
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

// ✅ NEW: Get entity by ID (to fetch CreatedBy user)
func (r *Repository) GetEntityByID(ctx context.Context, entityID uint) (entity.Entity, error) {
	var ent entity.Entity
	err := r.db.WithContext(ctx).
		Where("id = ?", entityID).
		First(&ent).Error
	return ent, err
}

// ✅ NEW: Link approved entity to user
func (r *Repository) LinkEntityToUser(ctx context.Context, userID uint, entityID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("id = ?", userID).
		Update("entity_id", entityID).Error
}
