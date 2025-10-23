package seva

import (
    "context"
    "errors"
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

    // Filters, Search, Pagination
    SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error)

    // Counts
    GetBookingCountsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error)

    GetDetailedBookingsWithFilters(ctx context.Context, entityID uint, status, sevaType, startDate, endDate, search string, limit, offset int) ([]DetailedBooking, error)
    GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)
    GetBookingStatusCounts(ctx context.Context, entityID uint) (BookingStatusCounts, error)

    GetPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)

    // NEW: inject notification service
    SetNotifService(n notification.Service)
}

type service struct {
    repo     Repository
    auditSvc auditlog.Service
    notifSvc notification.Service // NEW
}

func NewService(repo Repository, auditSvc auditlog.Service) Service {
    return &service{
        repo:     repo,
        auditSvc: auditSvc,
    }
}

// NEW
func (s *service) SetNotifService(n notification.Service) {
    s.notifSvc = n
}

// Updated to use access context
func (s *service) CreateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error {
    // Check write permissions
    if !accessContext.CanWrite() {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_CREATE_FAILED", map[string]interface{}{
            "reason":    "write access denied",
            "seva_name": seva.Name,
        }, ip, "failure")
        return errors.New("write access denied")
    }

    entityID := accessContext.GetAccessibleEntityID()
    if entityID == nil {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_CREATE_FAILED", map[string]interface{}{
            "reason":    "no accessible entity",
            "seva_name": seva.Name,
        }, ip, "failure")
        return errors.New("no accessible entity")
    }

    // Validate status
    validStatuses := map[string]bool{"upcoming": true, "ongoing": true, "completed": true}
    if seva.Status == "" {
        seva.Status = "upcoming" // Default status
    } else if !validStatuses[seva.Status] {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_CREATE_FAILED", map[string]interface{}{
            "reason":         "invalid status",
            "seva_name":      seva.Name,
            "invalid_status": seva.Status,
        }, ip, "failure")
        return errors.New("invalid status. Must be 'upcoming', 'ongoing', or 'completed'")
    }

    seva.EntityID = *entityID

    // Create seva
    err := s.repo.CreateSeva(ctx, seva)
    if err != nil {
        // Audit failed creation
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_CREATE_FAILED", map[string]interface{}{
            "seva_name": seva.Name,
            "seva_type": seva.SevaType,
            "status":    seva.Status,
            "error":     err.Error(),
        }, ip, "failure")
        return err
    }

    // Audit successful creation
    s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_CREATED", map[string]interface{}{
        "seva_id":   seva.ID,
        "seva_name": seva.Name,
        "seva_type": seva.SevaType,
        "price":     seva.Price,
        "status":    seva.Status,
        "role":      accessContext.RoleName,
    }, ip, "success")

    // In-app notify devotees and volunteers of the entity
    if s.notifSvc != nil {
        _ = s.notifSvc.CreateInAppForEntityRoles(
            ctx,
            *entityID,
            []string{"devotee", "volunteer"},
            "New Seva",
            seva.Name+" has been added",
            "seva",
        )
    }

    return nil
}

// FIXED UpdateSeva - Removed entityID overwrite
func (s *service) UpdateSeva(ctx context.Context, seva *Seva, accessContext middleware.AccessContext, ip string) error {
    // Check write permissions
    if !accessContext.CanWrite() {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_UPDATE_FAILED", map[string]interface{}{
            "reason":  "write access denied",
            "seva_id": seva.ID,
        }, ip, "failure")
        return errors.New("write access denied")
    }

    entityID := accessContext.GetAccessibleEntityID()
    if entityID == nil {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_UPDATE_FAILED", map[string]interface{}{
            "reason":  "no accessible entity",
            "seva_id": seva.ID,
        }, ip, "failure")
        return errors.New("no accessible entity")
    }

    // Validate status if provided
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

    // Update seva
    err := s.repo.UpdateSeva(ctx, seva)
    if err != nil {
        // Audit failed update
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATE_FAILED", map[string]interface{}{
            "seva_id":   seva.ID,
            "seva_name": seva.Name,
            "error":     err.Error(),
        }, ip, "failure")
        return err
    }

    // Audit successful update
    s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_UPDATED", map[string]interface{}{
        "seva_id":   seva.ID,
        "seva_name": seva.Name,
        "seva_type": seva.SevaType,
        "price":     seva.Price,
        "status":    seva.Status,
        "role":      accessContext.RoleName,
    }, ip, "success")

    // In-app notify devotees and volunteers of the entity
    if s.notifSvc != nil && accessContext.GetAccessibleEntityID() != nil {
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

// Updated to use access context with permanent delete
func (s *service) DeleteSeva(ctx context.Context, sevaID uint, accessContext middleware.AccessContext, ip string) error {
    // Check write permissions
    if !accessContext.CanWrite() {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, accessContext.GetAccessibleEntityID(), "SEVA_DELETE_FAILED", map[string]interface{}{
            "reason":  "write access denied",
            "seva_id": sevaID,
        }, ip, "failure")
        return errors.New("write access denied")
    }

    entityID := accessContext.GetAccessibleEntityID()
    if entityID == nil {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, nil, "SEVA_DELETE_FAILED", map[string]interface{}{
            "reason":  "no accessible entity",
            "seva_id": sevaID,
        }, ip, "failure")
        return errors.New("no accessible entity")
    }

    // Get seva details for audit logging before deletion
    seva, err := s.repo.GetSevaByID(ctx, sevaID)
    if err != nil {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
            "seva_id": sevaID,
            "reason":  "seva not found",
            "error":   err.Error(),
        }, ip, "failure")
        return err
    }

    // Verify seva belongs to accessible entity
    if seva.EntityID != *entityID {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
            "seva_id": sevaID,
            "reason":  "access denied to this seva",
        }, ip, "failure")
        return errors.New("access denied to this seva")
    }

    // Check if there are any bookings for this seva
    bookings, err := s.repo.ListBookingsByEntityID(ctx, *entityID)
    if err == nil {
        hasBookings := false
        for _, booking := range bookings {
            if booking.SevaID == sevaID {
                hasBookings = true
                break
            }
        }

        if hasBookings {
            s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
                "seva_id":   sevaID,
                "seva_name": seva.Name,
                "reason":    "seva has existing bookings",
            }, ip, "failure")
            return errors.New("cannot delete seva with existing bookings")
        }
    }

    // Perform permanent delete
    err = s.repo.DeleteSeva(ctx, sevaID)
    if err != nil {
        s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETE_FAILED", map[string]interface{}{
            "seva_id":   sevaID,
            "seva_name": seva.Name,
            "error":     err.Error(),
        }, ip, "failure")
        return err
    }

    // Audit successful deletion
    s.auditSvc.LogAction(ctx, &accessContext.UserID, entityID, "SEVA_DELETED_PERMANENTLY", map[string]interface{}{
        "seva_id":   sevaID,
        "seva_name": seva.Name,
        "seva_type": seva.SevaType,
        "price":     seva.Price,
        "status":    seva.Status,
        "role":      accessContext.RoleName,
    }, ip, "success")

    // In-app notify devotees and volunteers of the entity
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

// Enhanced seva listing with filters for temple admin
func (s *service) GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error) {
    return s.repo.GetSevasWithFilters(ctx, entityID, sevaType, search, status, limit, offset)
}

// Devotee only - keep unchanged but validate seva status
func (s *service) BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint, ip string) error {
    if userRole != "devotee" {
        // Audit failed attempt
        s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
            "reason":  "unauthorized access",
            "seva_id": booking.SevaID,
        }, ip, "failure")
        return errors.New("unauthorized: only devotee can book sevas")
    }

    // Validate Seva exists and is bookable
    seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
    if err != nil {
        // Audit failed booking
        s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
            "seva_id": booking.SevaID,
            "reason":  "seva not found",
            "error":   err.Error(),
        }, ip, "failure")
        return err
    }

    // Check if seva is bookable (only upcoming and ongoing sevas can be booked)
    if seva.Status != "upcoming" && seva.Status != "ongoing" {
        s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
            "seva_id":     booking.SevaID,
            "seva_name":   seva.Name,
            "seva_status": seva.Status,
            "reason":      "seva is not bookable",
        }, ip, "failure")
        return errors.New("seva is not available for booking")
    }

    booking.UserID = userID
    booking.EntityID = entityID
    booking.BookingTime = time.Now()
    booking.Status = "pending"

    // Create booking
    err = s.repo.BookSeva(ctx, booking)
    if err != nil {
        // Audit failed booking
        s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKING_FAILED", map[string]interface{}{
            "seva_id":   booking.SevaID,
            "seva_name": seva.Name,
            "error":     err.Error(),
        }, ip, "failure")
        return err
    }

    // Audit successful booking
    s.auditSvc.LogAction(ctx, &userID, &entityID, "SEVA_BOOKED", map[string]interface{}{
        "booking_id":     booking.ID,
        "seva_id":        booking.SevaID,
        "seva_name":      seva.Name,
        "seva_type":      seva.SevaType,
        "seva_status":    seva.Status,
        "booking_status": booking.Status,
    }, ip, "success")

    // In-app notifications:
    // - inform temple admins and standard users in this entity
    // - inform the booking devotee
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

// Temple Admin only - keep existing logic since it's used through access context
func (s *service) UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string, userID uint, ip string) error {
    // Get booking details for audit
    booking, err := s.repo.GetBookingByID(ctx, bookingID)
    if err != nil {
        // Audit failed attempt
        s.auditSvc.LogAction(ctx, &userID, nil, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
            "booking_id": bookingID,
            "new_status": newStatus,
            "reason":     "booking not found",
            "error":      err.Error(),
        }, ip, "failure")
        return err
    }

    // Get seva details for better audit logging
    seva, _ := s.repo.GetSevaByID(ctx, booking.SevaID)

    // Update status
    err = s.repo.UpdateBookingStatus(ctx, bookingID, newStatus)
    if err != nil {
        // Audit failed update
        s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, "SEVA_BOOKING_STATUS_UPDATE_FAILED", map[string]interface{}{
            "booking_id": bookingID,
            "seva_id":    booking.SevaID,
            "new_status": newStatus,
            "error":      err.Error(),
        }, ip, "failure")
        return err
    }

    // Audit successful status update with specific action
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
        "old_status": booking.Status,
        "new_status": newStatus,
    }

    if seva != nil {
        auditDetails["seva_name"] = seva.Name
        auditDetails["seva_type"] = seva.SevaType
        auditDetails["seva_status"] = seva.Status
    }

    s.auditSvc.LogAction(ctx, &userID, &booking.EntityID, action, auditDetails, ip, "success")

    // Notify the devotee about status change
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

// Temple admin only: Full booking table with names, types, etc.
func (s *service) GetDetailedBookingsForEntity(ctx context.Context, entityID uint) ([]DetailedBooking, error) {
    return s.repo.ListBookingsWithDetails(ctx, entityID)
}

// GetBookingByID: Public or Temple admin
func (s *service) GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error) {
    return s.repo.GetBookingByID(ctx, bookingID)
}

// SearchBookings: Temple admin with filters (pagination, search, etc.)
func (s *service) SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) {
    return s.repo.SearchBookingsWithFilters(ctx, filter)
}

// Count booking statuses: Temple admin dashboard metrics
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
