package donation

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, donation *Donation) error
	GetByOrderID(ctx context.Context, orderID string) (*Donation, error)
	UpdatePaymentStatus(ctx context.Context, orderID string, status string, paymentID *string) error
	ListByUserID(ctx context.Context, userID uint) ([]Donation, error)
	ListByEntityID(ctx context.Context, entityID uint, page int, limit int, status string) ([]Donation, int64, error)
	ListByEntityWithFilters(ctx context.Context, entityID uint, page, limit int, status, from, to, dType, method, minAmount, maxAmount, search string) ([]Donation, int64, error)
	ListTopDonorsByEntity(ctx context.Context, entityID uint, limit int) ([]TopDonor, error)
	ListAllSuccessfulDonorsByEntity(ctx context.Context, entityID uint) ([]Donation, error)

	// Dashboard summaries
	GetDonationDashboard(ctx context.Context, entityID uint) (DonationDashboardResponse, error)
	GetDonationSummary(ctx context.Context, entityID uint, since time.Time) (DonationDashboardResponse, error)

	GetTopDonors(ctx context.Context, entityID uint) ([]TopDonor, error)
	GetAllDonors(ctx context.Context, entityID uint) ([]Donor, error)
	GetTotalAmountByEntityID(ctx context.Context, entityID uint) (float64, error)
	GetDonationCountByEntityID(ctx context.Context, entityID uint) (int, error)
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
	if status == StatusSuccess {
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

func (r *repository) ListByEntityID(ctx context.Context, entityID uint, page int, limit int, status string) ([]Donation, int64, error) {
	var donations []Donation
	var total int64

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&Donation{}).Where("entity_id = ?", entityID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&donations).Error

	return donations, total, err
}

func (r *repository) ListByEntityWithFilters(
	ctx context.Context,
	entityID uint,
	page, limit int,
	status, from, to, dType, method, minAmount, maxAmount, search string,
) ([]Donation, int64, error) {
	var donations []Donation
	var total int64

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).
		Model(&Donation{}).
		Where("entity_id = ?", entityID)

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if from != "" {
		db = db.Where("donated_at >= ?", from)
	}
	if to != "" {
		db = db.Where("donated_at <= ?", to)
	}
	if dType != "" {
		db = db.Where("donation_type = ?", dType)
	}
	if method != "" {
		db = db.Where("method = ?", method)
	}
	if minAmount != "" {
		if amt, err := strconv.ParseFloat(minAmount, 64); err == nil {
			db = db.Where("amount >= ?", amt)
		}
	}
	if maxAmount != "" {
		if amt, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			db = db.Where("amount <= ?", amt)
		}
	}
	if search != "" {
		keyword := "%" + search + "%"
		db = db.Where("order_id LIKE ? OR reference_id LIKE ?", keyword, keyword)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("donated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&donations).Error; err != nil {
		return nil, 0, err
	}

	return donations, total, nil
}

func (r *repository) ListTopDonorsByEntity(ctx context.Context, entityID uint, limit int) ([]TopDonor, error) {
	var donors []TopDonor
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select("u.full_name AS name, SUM(d.amount) AS amount").
		Joins("JOIN users u ON u.id = d.user_id").
		Where("d.entity_id = ? AND d.status = ?", entityID, StatusSuccess).
		Group("u.full_name").
		Order("amount DESC").
		Limit(limit).
		Scan(&donors).Error
	return donors, err
}

func (r *repository) ListAllSuccessfulDonorsByEntity(ctx context.Context, entityID uint) ([]Donation, error) {
	var donations []Donation
	err := r.db.WithContext(ctx).
		Where("entity_id = ? AND status = ?", entityID, StatusSuccess).
		Order("donated_at DESC").
		Find(&donations).Error
	return donations, err
}

func (r *repository) GetDonationDashboard(ctx context.Context, entityID uint) (DonationDashboardResponse, error) {
	var res DonationDashboardResponse

	// Summary fields
	err := r.db.WithContext(ctx).
		Model(&Donation{}).
		Select("COALESCE(SUM(amount), 0) AS total_donations, COUNT(DISTINCT user_id) AS total_donors, COALESCE(AVG(amount),0) AS average_donation").
		Where("entity_id = ? AND status = ?", entityID, StatusSuccess).
		Scan(&res).Error
	if err != nil {
		return res, err
	}

	// This month total
	beginningOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	err = r.db.WithContext(ctx).
		Model(&Donation{}).
		Select("COALESCE(SUM(amount),0)").
		Where("entity_id = ? AND status = ? AND donated_at >= ?", entityID, StatusSuccess, beginningOfMonth).
		Scan(&res.ThisMonth).Error
	if err != nil {
		return res, err
	}

	// 🔁 Manually populate RecentDonors (last 5)
	err = r.db.WithContext(ctx).
		Table("donations d").
		Select("u.full_name AS name, u.email, d.amount, d.donated_at AS date, d.method, d.status").
		Joins("JOIN users u ON u.id = d.user_id").
		Where("d.entity_id = ? AND d.status = ?", entityID, StatusSuccess).
		Order("d.donated_at DESC").
		Limit(5).
		Scan(&res.RecentDonors).Error

	return res, err
}

func (r *repository) GetDonationSummary(ctx context.Context, entityID uint, since time.Time) (DonationDashboardResponse, error) {
	var res DonationDashboardResponse
	err := r.db.WithContext(ctx).
		Model(&Donation{}).
		Select("COALESCE(SUM(amount),0) AS total_donations, COUNT(DISTINCT user_id) AS total_donors, COALESCE(AVG(amount),0) AS average_donation").
		Where("entity_id = ? AND status = ? AND donated_at >= ?", entityID, StatusSuccess, since).
		Scan(&res).Error
	return res, err
}

func (r *repository) GetTopDonors(ctx context.Context, entityID uint) ([]TopDonor, error) {
	return r.ListTopDonorsByEntity(ctx, entityID, 5)
}

func (r *repository) GetAllDonors(ctx context.Context, entityID uint) ([]Donor, error) {
	var donors []Donor
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select("u.full_name AS name, u.email as email, d.amount, d.donated_at as date, d.method, d.status").
		Joins("JOIN users u ON u.id = d.user_id").
		Where("d.entity_id = ? AND d.status = ?", entityID, StatusSuccess).
		Order("d.donated_at DESC").
		Scan(&donors).Error
	return donors, err
}

func (r *repository) GetTotalAmountByEntityID(ctx context.Context, entityID uint) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&Donation{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("entity_id = ? AND status = ?", entityID, StatusSuccess).
		Scan(&total).Error
	return total, err
}

func (r *repository) GetDonationCountByEntityID(ctx context.Context, entityID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Donation{}).
		Where("entity_id = ? AND status = ?", entityID, StatusSuccess).
		Count(&count).Error
	return int(count), err
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