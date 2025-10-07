package entity

import (
	"context"
	"errors"
<<<<<<< HEAD
	"log"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/auth"
=======
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
)

type MembershipService interface {
	UpdateMembershipStatus(userID, entityID uint, status string) error
}

type Service struct {
	Repo              *Repository
	MembershipService MembershipService
<<<<<<< HEAD
	AuditService      auditlog.Service
}



=======
	AuditService      auditlog.Service // 🆕 ADD AUDIT SERVICE
}

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func NewService(r *Repository, ms MembershipService, as auditlog.Service) *Service {
	return &Service{
		Repo:              r,
		MembershipService: ms,
<<<<<<< HEAD
		AuditService:      as,
=======
		AuditService:      as, // 🆕 INJECT AUDIT SERVICE
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	}
}

var (
	ErrMissingFields = errors.New("temple name, deity, phone, and email are required")
)

// ========== ENTITY CORE ==========

<<<<<<< HEAD
// CreateEntity - Create temple with auto-approval for superadmin (role_id = 1)
func (s *Service) CreateEntity(e *Entity, userID uint, userRoleID uint, ip string) error {
=======
// Temple Admin → Create Entity
func (s *Service) CreateEntity(e *Entity, userID uint, ip string) error {
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	// Validate required fields
	if strings.TrimSpace(e.Name) == "" ||
		e.MainDeity == nil || strings.TrimSpace(*e.MainDeity) == "" ||
		strings.TrimSpace(e.Phone) == "" ||
		strings.TrimSpace(e.Email) == "" {
<<<<<<< HEAD

		auditDetails := map[string]interface{}{
			"temple_name": strings.TrimSpace(e.Name),
			"email":       strings.TrimSpace(e.Email),
			"role_id":     userRoleID,
			"error":       "Missing required fields",
		}
		s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")

=======
		
		// 🔍 LOG FAILED TEMPLE CREATION ATTEMPT
		auditDetails := map[string]interface{}{
			"temple_name": strings.TrimSpace(e.Name),
			"email":       strings.TrimSpace(e.Email),
			"error":       "Missing required fields",
		}
		s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")
		
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		return ErrMissingFields
	}

	now := time.Now()

<<<<<<< HEAD
	// AUTO-APPROVE LOGIC: Check if creator is superadmin (role_id = 1)
	if e.Status == "" {
		if userRoleID == 1 {
			e.Status = "approved"
			log.Printf("🎉 Temple auto-approved: Created by superadmin (user_id: %d, role_id: %d)", userID, userRoleID)
		} else {
			e.Status = "pending"
			log.Printf("📝 Temple pending approval: Created by role_id: %d (user_id: %d)", userRoleID, userID)
		}
	}

	// Set metadata
	e.CreatedAt = now
	e.UpdatedAt = now
	e.CreatorRoleID = &userRoleID
=======
	// Set metadata
	e.Status = "pending"
	e.CreatedBy = userID
	e.CreatedAt = now
	e.UpdatedAt = now
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6

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

<<<<<<< HEAD
=======
	// Trim main deity if present
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	if e.MainDeity != nil {
		trimmed := strings.TrimSpace(*e.MainDeity)
		e.MainDeity = &trimmed
	}

<<<<<<< HEAD
=======
	// Trim document URLs
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	e.RegistrationCertURL = strings.TrimSpace(e.RegistrationCertURL)
	e.TrustDeedURL = strings.TrimSpace(e.TrustDeedURL)
	e.PropertyDocsURL = strings.TrimSpace(e.PropertyDocsURL)
	e.AdditionalDocsURLs = strings.TrimSpace(e.AdditionalDocsURLs)

<<<<<<< HEAD
	// Save entity to database
	if err := s.Repo.CreateEntity(e); err != nil {
		auditDetails := map[string]interface{}{
			"temple_name": e.Name,
			"email":       e.Email,
			"role_id":     userRoleID,
			"error":       err.Error(),
		}
		s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")

		return err
	}

	// ONLY create approval request if NOT auto-approved (i.e., status is pending)
	if e.Status == "pending" {
		req := &auth.ApprovalRequest{
			UserID:      userID,
			EntityID:    &e.ID,
			RequestType: "temple_approval",
			Status:      "pending",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := s.Repo.CreateApprovalRequest(req); err != nil {
			auditDetails := map[string]interface{}{
				"temple_name": e.Name,
				"temple_id":   e.ID,
				"email":       e.Email,
				"role_id":     userRoleID,
				"error":       err.Error(),
			}
			s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_APPROVAL_REQUEST_FAILED", auditDetails, ip, "failure")

			return err
		}

		log.Printf("✅ Approval request created for temple ID: %d", e.ID)
	} else {
		log.Printf("⚡ Skipped approval request - Temple auto-approved (ID: %d)", e.ID)
	}

	// Log successful temple creation with appropriate action type
=======
	// Save entity
	if err := s.Repo.CreateEntity(e); err != nil {
		// 🔍 LOG FAILED TEMPLE CREATION (DB ERROR)
		auditDetails := map[string]interface{}{
			"temple_name": e.Name,
			"email":       e.Email,
			"error":       err.Error(),
		}
		s.AuditService.LogAction(context.Background(), &userID, nil, "TEMPLE_CREATE_FAILED", auditDetails, ip, "failure")
		
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

	if err := s.Repo.CreateApprovalRequest(req); err != nil {
		// 🔍 LOG FAILED APPROVAL REQUEST CREATION
		auditDetails := map[string]interface{}{
			"temple_name": e.Name,
			"temple_id":   e.ID,
			"email":       e.Email,
			"error":       err.Error(),
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_APPROVAL_REQUEST_FAILED", auditDetails, ip, "failure")
		
		return err
	}

	// 🔍 LOG SUCCESSFUL TEMPLE CREATION
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	auditDetails := map[string]interface{}{
		"temple_name":   e.Name,
		"temple_id":     e.ID,
		"temple_type":   e.TempleType,
		"email":         e.Email,
		"phone":         e.Phone,
		"city":          e.City,
		"state":         e.State,
		"main_deity":    e.MainDeity,
		"status":        e.Status,
<<<<<<< HEAD
		"role_id":       userRoleID,
		"auto_approved": e.Status == "approved",
	}

	actionType := "TEMPLE_CREATED"
	if e.Status == "approved" {
		actionType = "TEMPLE_CREATED_AUTO_APPROVED"
	}

	s.AuditService.LogAction(context.Background(), &userID, &e.ID, actionType, auditDetails, ip, "success")
=======
	}
	s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_CREATED", auditDetails, ip, "success")
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6

	return nil
}

<<<<<<< HEAD
// GetAllEntities - Super Admin → Get all temples
=======
// Super Admin → Get all temples
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

<<<<<<< HEAD
// GetEntitiesByCreator - Temple Admin → Get entities created by specific user
=======
// Temple Admin → Get entities created by specific user
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (s *Service) GetEntitiesByCreator(creatorID uint) ([]Entity, error) {
	return s.Repo.GetEntitiesByCreator(creatorID)
}

<<<<<<< HEAD
// GetEntityByID - Anyone → View a temple by ID
=======
// Anyone → View a temple by ID
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (s *Service) GetEntityByID(id int) (Entity, error) {
	return s.Repo.GetEntityByID(id)
}

<<<<<<< HEAD
// UpdateEntity - Temple Admin → Update own temple
func (s *Service) UpdateEntity(e Entity, userID uint, ip string) error {
	existingEntity, err := s.Repo.GetEntityByID(int(e.ID))
	if err != nil {
=======
// Temple Admin → Update own temple
func (s *Service) UpdateEntity(e Entity, userID uint, ip string) error {
	// Get existing entity for comparison
	existingEntity, err := s.Repo.GetEntityByID(int(e.ID))
	if err != nil {
		// 🔍 LOG FAILED TEMPLE UPDATE ATTEMPT (NOT FOUND)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		auditDetails := map[string]interface{}{
			"temple_id": e.ID,
			"error":     "Temple not found",
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATE_FAILED", auditDetails, ip, "failure")
<<<<<<< HEAD

=======
		
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		return err
	}

	e.UpdatedAt = time.Now()
<<<<<<< HEAD

	if err := s.Repo.UpdateEntity(e); err != nil {
=======
	
	if err := s.Repo.UpdateEntity(e); err != nil {
		// 🔍 LOG FAILED TEMPLE UPDATE (DB ERROR)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		auditDetails := map[string]interface{}{
			"temple_id":   e.ID,
			"temple_name": e.Name,
			"error":       err.Error(),
		}
		s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATE_FAILED", auditDetails, ip, "failure")
<<<<<<< HEAD

		return err
	}

	auditDetails := map[string]interface{}{
		"temple_id":      e.ID,
		"temple_name":    e.Name,
		"previous_name":  existingEntity.Name,
		"temple_type":    e.TempleType,
		"email":          e.Email,
		"phone":          e.Phone,
		"city":           e.City,
		"state":          e.State,
		"main_deity":     e.MainDeity,
		"description":    e.Description,
		"updated_fields": getUpdatedFields(existingEntity, e),
=======
		
		return err
	}

	// 🔍 LOG SUCCESSFUL TEMPLE UPDATE
	auditDetails := map[string]interface{}{
		"temple_id":         e.ID,
		"temple_name":       e.Name,
		"previous_name":     existingEntity.Name,
		"temple_type":       e.TempleType,
		"email":             e.Email,
		"phone":             e.Phone,
		"city":              e.City,
		"state":             e.State,
		"main_deity":        e.MainDeity,
		"description":       e.Description,
		"updated_fields":    getUpdatedFields(existingEntity, e),
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	}
	s.AuditService.LogAction(context.Background(), &userID, &e.ID, "TEMPLE_UPDATED", auditDetails, ip, "success")

	return nil
}
<<<<<<< HEAD
// ToggleEntityStatus toggles the active/inactive status of an entity
// This should be added ONLY ONCE in your service.go file
// If you have this method declared twice, remove one of them
func (s *Service) ToggleEntityStatus(id int, isActive bool, userID uint, ip string) error{
	// Get existing entity
	existingEntity, err := s.Repo.GetEntityByID(id)
	if err != nil {
		auditDetails := map[string]interface{}{
			"temple_id": id,
			"error":     "Temple not found",
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_STATUS_TOGGLE_FAILED", auditDetails, ip, "failure")
		return err
	}

	// Update the status
	if err := s.Repo.UpdateEntityStatus(id, isActive); err != nil {
		auditDetails := map[string]interface{}{
			"temple_id":   id,
			"temple_name": existingEntity.Name,
			"new_status":  isActive,
			"error":       err.Error(),
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_STATUS_TOGGLE_FAILED", auditDetails, ip, "failure")
		return err
	}

	// Log successful status toggle
	statusText := "inactive"
	if isActive {
		statusText = "active"
	}

	auditDetails := map[string]interface{}{
		"temple_id":       id,
		"temple_name":     existingEntity.Name,
		"previous_status": existingEntity.IsActive,
		"new_status":      isActive,
		"status_text":     statusText,
	}
	entityID := uint(id)
	s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_STATUS_TOGGLED", auditDetails, ip, "success")

	return nil
}

// DeleteEntity - Super Admin → Delete temple
func (s *Service) DeleteEntity(id int, userID uint, ip string) error {
	existingEntity, err := s.Repo.GetEntityByID(id)
	if err != nil {
=======

// Super Admin → Delete temple
func (s *Service) DeleteEntity(id int, userID uint, ip string) error {
	// Get existing entity for audit log
	existingEntity, err := s.Repo.GetEntityByID(id)
	if err != nil {
		// 🔍 LOG FAILED TEMPLE DELETION ATTEMPT (NOT FOUND)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		auditDetails := map[string]interface{}{
			"temple_id": id,
			"error":     "Temple not found",
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETE_FAILED", auditDetails, ip, "failure")
<<<<<<< HEAD

=======
		
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		return err
	}

	if err := s.Repo.DeleteEntity(id); err != nil {
<<<<<<< HEAD
=======
		// 🔍 LOG FAILED TEMPLE DELETION (DB ERROR)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		auditDetails := map[string]interface{}{
			"temple_id":   id,
			"temple_name": existingEntity.Name,
			"error":       err.Error(),
		}
		entityID := uint(id)
		s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETE_FAILED", auditDetails, ip, "failure")
<<<<<<< HEAD

		return err
	}

=======
		
		return err
	}

	// 🔍 LOG SUCCESSFUL TEMPLE DELETION
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	auditDetails := map[string]interface{}{
		"temple_id":   id,
		"temple_name": existingEntity.Name,
		"temple_type": existingEntity.TempleType,
		"email":       existingEntity.Email,
		"city":        existingEntity.City,
		"state":       existingEntity.State,
	}
	entityID := uint(id)
	s.AuditService.LogAction(context.Background(), &userID, &entityID, "TEMPLE_DELETED", auditDetails, ip, "success")

	return nil
}

// ========== DEVOTEE MANAGEMENT ==========

<<<<<<< HEAD
// GetDevotees - Temple Admin → Get devotees for specific entity
=======
// Temple Admin → Get devotees for specific entity
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (s *Service) GetDevotees(entityID uint) ([]DevoteeDTO, error) {
	return s.Repo.GetDevoteesByEntityID(entityID)
}

<<<<<<< HEAD
// GetDevoteeStats - Temple Admin → Get devotee statistics for entity
=======
// Temple Admin → Get devotee statistics for entity
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
<<<<<<< HEAD
		Total    int64 `json:"total"`
		ThisWeek int64 `json:"this_week"`
	} `json:"upcoming_events"`
}

// GetDashboardSummary - Temple Admin → Dashboard Summary
func (s *Service) GetDashboardSummary(entityID uint) (DashboardSummary, error) {
	var summary DashboardSummary

=======
		Total     int64 `json:"total"`
		ThisWeek  int64 `json:"this_week"`
	} `json:"upcoming_events"`
}

// Temple Admin → Dashboard Summary
func (s *Service) GetDashboardSummary(entityID uint) (DashboardSummary, error) {
	var summary DashboardSummary

	// ============= DEVOTEES =============
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
=======
	// ============= SEVA BOOKINGS =============
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
=======
	// ============= DONATIONS =============
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	monthDonationAmount, percentChange, err := s.Repo.GetMonthDonationsWithChange(entityID)
	if err != nil {
		return summary, err
	}
	summary.MonthDonations.Amount = monthDonationAmount
	summary.MonthDonations.PercentChange = percentChange

<<<<<<< HEAD
=======
	// ============= UPCOMING EVENTS =============
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

// Helper function to track what fields were updated
func getUpdatedFields(old, new Entity) []string {
	var updatedFields []string

	if old.Name != new.Name {
		updatedFields = append(updatedFields, "name")
	}
<<<<<<< HEAD
	if (old.MainDeity == nil && new.MainDeity != nil) ||
		(old.MainDeity != nil && new.MainDeity == nil) ||
		(old.MainDeity != nil && new.MainDeity != nil && *old.MainDeity != *new.MainDeity) {
=======
	if (old.MainDeity == nil && new.MainDeity != nil) || 
	   (old.MainDeity != nil && new.MainDeity == nil) ||
	   (old.MainDeity != nil && new.MainDeity != nil && *old.MainDeity != *new.MainDeity) {
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		updatedFields = append(updatedFields, "main_deity")
	}
	if old.TempleType != new.TempleType {
		updatedFields = append(updatedFields, "temple_type")
	}
	if old.Email != new.Email {
		updatedFields = append(updatedFields, "email")
	}
	if old.Phone != new.Phone {
		updatedFields = append(updatedFields, "phone")
	}
	if old.Description != new.Description {
		updatedFields = append(updatedFields, "description")
	}
	if old.StreetAddress != new.StreetAddress {
		updatedFields = append(updatedFields, "street_address")
	}
	if old.City != new.City {
		updatedFields = append(updatedFields, "city")
	}
	if old.State != new.State {
		updatedFields = append(updatedFields, "state")
	}
	if old.District != new.District {
		updatedFields = append(updatedFields, "district")
	}
	if old.Pincode != new.Pincode {
		updatedFields = append(updatedFields, "pincode")
	}

	return updatedFields
<<<<<<< HEAD
}
=======
}
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
