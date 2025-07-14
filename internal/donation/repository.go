package donation

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, donation *Donation) error
	GetByOrderID(ctx context.Context, orderID string) (*Donation, error)
	UpdatePaymentStatus(ctx context.Context, orderID string, status string, paymentID *string) error
	ListByUserID(ctx context.Context, userID uint) ([]Donation, error)
	ListByEntityID(ctx context.Context, entityID uint, page int, limit int, status string) ([]Donation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, donation *Donation) error {
	return r.db.WithContext(ctx).Create(donation).Error
}

func (r *repository) GetByOrderID(ctx context.Context, orderID string) (*Donation, error) {
	var donation Donation
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&donation).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, orderID string, status string, paymentID *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"payment_id": paymentID,
	}

	// Record donation time only on successful payment
	if status == "SUCCESS" {
		updates["donated_at"] = time.Now()
	}

	return r.db.WithContext(ctx).
		Model(&Donation{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

func (r *repository) ListByUserID(ctx context.Context, userID uint) ([]Donation, error) {
	var donations []Donation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&donations).Error
	return donations, err
}

func (r *repository) ListByEntityID(ctx context.Context, entityID uint, page int, limit int, status string) ([]Donation, error) {
	var donations []Donation

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // default pagination fallback
	}
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).
		Where("entity_id = ?", entityID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&donations).Error

	return donations, err
}







// package donation

// import (
// 	"context"

// 	"gorm.io/gorm"
// )

// type Repository interface {
// 	Create(ctx context.Context, donation *Donation) error
// 	GetByOrderID(ctx context.Context, orderID string) (*Donation, error)
// 	UpdatePaymentStatus(ctx context.Context, orderID string, status string, paymentID *string) error
// 	ListByUserID(ctx context.Context, userID uint) ([]Donation, error)
// 	ListByEntityID(ctx context.Context, entityID uint) ([]Donation, error)
// }

// type repository struct {
// 	db *gorm.DB
// }

// func NewRepository(db *gorm.DB) Repository {
// 	return &repository{db: db}
// }

// func (r *repository) Create(ctx context.Context, donation *Donation) error {
// 	return r.db.WithContext(ctx).Create(donation).Error
// }

// func (r *repository) GetByOrderID(ctx context.Context, orderID string) (*Donation, error) {
// 	var donation Donation
// 	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&donation).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &donation, nil
// }

// func (r *repository) UpdatePaymentStatus(ctx context.Context, orderID string, status string, paymentID *string) error {
// 	return r.db.WithContext(ctx).
// 		Model(&Donation{}).
// 		Where("order_id = ?", orderID).
// 		Updates(map[string]interface{}{
// 			"status":     status,
// 			"payment_id": paymentID,
// 		}).Error
// }

// func (r *repository) ListByUserID(ctx context.Context, userID uint) ([]Donation, error) {
// 	var donations []Donation
// 	err := r.db.WithContext(ctx).
// 		Where("user_id = ?", userID).
// 		Order("created_at DESC").
// 		Find(&donations).Error
// 	return donations, err
// }

// func (r *repository) ListByEntityID(ctx context.Context, entityID uint) ([]Donation, error) {
// 	var donations []Donation
// 	err := r.db.WithContext(ctx).
// 		Where("entity_id = ?", entityID).
// 		Order("created_at DESC").
// 		Find(&donations).Error
// 	return donations, err
// }