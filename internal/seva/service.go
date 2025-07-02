package seva

import (
	"context"
	"errors"
)

type Service interface {
	CreateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error
	UpdateSeva(ctx context.Context, seva *Seva, userRole string, entityID uint) error
	DeleteSeva(ctx context.Context, sevaID uint, userRole string) error
	GetSevasByEntity(ctx context.Context, entityID uint) ([]Seva, error)
	GetSevaByID(ctx context.Context, id uint) (*Seva, error)

	BookSeva(ctx context.Context, booking *SevaBooking, userRole string, userID uint, entityID uint) error
	GetBookingsForUser(ctx context.Context, userID uint) ([]SevaBooking, error)
	GetBookingsForEntity(ctx context.Context, entityID uint) ([]SevaBooking, error)

	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error
	CancelBooking(ctx context.Context, bookingID uint, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// templeadmin only
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

// Devotee only
func (s *service) CancelBooking(ctx context.Context, bookingID uint, userID uint) error {
	// Optional: add extra validation if needed in future
	return s.repo.CancelBooking(ctx, bookingID, userID)
}
