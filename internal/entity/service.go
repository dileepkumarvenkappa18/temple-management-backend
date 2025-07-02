package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
)

type Service struct {
	Repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{Repo: r}
}

var (
	ErrMissingFields = errors.New("temple name and email are required")
)

// ========== ENTITY CORE ==========

func (s *Service) CreateEntity(e *Entity, userID uint) error {
	// Validate required fields
	if strings.TrimSpace(e.Name) == "" || strings.TrimSpace(e.Email) == "" {
		return ErrMissingFields
	}

	now := time.Now()

	// Populate required base fields
	e.Status = "pending"
	e.CreatedBy = userID
	e.CreatedAt = now
	e.UpdatedAt = now

	// Save the entity
	if err := s.Repo.CreateEntity(e); err != nil {
		return err
	}

	// Create approval request
	req := &auth.ApprovalRequest{
		UserID:      userID,
		EntityID:    &e.ID,
		RequestType: "temple_approval",
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.Repo.CreateApprovalRequest(req)
}

func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

func (s *Service) GetEntityByID(id int) (Entity, error) {
	return s.Repo.GetEntityByID(id)
}

func (s *Service) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	return s.Repo.UpdateEntity(e)
}

func (s *Service) DeleteEntity(id int) error {
	return s.Repo.DeleteEntity(id)
}
