package entity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/donation"
)

type Service struct {
	Repo         Repository
	DonationRepo donation.Repository
}

func NewService(r Repository, d donation.Repository) *Service {
	return &Service{
		Repo:         r,
		DonationRepo: d,
	}
}

var (
	ErrMissingFields = errors.New("temple name, deity, phone, and email are required")
)

// CreateEntity handles the creation of a new temple by a Temple Admin
func (s *Service) CreateEntity(e *Entity, userID uint) error {
	// Basic validation
	if strings.TrimSpace(e.Name) == "" ||
		e.MainDeity == nil || strings.TrimSpace(*e.MainDeity) == "" ||
		strings.TrimSpace(e.Phone) == "" ||
		strings.TrimSpace(e.Email) == "" {
		return ErrMissingFields
	}

	now := time.Now()
	e.Status = "pending"
	e.CreatedBy = userID
	e.CreatedAt = now
	e.UpdatedAt = now

	// Field sanitization
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
	e.MapLink = strings.TrimSpace(e.MapLink)
	e.RegistrationCertURL = strings.TrimSpace(e.RegistrationCertURL)
	e.TrustDeedURL = strings.TrimSpace(e.TrustDeedURL)
	e.PropertyDocsURL = strings.TrimSpace(e.PropertyDocsURL)
	e.AdditionalDocsURLs = strings.TrimSpace(e.AdditionalDocsURLs)

	if e.MainDeity != nil {
		mainDeity := strings.TrimSpace(*e.MainDeity)
		e.MainDeity = &mainDeity
	}

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

// GetAllEntities fetches all temples (for Super Admin)
func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

// GetEntitiesByCreator fetches temples created by a specific user (Temple Admin)
func (s *Service) GetEntitiesByCreator(userID uint) ([]Entity, error) {
	return s.Repo.GetEntitiesByCreator(userID)
}

// GetEntityByID retrieves a temple by ID (public)
func (s *Service) GetEntityByID(id int) (Entity, error) {
	return s.Repo.GetEntityByID(id)
}

// UpdateEntity modifies an existing temple (Temple Admin)
func (s *Service) UpdateEntity(e Entity) error {
	e.UpdatedAt = time.Now()
	return s.Repo.UpdateEntity(e)
}

// DeleteEntity removes a temple (Super Admin)
func (s *Service) DeleteEntity(id int) error {
	return s.Repo.DeleteEntity(id)
}

// GetEntityDashboardStats returns donation summary stats for a temple
func (s *Service) GetEntityDashboardStats(ctx context.Context, entityID int) (map[string]interface{}, error) {
	totalAmount, err := s.DonationRepo.GetTotalAmountByEntityID(ctx, uint(entityID))
	if err != nil {
		return nil, err
	}

	donationCount, err := s.DonationRepo.GetDonationCountByEntityID(ctx, uint(entityID))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalAmount":     totalAmount,
		"donationCount":   donationCount,
		"lastUpdatedTime": time.Now(),
	}, nil
}

// GetEntityDashboard is a context-less wrapper for dashboard stats
func (s *Service) GetEntityDashboard(entityID uint) (interface{}, error) {
	ctx := context.Background()
	return s.GetEntityDashboardStats(ctx, int(entityID))
}

// GetTopDonors returns the top donors of a temple
func (s *Service) GetTopDonors(ctx context.Context, templeID uint) ([]donation.TopDonor, error) {
	return s.DonationRepo.GetTopDonors(ctx, templeID)
}