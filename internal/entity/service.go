package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
)

type MembershipService interface {
	UpdateMembershipStatus(userID, entityID uint, status string) error
}

type Service struct {
	Repo              *Repository
	MembershipService MembershipService
}


func NewService(r *Repository, ms MembershipService) *Service {
	return &Service{
		Repo:              r,
		MembershipService: ms,
	}
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

// Super Admin → Get all temples
func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

// Temple Admin → Get entities created by specific user
func (s *Service) GetEntitiesByCreator(creatorID uint) ([]Entity, error) {
	return s.Repo.GetEntitiesByCreator(creatorID)
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

// ========== DEVOTEE MANAGEMENT ==========

// Temple Admin → Get devotees for specific entity
func (s *Service) GetDevotees(entityID uint) ([]DevoteeDTO, error) {
	return s.Repo.GetDevoteesByEntityID(entityID)
}

// Temple Admin → Get devotee statistics for entity
func (s *Service) GetDevoteeStats(entityID uint) (DevoteeStats, error) {
	return s.Repo.GetDevoteeStats(entityID)
}


// DashboardSummary is the structured JSON response
type DashboardSummary struct {
	RegisteredDevotees struct {
		Total     int64 `json:"total"`
		ThisMonth int64 `json:"this_month"`
	} `json:"registered_devotees"`

	SevaBookings struct {
		Today     int64 `json:"today"`
		ThisMonth int64 `json:"this_month"`
	} `json:"seva_bookings"`

	MonthDonations struct {
		Amount        float64 `json:"amount"`
		PercentChange float64 `json:"percent_change"`
	} `json:"month_donations"`

	UpcomingEvents struct {
		Total     int64 `json:"total"`
		ThisWeek  int64 `json:"this_week"`
	} `json:"upcoming_events"`
}

// Temple Admin → Dashboard Summary
func (s *Service) GetDashboardSummary(entityID uint) (DashboardSummary, error) {
	var summary DashboardSummary

	// ============= DEVOTEES =============
	totalDevotees, err := s.Repo.CountDevotees(entityID)
	if err != nil {
		return summary, err
	}
	thisMonthDevotees, err := s.Repo.CountDevoteesThisMonth(entityID)
	if err != nil {
		return summary, err
	}
	summary.RegisteredDevotees.Total = totalDevotees
	summary.RegisteredDevotees.ThisMonth = thisMonthDevotees

	// ============= SEVA BOOKINGS =============
	todaySevas, err := s.Repo.CountSevaBookingsToday(entityID)
	if err != nil {
		return summary, err
	}
	monthSevas, err := s.Repo.CountSevaBookingsThisMonth(entityID)
	if err != nil {
		return summary, err
	}
	summary.SevaBookings.Today = todaySevas
	summary.SevaBookings.ThisMonth = monthSevas

	// ============= DONATIONS =============
	monthDonationAmount, percentChange, err := s.Repo.GetMonthDonationsWithChange(entityID)
	if err != nil {
		return summary, err
	}
	summary.MonthDonations.Amount = monthDonationAmount
	summary.MonthDonations.PercentChange = percentChange

	// ============= UPCOMING EVENTS =============
	totalUpcoming, err := s.Repo.CountUpcomingEvents(entityID)
	if err != nil {
		return summary, err
	}
	thisWeekUpcoming, err := s.Repo.CountUpcomingEventsThisWeek(entityID)
	if err != nil {
		return summary, err
	}
	summary.UpcomingEvents.Total = totalUpcoming
	summary.UpcomingEvents.ThisWeek = thisWeekUpcoming

	return summary, nil
}
