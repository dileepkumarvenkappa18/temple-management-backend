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

// ====================== TENANT ======================

func (r *Repository) GetTenantsWithFilters(ctx context.Context, status string, limit, page int) ([]auth.User, int64, error) {
	var tenants []auth.User
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ?", "templeadmin")

	if status != "" {
		query = query.Where("LOWER(users.status) = LOWER(?)", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

func (r *Repository) ApproveTenant(ctx context.Context, userID, adminID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&auth.User{}).
			Where("id = ?", userID).
			Update("status", "active").Error; err != nil {
			return err
		}

		return tx.Model(&auth.ApprovalRequest{}).
			Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
			Updates(map[string]interface{}{
				"status":      "approved",
				"approved_by": adminID,
				"approved_at": time.Now(),
				"updated_at":  time.Now(),
			}).Error
	})
}

func (r *Repository) RejectTenant(ctx context.Context, userID, adminID uint, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&auth.User{}).
			Where("id = ?", userID).
			Update("status", "rejected").Error; err != nil {
			return err
		}

		return tx.Model(&auth.ApprovalRequest{}).
			Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
			Updates(map[string]interface{}{
				"status":      "rejected",
				"approved_by": adminID,
				"admin_notes": reason,
				"rejected_at": time.Now(),
				"updated_at":  time.Now(),
			}).Error
	})
}

// ====================== ENTITY ======================

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

func (r *Repository) ApproveEntity(ctx context.Context, entityID, adminID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Entity{}).
			Where("id = ?", entityID).
			Updates(map[string]interface{}{
				"status":     "approved",
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return tx.Model(&auth.ApprovalRequest{}).
			Where("entity_id = ? AND request_type = ?", entityID, "entity").
			Updates(map[string]interface{}{
				"status":      "approved",
				"approved_by": adminID,
				"approved_at": time.Now(),
				"updated_at":  time.Now(),
			}).Error
	})
}

func (r *Repository) RejectEntity(ctx context.Context, entityID, adminID uint, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Entity{}).
			Where("id = ?", entityID).
			Update("status", "rejected").Error; err != nil {
			return err
		}

		return tx.Model(&auth.ApprovalRequest{}).
			Where("entity_id = ? AND request_type = ?", entityID, "entity").
			Updates(map[string]interface{}{
				"status":      "rejected",
				"approved_by": adminID,
				"admin_notes": reason,
				"rejected_at": time.Now(),
				"updated_at":  time.Now(),
			}).Error
	})
}

// ====================== DASHBOARD ======================

func (r *Repository) GetDashboardMetrics(ctx context.Context) (DashboardMetrics, error) {
	var metrics DashboardMetrics

	if err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("status = ?", "pending").
		Count(&metrics.PendingApprovals).Error; err != nil {
		return metrics, err
	}

	if err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("status = ?", "approved").
		Count(&metrics.ActiveTemples).Error; err != nil {
		return metrics, err
	}

	if err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Count(&metrics.TotalUsers).Error; err != nil {
		return metrics, err
	}

	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	if err := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("status = ? AND updated_at >= ?", "rejected", startOfMonth).
		Count(&metrics.RejectedCount).Error; err != nil {
		return metrics, err
	}

	return metrics, nil
}

// ====================== SUPPORTING ======================

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetEntityByID(ctx context.Context, entityID uint) (*entity.Entity, error) {
	var ent entity.Entity
	if err := r.db.WithContext(ctx).First(&ent, entityID).Error; err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *Repository) LinkEntityToUser(ctx context.Context, userID, entityID uint) error {
	link := entity.UserEntity{
		UserID:   userID,
		EntityID: entityID,
	}
	return r.db.WithContext(ctx).Create(&link).Error
}