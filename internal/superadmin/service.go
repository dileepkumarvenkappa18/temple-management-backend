package superadmin

import (
	"context"
	"errors"
	"fmt"
	"time"
	"strings"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"github.com/sharath018/temple-management-backend/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ====================== TENANT APPROVAL ======================

func (s *Service) ApproveTenant(ctx context.Context, userID uint, adminID uint) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.Status == "active" {
		return errors.New("user is already approved")
	}
	if user.Status == "rejected" {
		return errors.New("cannot approve a rejected user")
	}
	if user.Status != "pending" {
		return fmt.Errorf("cannot approve user with status: %s", user.Status)
	}

	if err := s.repo.ApproveTenant(ctx, userID, adminID); err != nil {
		return fmt.Errorf("failed to activate tenant user: %w", err)
	}

	if err := s.repo.MarkTenantApprovalApproved(ctx, userID, adminID); err != nil {
		return fmt.Errorf("failed to mark approval_request approved: %w", err)
	}

	utils.SendTenantApprovalEmail(user.Email, user.FullName)
	return nil
}

func (s *Service) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.Status == "rejected" {
		return errors.New("user is already rejected")
	}
	if user.Status == "active" {
		return errors.New("cannot reject an already approved user")
	}
	if user.Status != "pending" {
		return fmt.Errorf("cannot reject user with status: %s", user.Status)
	}

	if err := s.repo.RejectTenant(ctx, userID, adminID, reason); err != nil {
		return fmt.Errorf("failed to reject tenant: %w", err)
	}

	utils.SendTenantRejectionEmail(user.Email, user.FullName, reason)
	return nil
}

func (s *Service) GetPendingTenants(ctx context.Context, page int, limit int, search, status string) ([]auth.User, int64, error) {
	return s.repo.GetPendingTenants(ctx, page, limit, strings.TrimSpace(search), strings.TrimSpace(status))
}

// ====================== ENTITY APPROVAL ======================

func (s *Service) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	if ent.Status == "approved" {
		return errors.New("entity is already approved")
	}
	if ent.Status == "rejected" {
		return errors.New("cannot approve a rejected entity")
	}
	if ent.Status != "pending" {
		return fmt.Errorf("cannot approve entity with status: %s", ent.Status)
	}

	if err := s.repo.ApproveEntity(ctx, entityID, adminID); err != nil {
		return fmt.Errorf("failed to update entity status: %w", err)
	}

	if err := s.repo.LinkEntityToUser(ctx, ent.CreatedBy, ent.ID); err != nil {
		return fmt.Errorf("failed to link entity to user: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, ent.CreatedBy)
	if err == nil {
		utils.SendEntityApprovalEmail(user.Email, user.FullName, ent.Name)
	}

	return nil
}

func (s *Service) GetPendingEntities(ctx context.Context, page int, limit int, search, status string) ([]entity.Entity, int64, error) {
	return s.repo.GetPendingEntities(ctx, page, limit, strings.TrimSpace(search), strings.TrimSpace(status))
}

func (s *Service) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	if ent.Status == "rejected" {
		return errors.New("entity is already rejected")
	}
	if ent.Status == "approved" {
		return errors.New("cannot reject an already approved entity")
	}
	if ent.Status != "pending" {
		return fmt.Errorf("cannot reject entity with status: %s", ent.Status)
	}

	if err := s.repo.RejectEntity(ctx, entityID, adminID, reason, time.Now()); err != nil {
		return fmt.Errorf("failed to reject entity: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, ent.CreatedBy)
	if err == nil {
		utils.SendEntityRejectionEmail(user.Email, user.FullName, ent.Name, reason)
	}

	return nil
}
