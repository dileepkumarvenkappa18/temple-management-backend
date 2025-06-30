package seva

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateSeva(ctx context.Context, seva *Seva) error
	GetSevaByID(ctx context.Context, id uint) (*Seva, error)
	ListSevasByEntityID(ctx context.Context, entityID uint) ([]Seva, error)
	UpdateSeva(ctx context.Context, seva *Seva) error
	DeleteSeva(ctx context.Context, id uint) error

	BookSeva(ctx context.Context, booking *SevaBooking) error
	ListBookingsByUserID(ctx context.Context, userID uint) ([]SevaBooking, error)
	ListBookingsByEntityID(ctx context.Context, entityID uint) ([]SevaBooking, error)

	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error
	CancelBooking(ctx context.Context, bookingID uint, userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// Seva Methods
func (r *repository) CreateSeva(ctx context.Context, seva *Seva) error {
	return r.db.WithContext(ctx).Create(seva).Error
}

func (r *repository) GetSevaByID(ctx context.Context, id uint) (*Seva, error) {
	var seva Seva
	err := r.db.WithContext(ctx).First(&seva, id).Error
	return &seva, err
}

func (r *repository) ListSevasByEntityID(ctx context.Context, entityID uint) ([]Seva, error) {
	var sevas []Seva
	err := r.db.WithContext(ctx).Where("entity_id = ? AND is_active = true", entityID).Find(&sevas).Error
	return sevas, err
}

func (r *repository) UpdateSeva(ctx context.Context, seva *Seva) error {
	return r.db.WithContext(ctx).Save(seva).Error
}

func (r *repository) DeleteSeva(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Seva{}, id).Error
}

// Booking Methods
func (r *repository) BookSeva(ctx context.Context, booking *SevaBooking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *repository) ListBookingsByUserID(ctx context.Context, userID uint) ([]SevaBooking, error) {
	var bookings []SevaBooking
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *repository) ListBookingsByEntityID(ctx context.Context, entityID uint) ([]SevaBooking, error) {
	var bookings []SevaBooking
	err := r.db.WithContext(ctx).Where("entity_id = ?", entityID).Find(&bookings).Error
	return bookings, err
}

// ✅ New Method: Update Booking Status (used by temple admin)
func (r *repository) UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error {
	return r.db.WithContext(ctx).
		Model(&SevaBooking{}).
		Where("id = ?", bookingID).
		Update("status", newStatus).Error
}

// ✅ New Method: Cancel Booking (used by devotee)
func (r *repository) CancelBooking(ctx context.Context, bookingID uint, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&SevaBooking{}).
		Where("id = ? AND user_id = ?", bookingID, userID).
		Update("status", "cancelled").Error
}
