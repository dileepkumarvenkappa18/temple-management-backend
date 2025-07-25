package superadmin

import (
	"context"
	"errors"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
)

// Service handles the business logic for superadmin functionalities.
type Service struct {
	repo *Repository
}

// NewService creates a new instance of the superadmin Service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ============ TENANT LOGIC ============ //

// ApproveTenant activates a tenant user and marks the approval request.
func (s *Service) ApproveTenant(ctx context.Context, userID, adminID uint) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	switch user.Status {
	case "active":
		return errors.New("tenant already approved")
	case "rejected":
		return errors.New("tenant already rejected")
	}

	return s.repo.ApproveTenant(ctx, userID, adminID)
}

// RejectTenant marks a tenant's request as rejected with a reason.
func (s *Service) RejectTenant(ctx context.Context, userID, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	switch user.Status {
	case "rejected":
		return errors.New("tenant already rejected")
	case "active":
		return errors.New("tenant already approved")
	}

	return s.repo.RejectTenant(ctx, userID, adminID, reason)
}

// GetTenantsWithFilters fetches paginated tenants filtered by status.
func (s *Service) GetTenantsWithFilters(ctx context.Context, status string, limit, page int) ([]auth.User, int64, error) {
	return s.repo.GetTenantsWithFilters(ctx, status, limit, page)
}

// ============ ENTITY LOGIC ============ //

// ApproveEntity sets an entity as approved and links it to the user.
func (s *Service) ApproveEntity(ctx context.Context, entityID, adminID uint) error {
	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	switch ent.Status {
	case "approved":
		return errors.New("entity already approved")
	case "rejected":
		return errors.New("entity already rejected")
	}

	if err := s.repo.ApproveEntity(ctx, entityID, adminID); err != nil {
		return err
	}

	return s.repo.LinkEntityToUser(ctx, ent.CreatedBy, entityID)
}

// RejectEntity marks an entity as rejected with a reason.
func (s *Service) RejectEntity(ctx context.Context, entityID, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	switch ent.Status {
	case "rejected":
		return errors.New("entity already rejected")
	case "approved":
		return errors.New("entity already approved")
	}

	return s.repo.RejectEntity(ctx, entityID, adminID, reason)
}

// GetEntitiesWithFilters returns paginated entities filtered by status.
func (s *Service) GetEntitiesWithFilters(ctx context.Context, status string, limit, page int) ([]entity.Entity, int64, error) {
	return s.repo.GetEntitiesWithFilters(ctx, status, limit, page)
}

// ============ DASHBOARD ============ //

// GetDashboardMetrics provides overview metrics for the admin dashboard.
func (s *Service) GetDashboardMetrics(ctx context.Context) (DashboardMetrics, error) {
	return s.repo.GetDashboardMetrics(ctx)
}