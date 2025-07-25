package seva

import (
	"context"
	"errors"
)

type Service interface {
	// Seva Core
	CreateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error
	UpdateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error
	DeleteSeva(ctx context.Context, sevaID uint, userRole string) error
	GetSevasByEntity(ctx context.Context, entityID uint) ([]Seva, error)
	GetSevaByID(ctx context.Context, id uint) (*Seva, error)

	// Booking Core
	BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint) error
	GetBookingsForUser(ctx context.Context, userID uint) ([]SevaBooking, error)
	GetBookingsForEntity(ctx context.Context, entityID uint) ([]SevaBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error


	// Composite Booking Details
	GetDetailedBookingsForEntity(ctx context.Context, entityID uint) ([]DetailedBooking, error)

	// Filters, Search, Pagination
	SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) // ✅ New

	// Counts
	GetBookingCountsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error) // ✅ New

	GetDetailedBookingsWithFilters(ctx context.Context, entityID uint, status, sevaType, startDate, endDate, search string, limit, offset int) ([]DetailedBooking, error)
	GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)
	GetBookingStatusCounts(ctx context.Context, entityID uint) (BookingStatusCounts, error)

	GetPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)

}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// Templeadmin only
func (s *service) CreateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error {
	if userRole != "templeadmin" {
		return errors.New("unauthorized: only templeadmin can create sevas")
	}
	seva.EntityID = entityID
	return s.repo.CreateSeva(ctx, seva)
}

func (s *service) UpdateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error {
	if userRole != "templeadmin" {
		return errors.New("unauthorized: only templeadmin can update sevas")
	}
	seva.EntityID = entityID
	return s.repo.UpdateSeva(ctx, seva)
}

func (s *service) DeleteSeva(ctx context.Context, sevaID uint, userRole string) error {
	if userRole != "templeadmin" {
		return errors.New("unauthorized: only templeadmin can delete sevas")
	}
	return s.repo.DeleteSeva(ctx, sevaID)
}

func (s *service) GetSevasByEntity(ctx context.Context, entityID uint) ([]Seva, error) {
	return s.repo.ListSevasByEntityID(ctx, entityID)
}

func (s *service) GetSevaByID(ctx context.Context, id uint) (*Seva, error) {
	return s.repo.GetSevaByID(ctx, id)
}

// Devotee only
func (s *service) BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint) error {
	if userRole != "devotee" {
		return errors.New("unauthorized: only devotee can book sevas")
	}

	seva, err := s.repo.GetSevaByID(ctx, booking.SevaID)
	if err != nil {
		return err
	}

	timeSlot := booking.BookingTime.Format("15:04")
	count, err := s.repo.CountBookingsForSlot(ctx, booking.SevaID, booking.BookingDate, timeSlot)
	if err != nil {
		return err
	}
	if int(count) >= seva.MaxBookingsPerDay {
		return errors.New("booking limit reached for selected time slot")
	}

	booking.UserID = userID
	booking.EntityID = entityID
	return s.repo.BookSeva(ctx, booking)
}

func (s *service) GetBookingsForUser(ctx context.Context, userID uint) ([]SevaBooking, error) {
	return s.repo.ListBookingsByUserID(ctx, userID)
}

func (s *service) GetBookingsForEntity(ctx context.Context, entityID uint) ([]SevaBooking, error) {
	return s.repo.ListBookingsByEntityID(ctx, entityID)
}

// Temple Admin only
func (s *service) UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error {
	return s.repo.UpdateBookingStatus(ctx, bookingID, newStatus)
}



// 🔄 Templeadmin only: Add/Replace Availability
// func (s *service) SetAvailabilityForSeva(ctx context.Context, sevaID uint, slots []SevaAvailability, userRole string) error {
// 	if userRole != "templeadmin" {
// 		return errors.New("unauthorized: only templeadmin can manage availability")
// 	}

// 	// Remove previous availability
// 	if err := s.repo.DeleteAvailabilityBySevaID(ctx, sevaID); err != nil {
// 		return err
// 	}

// 	// Save new availability
// 	for _, slot := range slots {
// 		slot.SevaID = sevaID
// 		if err := s.repo.CreateAvailability(ctx, &slot); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// 🔄 Public: View availability for seva
// func (s *service) GetAvailabilityForSeva(ctx context.Context, sevaID uint) ([]SevaAvailability, error) {
// 	return s.repo.GetAvailabilityBySevaID(ctx, sevaID)
// }

// 🔄 Templeadmin only: Full booking table with names, types, etc.
func (s *service) GetDetailedBookingsForEntity(ctx context.Context, entityID uint) ([]DetailedBooking, error) {
	return s.repo.ListBookingsWithDetails(ctx, entityID)
}

// ✅ GetBookingByID: Public or Templeadmin
func (s *service) GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error) {
	return s.repo.GetBookingByID(ctx, bookingID)
}

// ✅ SearchBookings: Templeadmin with filters (pagination, search, etc.)
func (s *service) SearchBookings(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) {
	return s.repo.SearchBookingsWithFilters(ctx, filter)
}

// ✅ Count booking statuses: Templeadmin dashboard metrics
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
