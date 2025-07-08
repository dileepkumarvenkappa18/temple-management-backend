package superadmin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ---------------- TENANT -----------------------

func (r *Repository) GetPendingTenants(ctx context.Context, page, limit int, search, status string) ([]auth.User, int64, error) {
	var tenants []auth.User
	var total int64

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if status == "" {
		status = "pending"
	}

	query := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.role_id = user_roles.id").
		Where("user_roles.role_name = ? AND users.status = ?", "templeadmin", status)

	// Optional Search
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(users.full_name) LIKE ? OR LOWER(users.email) LIKE ?", searchTerm, searchTerm)
	}

	// Count
	if err := query.Model(&auth.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tenants: %w", err)
	}

	// Query paginated data
	err := query.
		Order("users.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tenants).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tenants: %w", err)
	}

	return tenants, total, nil
}

func (r *Repository) ApproveTenant(ctx context.Context, userID uint, adminID uint) error {
	tx := r.db.WithContext(ctx).Begin()

	var user auth.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch user for approval: %w", err)
	}
	if user.Status != "pending" {
		tx.Rollback()
		return fmt.Errorf("user is already %s", user.Status)
	}

	if err := tx.Model(&auth.User{}).
		Where("id = ?", userID).
		Update("status", "active").Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
		Where("status = ?", "pending").
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": adminID,
			"approved_at": gorm.Expr("NOW()"),
			"updated_at":  gorm.Expr("NOW()"),
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) MarkTenantApprovalApproved(ctx context.Context, userID uint, adminID uint) error {
	return r.db.WithContext(ctx).
		Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": adminID,
			"approved_at": gorm.Expr("NOW()"),
			"updated_at":  gorm.Expr("NOW()"),
		}).Error
}

func (r *Repository) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	tx := r.db.WithContext(ctx).Begin()

	result := tx.Model(&auth.User{}).
		Where("id = ?", userID).
		Update("status", "rejected")
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("user not found or already rejected")
	}

	result = tx.Model(&auth.ApprovalRequest{}).
		Where("user_id = ? AND request_type = ?", userID, "tenant_approval").
		Updates(map[string]interface{}{
			"status":      "rejected",
			"approved_by": adminID,
			"admin_notes": reason,
			"rejected_at": gorm.Expr("NOW()"),
			"updated_at":  gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no approval request found for user")
	}

	return tx.Commit().Error
}

// ---------------- ENTITY -----------------------

func (r *Repository) GetPendingEntities(ctx context.Context, page, limit int, search, status string) ([]entity.Entity, int64, error) {
	var entities []entity.Entity
	var total int64

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if status == "" {
		status = "pending"
	}

	query := r.db.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("status = ?", status)

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count entities: %w", err)
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get entities: %w", err)
	}

	return entities, total, nil
}

func (r *Repository) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	tx := r.db.WithContext(ctx).Begin()

	var ent entity.Entity
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", entityID).
		First(&ent).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("entity not found: %w", err)
	}
	if ent.Status != "pending" {
		tx.Rollback()
		return fmt.Errorf("entity is already %s", ent.Status)
	}

	if err := tx.Model(&entity.Entity{}).
		Where("id = ?", entityID).
		Update("status", "approved").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string, rejectedAt time.Time) error {
	tx := r.db.WithContext(ctx).Begin()

	var ent entity.Entity
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", entityID).
		First(&ent).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("entity not found: %w", err)
	}
	if ent.Status != "pending" {
		tx.Rollback()
		return fmt.Errorf("cannot reject entity with status: %s", ent.Status)
	}

	if err := tx.Model(&entity.Entity{}).
		Where("id = ?", entityID).
		Update("status", "rejected").Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&auth.ApprovalRequest{}).
		Where("entity_id = ? AND request_type = ?", entityID, "entity").
		Where("status = ?", "pending").
		Updates(map[string]interface{}{
			"status":      "rejected",
			"approved_by": adminID,
			"admin_notes": reason,
			"rejected_at": gorm.Expr("NOW()"),
			"updated_at":  gorm.Expr("NOW()"),
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (auth.User, error) {
	var user auth.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error
	return user, err
}
