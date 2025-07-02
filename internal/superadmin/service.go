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
	// Step 1: Set user status to "active"
	if err := s.repo.ApproveTenant(ctx, userID, adminID); err != nil {
		return err
	}

	// Step 2: Update approval_requests (tenant approval)
	return s.repo.MarkTenantApprovalApproved(ctx, userID, adminID)
}

func (s *Service) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}
	return s.repo.RejectTenant(ctx, userID, adminID, reason)
}

func (s *Service) GetPendingTenants(ctx context.Context) ([]auth.User, error) {
	return s.repo.GetPendingTenants(ctx)
}

// ================== ENTITY ==================

func (s *Service) GetPendingEntities(ctx context.Context) ([]entity.Entity, error) {
	return s.repo.GetPendingEntities(ctx)
}

func (s *Service) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	// Step 1: Approve the entity (entity.status = approved, approval_request updated)
	err := s.repo.ApproveEntity(ctx, entityID, adminID)
	if err != nil {
		return err
	}

	// Step 2: Link entity to the user who created it
	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return err
	}

	return s.repo.LinkEntityToUser(ctx, ent.CreatedBy, ent.ID)
}

func (s *Service) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	rejectedAt := time.Now()

	return s.repo.RejectEntity(ctx, entityID, adminID, reason, rejectedAt)
}
