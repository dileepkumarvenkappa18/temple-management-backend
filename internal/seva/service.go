package seva

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/notification"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Service interface {
	// Seva Core
	CreateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error
	UpdateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error
	DeleteSeva(ctx context.Context, sevaID uint, accessContext middleware.AccessContext, ip string) error
	GetSevasByEntity(ctx context.Context, entityID uint) ([]Seva, error)
	GetSevaByID(ctx context.Context, id uint) (*Seva, error)

	// Enhanced seva listing with filters for temple admin
	GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error)

	// Booking Core
	BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint, ip string) error
	GetBookingsForUser(ctx context.Context, userID uint) ([]SevaBooking, error)
	GetBookingsForEntity(ctx context.Context, entityID uint) ([]SevaBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string, userID uint, ip string) error

	// Composite Booking Details
	GetDetailedBookingsForEntity(ctx context.Context, entityID uint) ([]DetailedBooking, error)
	CreateSevaBookingWithPayment(ctx context.Context, booking *SevaBooking, userID uint, entityID uint, ip string) error
	VerifySevaPayment(ctx context.Context, razorpayOrderID, razorpayPaymentID, razorpaySignature string, sevaID, userID uint, ip string) error

	// Filters, Search, Pagination
	SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error)

	// Counts
	GetBookingCountsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error)

	GetDetailedBookingsWithFilters(ctx context.Context, entityID uint, status, sevaType, startDate, endDate, search string, limit, offset int) ([]DetailedBooking, error)
	GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)
	GetBookingStatusCounts(ctx context.Context, entityID uint) (BookingStatusCounts, error)

	GetPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)

	// Get approved booking counts per seva
	GetApprovedBookingCountsPerSeva(ctx context.Context, entityID uint) (map[uint]int64, error)

	SetNotifService(n notification.Service)
}

type service struct {
	repo     Repository
	auditSvc auditlog.Service
	notifSvc notification.Service
}

func NewService(repo Repository, auditSvc auditlog.Service) Service {
	return &service{
		repo:     repo,
		auditSvc: auditSvc,
	}
}

func (s *service) SetNotifService(n notification.Service) {
	s.notifSvc = n
}

// isSuperAdmin returns true for roles that can manage any entity's sevas.
func isSuperAdmin(roleName string) bool {
	switch roleName {
	case "superadmin", "admin":
		return true
	}
	return false
}

// ─────────────────────────────────────────────
// CreateSeva
// FIX: Always resolve entity ID from accessContext first, then fall back to
//      seva.EntityID. Previously, if seva.EntityID == 0 AND accessContext
//      returned nil, the request failed intermittently (race / missing claim).
//      Now we resolve once, fail fast with a clear error, and stamp the seva.
// ─────────────────────────────────────────────
func (s *service) CreateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error {
	if !accessContext.CanWrite() {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_CREATE_FAILED", map[string]interface{}{
			"reason":    "write access denied",
			"seva_name": seva.Name,
		}, ip, "failure")
		return errors.New("write access denied")
	}

	// FIX: Trust seva.EntityID first (set by handler from route.params.id).
	// Only fall back to JWT access context if entity_id was not provided by caller.
	var resolvedEntityID uint
	if seva.EntityID != 0 {
		resolvedEntityID = seva.EntityID
	} else if ctxEntityID := accessContext.GetAccessibleEntityID(); ctxEntityID != nil {
		resolvedEntityID = *ctxEntityID
	} else {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_CREATE_FAILED", map[string]interface{}{
			"reason":    "no accessible entity",
			"seva_name": seva.Name,
		}, ip, "failure")
		return errors.New("no accessible entity")
	}

	// Stamp the seva with the resolved entity so the repo always has it.
	seva.EntityID = resolvedEntityID

	validStatuses := map[string]bool{"upcoming": true, "ongoing": true, "completed": true}
	if seva.Status == "" {
		seva.Status = "upcoming"
	} else if !validStatuses[seva.Status] {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, &resolvedEntityID, "SEVA_CREATE_FAILED", map[string]interface{}{
			"reason":         "invalid status",
			"seva_name":      seva.Name,
			"invalid_status": seva.Status,
		}, ip, "failure")
		return errors.New("invalid status. Must be 'upcoming', 'ongoing', or 'completed'")
	}

	seva.BookedSlots = 0
	seva.RemainingSlots = seva.AvailableSlots

	err := s.repo.CreateSeva(ctx, seva)
	if err != nil {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, &resolvedEntityID, "SEVA_CREATE_FAILED", map[string]interface{}{
			"seva_name": seva.Name,
			"seva_type": seva.SevaType,
			"status":    seva.Status,
			"error":     err.Error(),
		}, ip, "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &accessContext.UserID, &resolvedEntityID, "SEVA_CREATED", map[string]interface{}{
		"seva_id":         seva.ID,
		"seva_name":       seva.Name,
		"seva_type":       seva.SevaType,
		"price":           seva.Price,
		"status":          seva.Status,
		"available_slots": seva.AvailableSlots,
		"booked_slots":    seva.BookedSlots,
		"remaining_slots": seva.RemainingSlots,
		"role":            accessContext.RoleName,
	}, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppForEntityRoles(
			ctx,
			resolvedEntityID,
			[]string{"devotee", "volunteer"},
			"New Seva",
			seva.Name+" has been added",
			"seva",
		)
	}

	return nil
}

// UpdateSeva
// FIX: Added ownership check — verify the seva being updated actually belongs
//      to the caller's entity before proceeding. This mirrors the DeleteSeva
//      check and is why the handler was returning 403 (it had its own guard
//      because the service lacked one). With this in place, the handler's
//      pre-check is no longer needed and both paths are consistent.
// ────────────────────────────────────────────
func (s *service) UpdateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error {
	if !accessContext.CanWrite() {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_UPDATE_FAILED", map[string]interface{}{
			"reason":  "write access denied",
			"seva_id": seva.ID,
		}, ip, "failure")
		return errors.New("write access denied")
	}

	// FIX: Trust seva.EntityID first (set by handler from the existing seva's entity).
	// JWT GetAccessibleEntityID() returns 3, but the seva belongs to entity 2.
	// Ownership check must use the seva's actual entity, not the JWT entity.
	var entityID *uint
	if seva.EntityID != 0 {
		entityID = &seva.EntityID
	} else if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
		entityID = ctxID
	}

	if entityID == nil && !isSuperAdmin(accessContext.RoleName) {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_UPDATE_FAILED", map[string]interface{}{
			"reason":  "no accessible entity",
			"seva_id": seva.ID,
		}, ip, "failure")
		return errors.New("no accessible entity")
	}

	existing, err := s.repo.GetSevaByID(ctx, seva.ID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATE_FAILED", map[string]interface{}{
			"seva_id": seva.ID,
			"reason":  "seva not found",
			"error":   err.Error(),
		}, ip, "failure")
		return err
	}

	// FIX: Ownership check — verify caller can access the seva's entity via
	// canAccessSeva-style tenant check, not just a direct entity ID match.
	// A templeadmin with TenantID can manage sevas across all their entities.
	if !isSuperAdmin(accessContext.RoleName) {
		canAccess := false
		// Direct entity match
		if accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == existing.EntityID {
			canAccess = true
		}
		// Assigned entity match (standarduser/monitoringuser)
		if accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == existing.EntityID {
			canAccess = true
		}
		// Tenant ownership check — templeadmin can manage all entities under their tenant
		if !canAccess && accessContext.TenantID > 0 {
			tenantID, err := s.repo.GetTenantIDByEntityID(existing.EntityID)
			if err == nil && tenantID == accessContext.TenantID {
				canAccess = true
			}
		}
		if !canAccess {
			s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATE_FAILED", map[string]interface{}{
				"seva_id":          seva.ID,
				"seva_entity_id":   existing.EntityID,
				"caller_entity_id": *entityID,
				"reason":           "unauthorized: cannot update this seva",
			}, ip, "failure")
			return errors.New("unauthorized: cannot update this seva")
		}
	}

	if seva.Status != "" {
		validStatuses := map[string]bool{"upcoming": true, "ongoing": true, "completed": true}
		if !validStatuses[seva.Status] {
			s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATE_FAILED", map[string]interface{}{
				"reason":         "invalid status",
				"seva_id":        seva.ID,
				"invalid_status": seva.Status,
			}, ip, "failure")
			return errors.New("invalid status. Must be 'upcoming', 'ongoing', or 'completed'")
		}
	}

	if seva.BookedSlots == 0 {
		seva.BookedSlots = existing.BookedSlots
	}

	seva.RemainingSlots = seva.AvailableSlots - seva.BookedSlots
	if seva.RemainingSlots < 0 {
		seva.RemainingSlots = 0
	}

	err = s.repo.UpdateSeva(ctx, seva)
	if err != nil {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATE_FAILED", map[string]interface{}{
			"seva_id":   seva.ID,
			"seva_name": seva.Name,
			"error":     err.Error(),
		}, ip, "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATED", map[string]interface{}{
		"seva_id":         seva.ID,
		"seva_name":       seva.Name,
		"seva_type":       seva.SevaType,
		"price":           seva.Price,
		"status":          seva.Status,
		"available_slots": seva.AvailableSlots,
		"booked_slots":    seva.BookedSlots,
		"remaining_slots": seva.RemainingSlots,
		"role":            accessContext.RoleName,
	}, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppForEntityRoles(
			ctx,
			*entityID,
			[]string{"devotee", "volunteer"},
			"Seva Updated",
			seva.Name+" has been updated",
			"seva",
		)
	}

	return nil
}
// ─────────────────────────────────────────────
// DeleteSeva  (unchanged logic, only formatting)
// ─────────────────────────────────────────────
func (s *service) DeleteSeva(ctx context.Context, sevaID uint, accessContext middleware.AccessContext, ip string) error {
	if !accessContext.CanWrite() {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_DELETE_FAILED", map[string]interface{}{
			"reason":  "write access denied",
			"seva_id": sevaID,
		}, ip, "failure")
		return errors.New("write access denied")
	}

	entityID := accessContext.GetAccessibleEntityID()
	// Non-superadmin roles must be linked to an entity
	if entityID == nil && !isSuperAdmin(accessContext.RoleName) {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_DELETE_FAILED", map[string]interface{}{
			"reason":  "no accessible entity",
			"seva_id": sevaID,
		}, ip, "failure")
		return errors.New("no accessible entity")
	}

	seva, err := s.repo.GetSevaByID(ctx, sevaID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
			"seva_id": sevaID,
			"reason":  "seva not found",
			"error":   err.Error(),
		}, ip, "failure")
		return err
	}

	// ── Ownership check ────────────────────────────────────────────────────────
	// Superadmin can delete any seva; all other roles must own it.
	if !isSuperAdmin(accessContext.RoleName) && seva.EntityID != *entityID {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
			"seva_id":          sevaID,
			"seva_entity_id":   seva.EntityID,
			"caller_entity_id": *entityID,
			"reason":           "unauthorized: cannot delete this seva",
		}, ip, "failure")
		return errors.New("unauthorized: cannot delete this seva")
	}
	// ─────────────────────────────────────────────────────────────────────────

	bookings, err := s.repo.ListBookingsByEntityID(ctx, *entityID)
	if err == nil {
		for _, booking := range bookings {
			if booking.SevaID == sevaID {
				s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
					"seva_id":   sevaID,
					"seva_name": seva.Name,
					"reason":    "seva has existing bookings",
				}, ip, "failure")
				return errors.New("cannot delete seva with existing bookings")
			}
		}
	}

	err = s.repo.DeleteSeva(ctx, sevaID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
			"seva_id":   sevaID,
			"seva_name": seva.Name,
			"error":     err.Error(),
		}, ip, "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETED_PERMANENTLY", map[string]interface{}{
		"seva_id":   sevaID,
		"seva_name": seva.Name,
		"seva_type": seva.SevaType,
		"price":     seva.Price,
		"status":    seva.Status,
		"role":      accessContext.RoleName,
	}, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppForEntityRoles(
			ctx,
			*entityID,
			[]string{"devotee", "volunteer"},
			"Seva Deleted",
			seva.Name+" has been removed",
			"seva",
		)
	}

	return nil
}

func (s *service) GetSevasByEntity(ctx context.Context, entityID uint) ([]Seva, error) {
	return s.repo.ListSevasByEntityID(ctx, entityID)
}

func (s *service) GetSevaByID(ctx context.Context, id uint) (*Seva, error) {
	return s.repo.GetSevaByID(ctx, id)
}

func (s *service) GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error) {
	return s.repo.GetSevasWithFilters(ctx, entityID, sevaType, search, status, limit, offset)
}

// BookSeva — slot availability check uses RemainingSlots
func (s *service) BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint, ip string) error {
	if userRole != "devotee" {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"reason":  "unauthorized access",
			"seva_id": booking.SevaID,
		}, ip, "failure")
		return errors.New("unauthorized: only devotee can book sevas")
	}

	seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id": booking.SevaID,
			"reason":  "seva not found",
			"error":   err.Error(),
		}, ip, "failure")
		return err
	}

	if seva.Status != "upcoming" && seva.Status != "ongoing" {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id":     booking.SevaID,
			"seva_name":   seva.Name,
			"seva_status": seva.Status,
			"reason":      "seva is not bookable",
		}, ip, "failure")
		return errors.New("seva is not available for booking")
	}

	if seva.AvailableSlots > 0 && seva.RemainingSlots <= 0 {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id":         booking.SevaID,
			"seva_name":       seva.Name,
			"reason":          "no slots available",
			"available_slots": seva.AvailableSlots,
			"booked_slots":    seva.BookedSlots,
			"remaining_slots": seva.RemainingSlots,
		}, ip, "failure")
		return errors.New("no slots available for this seva")
	}

	booking.UserID = userID
	booking.EntityID = entityID
	booking.BookingTime = time.Now()
	booking.Status = "pending"

	err = s.repo.BookSeva(ctx, booking)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id":   booking.SevaID,
			"seva_name": seva.Name,
			"error":     err.Error(),
		}, ip, "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKED", map[string]interface{}{
		"booking_id":      booking.ID,
		"seva_id":         booking.SevaID,
		"seva_name":       seva.Name,
		"seva_type":       seva.SevaType,
		"seva_status":     seva.Status,
		"booking_status":  booking.Status,
		"available_slots": seva.AvailableSlots,
		"booked_slots":    seva.BookedSlots,
		"remaining_slots": seva.RemainingSlots,
	}, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppForEntityRoles(
			ctx,
			entityID,
			[]string{"templeadmin", "standarduser"},
			"New Seva Booking",
			"A new booking was created for "+seva.Name,
			"seva",
		)
		_ = s.notifSvc.CreateInAppNotification(
			ctx,
			userID,
			entityID,
			"Booking Created",
			"Your seva booking has been submitted",
			"seva",
		)
	}

	return nil
}

func (s *service) GetBookingsForUser(ctx context.Context, userID uint) ([]SevaBooking, error) {
	return s.repo.ListBookingsByUserID(ctx, userID)
}

func (s *service) GetBookingsForEntity(ctx context.Context, entityID uint) ([]SevaBooking, error) {
	return s.repo.ListBookingsByEntityID(ctx, entityID)
}

// UpdateBookingStatus — manages BookedSlots / RemainingSlots on status transitions
func (s *service) UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string, userID uint, ip string) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, nil, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
			"booking_id": bookingID,
			"new_status": newStatus,
			"reason":     "booking not found",
			"error":      err.Error(),
		}, ip, "failure")
		return err
	}

	seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
			"booking_id": bookingID,
			"new_status": newStatus,
			"reason":     "seva not found",
			"error":      err.Error(),
		}, ip, "failure")
		return err
	}

	oldStatus := booking.Status

	// Approving a booking (pending/rejected -> approved)
	if newStatus == "approved" && oldStatus != "approved" {
		if seva.RemainingSlots <= 0 {
			s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
				"booking_id":      bookingID,
				"new_status":      newStatus,
				"reason":          "no slots available",
				"available_slots": seva.AvailableSlots,
				"booked_slots":    seva.BookedSlots,
				"remaining_slots": seva.RemainingSlots,
			}, ip, "failure")
			return errors.New("no slots available for this seva")
		}

		if err := s.repo.IncrementBookedSlots(ctx, booking.SevaID); err != nil {
			s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
				"booking_id": bookingID,
				"new_status": newStatus,
				"reason":     "failed to update slots",
				"error":      err.Error(),
			}, ip, "failure")
			return fmt.Errorf("failed to update slots: %v", err)
		}
	}

	// Rejecting/Canceling an approved booking (approved -> rejected/pending)
	if oldStatus == "approved" && newStatus != "approved" {
		if err := s.repo.DecrementBookedSlots(ctx, booking.SevaID); err != nil {
			s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
				"booking_id": bookingID,
				"new_status": newStatus,
				"reason":     "failed to update slots",
				"error":      err.Error(),
			}, ip, "failure")
			return fmt.Errorf("failed to update slots: %v", err)
		}
	}

	err = s.repo.UpdateBookingStatus(ctx, bookingID, newStatus)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
			"booking_id": bookingID,
			"seva_id":    booking.SevaID,
			"new_status": newStatus,
			"error":      err.Error(),
		}, ip, "failure")
		return err
	}

	action := "SEVA_BOOKING_STATUS_UPDATED"
	switch newStatus {
	case "approved":
		action = "SEVA_BOOKING_APPROVED"
	case "rejected":
		action = "SEVA_BOOKING_REJECTED"
	}

	auditDetails := map[string]interface{}{
		"booking_id": bookingID,
		"seva_id":    booking.SevaID,
		"devotee_id": booking.UserID,
		"old_status": oldStatus,
		"new_status": newStatus,
	}

	if seva != nil {
		auditDetails["seva_name"] = seva.Name
		auditDetails["seva_type"] = seva.SevaType
		auditDetails["seva_status"] = seva.Status

		updatedSeva, _ := s.repo.GetSevaByID(ctx, booking.SevaID)
		if updatedSeva != nil {
			auditDetails["available_slots"] = updatedSeva.AvailableSlots
			auditDetails["booked_slots"] = updatedSeva.BookedSlots
			auditDetails["remaining_slots"] = updatedSeva.RemainingSlots
		}
	}

	s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, action, auditDetails, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppNotification(
			ctx,
			booking.UserID,
			booking.EntityID,
			"Seva Booking "+newStatus,
			"Your booking status is now "+newStatus,
			"seva",
		)
	}

	return nil
}

func (s *service) GetDetailedBookingsForEntity(ctx context.Context, entityID uint) ([]DetailedBooking, error) {
	return s.repo.ListBookingsWithDetails(ctx, entityID)
}

func (s *service) GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error) {
	return s.repo.GetBookingByID(ctx, bookingID)
}

func (s *service) SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) {
	return s.repo.SearchBookingsWithFilters(ctx, filter)
}

func (s *service) GetBookingCountsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error) {
	return s.repo.CountBookingsByStatus(ctx, entityID)
}

func (s *service) GetDetailedBookingsWithFilters(ctx context.Context, entityID uint, status, sevaType, startDate, endDate, search string, limit, offset int) ([]DetailedBooking, error) {
	filter := BookingFilter{
		EntityID:  entityID,
		Status:    status,
		SevaType:  sevaType,
		StartDate: startDate,
		EndDate:   endDate,
		Search:    search,
		Limit:     limit,
		Offset:    offset,
	}
	bookings, _, err := s.repo.SearchBookingsWithFilters(ctx, filter)
	return bookings, err
}

func (s *service) GetBookingStatusCounts(ctx context.Context, entityID uint) (BookingStatusCounts, error) {
	return s.repo.CountBookingsByStatus(ctx, entityID)
}

func (s *service) GetPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error) {
	return s.repo.ListPaginatedSevas(ctx, entityID, sevaType, search, limit, offset)
}

func (s *service) GetApprovedBookingCountsPerSeva(ctx context.Context, entityID uint) (map[uint]int64, error) {
	return s.repo.GetApprovedBookingsCountPerSeva(ctx, entityID)
}

// CreateSevaBookingWithPayment creates a pending booking with Razorpay order
func (s *service) CreateSevaBookingWithPayment(ctx context.Context, booking *SevaBooking, userID uint, entityID uint, ip string) error {
	seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id": booking.SevaID,
			"reason":  "seva not found",
			"error":   err.Error(),
		}, ip, "failure")
		return err
	}

	if seva.Status != "upcoming" && seva.Status != "ongoing" {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id":     booking.SevaID,
			"seva_status": seva.Status,
			"reason":      "seva is not bookable",
		}, ip, "failure")
		return errors.New("seva is not available for booking")
	}

	if seva.RemainingSlots <= 0 {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id":         booking.SevaID,
			"remaining_slots": seva.RemainingSlots,
			"reason":          "no slots available",
		}, ip, "failure")
		return errors.New("no slots available for this seva")
	}

	err = s.repo.BookSeva(ctx, booking)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
			"seva_id": booking.SevaID,
			"error":   err.Error(),
		}, ip, "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_PAYMENT_INITIATED", map[string]interface{}{
		"booking_id":        booking.ID,
		"seva_id":           booking.SevaID,
		"razorpay_order_id": booking.RazorpayOrderID,
		"amount":            booking.Amount,
		"remaining_slots":   seva.RemainingSlots,
	}, ip, "success")

	return nil
}

// VerifySevaPayment verifies payment and approves booking
func (s *service) VerifySevaPayment(ctx context.Context, razorpayOrderID, razorpayPaymentID, razorpaySignature string, sevaID, userID uint, ip string) error {
	booking, err := s.repo.GetBookingByOrderID(ctx, razorpayOrderID)
	if err != nil {
		s.auditSvc.LogAction(ctx, &userID, nil, "SEVA_PAYMENT_VERIFICATION_FAILED", map[string]interface{}{
			"razorpay_order_id": razorpayOrderID,
			"reason":            "booking not found",
			"error":             err.Error(),
		}, ip, "failure")
		return err
	}

	if booking.UserID != userID {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_PAYMENT_VERIFICATION_FAILED", map[string]interface{}{
			"booking_id": booking.ID,
			"reason":     "unauthorized access",
		}, ip, "failure")
		return errors.New("unauthorized access to booking")
	}

	seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
	if err != nil {
		return err
	}

	if seva.RemainingSlots <= 0 {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_PAYMENT_VERIFICATION_FAILED", map[string]interface{}{
			"booking_id":      booking.ID,
			"remaining_slots": seva.RemainingSlots,
			"reason":          "no slots available",
		}, ip, "failure")
		return errors.New("no slots available for this seva")
	}

	now := time.Now()
	booking.RazorpayPaymentID = razorpayPaymentID
	booking.RazorpaySignature = razorpaySignature
	booking.PaymentVerifiedAt = &now
	booking.Status = "approved"

	if err := s.repo.UpdateSevaBooking(ctx, booking); err != nil {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_PAYMENT_VERIFICATION_FAILED", map[string]interface{}{
			"booking_id": booking.ID,
			"error":      err.Error(),
		}, ip, "failure")
		return err
	}

	if err := s.repo.IncrementBookedSlots(ctx, booking.SevaID); err != nil {
		s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_SLOT_UPDATE_FAILED", map[string]interface{}{
			"booking_id": booking.ID,
			"seva_id":    booking.SevaID,
			"error":      err.Error(),
		}, ip, "failure")
		return fmt.Errorf("failed to update slots: %v", err)
	}

	s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_PAYMENT_VERIFIED", map[string]interface{}{
		"booking_id":          booking.ID,
		"seva_id":             booking.SevaID,
		"razorpay_payment_id": razorpayPaymentID,
		"amount":              booking.Amount,
		"status":              "approved",
	}, ip, "success")

	if s.notifSvc != nil {
		_ = s.notifSvc.CreateInAppNotification(
			ctx,
			userID,
			booking.EntityID,
			"Seva Booking Confirmed",
			fmt.Sprintf("Your seva booking has been confirmed. Payment ID: %s", razorpayPaymentID),
			"seva",
		)
	}

	return nil
}