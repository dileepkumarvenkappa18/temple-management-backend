package superadmin

import (
	"context"
	"errors"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ================== TENANT ==================

func (s *Service) ApproveTenant(ctx context.Context, userID uint, adminID uint) error {
	// Check existence and current status
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	if user.Status == "active" {
		return errors.New("tenant already approved")
	}
	if user.Status == "rejected" {
		return errors.New("tenant already rejected")
	}

	if err := s.repo.ApproveTenant(ctx, userID, adminID); err != nil {
		return err
	}
	return s.repo.MarkTenantApprovalApproved(ctx, userID, adminID)
}

func (s *Service) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	if user.Status == "rejected" {
		return errors.New("tenant already rejected")
	}
	if user.Status == "active" {
		return errors.New("tenant already approved")
	}

	return s.repo.RejectTenant(ctx, userID, adminID, reason)
}

func (s *Service) GetPendingTenants(ctx context.Context) ([]auth.User, error) {
	return s.repo.GetPendingTenants(ctx)
}

func (s *Service) GetTenantsWithFilters(ctx context.Context, status string, limit, page int) ([]TenantWithDetails, int64, error) {
	return s.repo.GetTenantsWithFilters(ctx, status, limit, page)
}

func (s *Service) UpdateTenantApprovalStatus(ctx context.Context, userID, adminID uint, action string, reason string) error {
	switch action {
	case "approve":
		return s.ApproveTenant(ctx, userID, adminID)
	case "reject":
		return s.RejectTenant(ctx, userID, adminID, reason)
	default:
		return errors.New("invalid action: must be approve or reject")
	}
}

// ================== ENTITY ==================

func (s *Service) GetPendingEntities(ctx context.Context) ([]entity.Entity, error) {
	return s.repo.GetPendingEntities(ctx)
}

func (s *Service) GetEntitiesWithFilters(ctx context.Context, status string, limit, page int) ([]entity.Entity, int64, error) {
	return s.repo.GetEntitiesWithFilters(ctx, status, limit, page)
}

func (s *Service) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	if ent.Status == "approved" {
		return errors.New("entity already approved")
	}
	if ent.Status == "rejected" {
		return errors.New("entity already rejected")
	}

	if err := s.repo.ApproveEntity(ctx, entityID, adminID); err != nil {
		return err
	}

	return s.repo.LinkEntityToUser(ctx, ent.CreatedBy, ent.ID)
}

func (s *Service) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	if ent.Status == "rejected" {
		return errors.New("entity already rejected")
	}
	if ent.Status == "approved" {
		return errors.New("entity already approved")
	}

	rejectedAt := time.Now()
	return s.repo.RejectEntity(ctx, entityID, adminID, reason, rejectedAt)
}

func (s *Service) UpdateEntityApprovalStatus(ctx context.Context, entityID, adminID uint, action string, reason string) error {
	switch action {
	case "approve":
		return s.ApproveEntity(ctx, entityID, adminID)
	case "reject":
		return s.RejectEntity(ctx, entityID, adminID, reason)
	default:
		return errors.New("invalid action: must be approve or reject")
	}
}

// ================== METRIC COUNTS ==================

// Tenant approval counts for SuperAdmin dashboard
func (s *Service) GetTenantApprovalCounts(ctx context.Context) (*TenantApprovalCount, error) {
	approved, err := s.repo.CountTenantsByStatus(ctx, "active") // assuming "active" means approved
	if err != nil {
		return nil, err
	}

	pending, err := s.repo.CountTenantsByStatus(ctx, "pending")
	if err != nil {
		return nil, err
	}

	rejected, err := s.repo.CountTenantsByStatus(ctx, "rejected")
	if err != nil {
		return nil, err
	}

	return &TenantApprovalCount{
		Approved: approved,
		Pending:  pending,
		Rejected: rejected,
	}, nil
}

// Temple (entity) approval counts for dashboard
func (s *Service) GetTempleApprovalCounts(ctx context.Context) (*TempleApprovalCount, error) {
	pending, err := s.repo.CountEntitiesByStatus(ctx, "PENDING")
	if err != nil {
		return nil, err
	}

	active, err := s.repo.CountEntitiesByStatus(ctx, "APPROVED")
	if err != nil {
		return nil, err
	}

	rejected, err := s.repo.CountEntitiesByStatus(ctx, "REJECTED")
	if err != nil {
		return nil, err
	}

	totalDevotees, err := s.repo.CountUsersByRole(ctx, "devotee")
	if err != nil {
		return nil, err
	}

	return &TempleApprovalCount{
		PendingApproval: pending,
		ActiveTemples:   active,
		Rejected:        rejected,
		TotalDevotees:   totalDevotees,
	}, nil
}