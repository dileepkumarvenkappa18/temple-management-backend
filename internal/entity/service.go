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
	ErrMissingFields = errors.New("temple name, deity, phone, and email are required")
)

// ========== ENTITY CORE ==========

// Temple Admin → Create Entity
func (s *Service) CreateEntity(e *Entity, userID uint) error {
	// Validate required fields
	if strings.TrimSpace(e.Name) == "" ||
		e.MainDeity == nil || strings.TrimSpace(*e.MainDeity) == "" ||
		strings.TrimSpace(e.Phone) == "" ||
		strings.TrimSpace(e.Email) == "" {
		return ErrMissingFields
	}

	now := time.Now()

	// Set metadata
	e.Status = "pending"
	e.CreatedBy = userID
	e.CreatedAt = now
	e.UpdatedAt = now

	// Sanitize inputs
	e.Name = strings.TrimSpace(e.Name)
	e.Email = strings.TrimSpace(e.Email)
	e.Phone = strings.TrimSpace(e.Phone)
	e.TempleType = strings.TrimSpace(e.TempleType)
	e.Description = strings.TrimSpace(e.Description)
	e.StreetAddress = strings.TrimSpace(e.StreetAddress)
	e.City = strings.TrimSpace(e.City)
	e.State = strings.TrimSpace(e.State)
	e.District = strings.TrimSpace(e.District)
	e.Pincode = strings.TrimSpace(e.Pincode)

	// Trim main deity if present
	if e.MainDeity != nil {
		trimmed := strings.TrimSpace(*e.MainDeity)
		e.MainDeity = &trimmed
	}

	// Trim document URLs
	e.RegistrationCertURL = strings.TrimSpace(e.RegistrationCertURL)
	e.TrustDeedURL = strings.TrimSpace(e.TrustDeedURL)
	e.PropertyDocsURL = strings.TrimSpace(e.PropertyDocsURL)
	e.AdditionalDocsURLs = strings.TrimSpace(e.AdditionalDocsURLs)

	// Save entity
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

// Super Admin / Admin → Get all temples
func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

// Anyone → View a temple by ID
func (s *Service) GetEntityByID(id int) (Entity, error) {
	return s.Repo.GetEntityByID(id)
}

// Temple Admin → Update own temple
func (s *Service) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	return s.Repo.UpdateEntity(e)
}

// Super Admin → Delete temple
func (s *Service) DeleteEntity(id int) error {
	return s.Repo.DeleteEntity(id)
}