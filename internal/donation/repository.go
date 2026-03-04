package donation

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Basic CRUD operations
	Create(ctx context.Context, donation *Donation) error
	GetByOrderID(ctx context.Context, orderID string) (*Donation, error)
	GetByIDWithUser(ctx context.Context, donationID uint) (*DonationWithUser, error)
	UpdatePaymentDetails(ctx context.Context, orderID string, params UpdatePaymentDetailsParams) error

	// Data retrieval with filtering
	ListByUserID(ctx context.Context, userID uint) ([]DonationWithUser, error)
	ListByUserIDAndEntity(ctx context.Context, userID uint, entityID uint) ([]DonationWithUser, error)
	ListWithFilters(ctx context.Context, filters DonationFilters) ([]DonationWithUser, int, error)

	// Analytics queries
	GetTotalStats(ctx context.Context, entityID uint) (*StatsResult, error)
	GetStatsInDateRange(ctx context.Context, entityID uint, from, to time.Time) (*StatsResult, error)
	GetUniqueDonorCount(ctx context.Context, entityID uint) (int, error)
	GetTopDonors(ctx context.Context, entityID uint, limit int) ([]TopDonor, error)
	GetDonationTrends(ctx context.Context, entityID uint, days int) ([]TrendData, error)
	GetDonationsByType(ctx context.Context, entityID uint) ([]TypeData, error)
	GetDonationsByMethod(ctx context.Context, entityID uint) ([]MethodData, error)

	// Recent donations
	GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error)
	GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error)
	GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int) ([]RecentDonation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==============================
// Basic CRUD Operations
// ==============================

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

const donationSelectFields = `
	d.id, d.user_id, d.entity_id, d.amount, d.donation_type, d.reference_id,
	d.method, d.status, d.order_id, d.payment_id, d.note, d.donated_at,
	d.created_at, d.updated_at,
	COALESCE(d.account_holder_name, '') as account_holder_name,
	COALESCE(d.account_number, '') as account_number,
	COALESCE(d.account_type, '') as account_type,
	COALESCE(d.ifsc_code, '') as ifsc_code,
	COALESCE(d.upi_id, '') as upi_id,
	COALESCE(NULLIF(u.full_name, ''), u.email, 'Anonymous') as user_name,
	COALESCE(u.email, '') as user_email,
	COALESCE(e.name, '') as entity_name
`

func (r *repository) GetByIDWithUser(ctx context.Context, donationID uint) (*DonationWithUser, error) {
	var result DonationWithUser
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(donationSelectFields).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.id = ?", donationID).
		First(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) UpdatePaymentDetails(ctx context.Context, orderID string, params UpdatePaymentDetailsParams) error {
	updates := map[string]interface{}{
		"status":     params.Status,
		"payment_id": params.PaymentID,
		"method":     params.Method,
		"amount":     params.Amount,
	}

	if params.DonatedAt != nil {
		updates["donated_at"] = params.DonatedAt
	}

	// Always update payee fields (even if empty string — clears stale data)
	updates["account_holder_name"] = params.AccountHolderName
	updates["account_number"] = params.AccountNumber
	updates["account_type"] = params.AccountType
	updates["ifsc_code"] = params.IFSCCode
	updates["upi_id"] = params.UPIID // FIX: string now, no pointer needed

	return r.db.WithContext(ctx).
		Model(&Donation{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

// ==============================
// Data Retrieval with Filtering
// ==============================

func (r *repository) ListByUserID(ctx context.Context, userID uint) ([]DonationWithUser, error) {
	var donations []DonationWithUser
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(donationSelectFields).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.user_id = ?", userID).
		Order("d.created_at DESC").
		Find(&donations).Error

	return donations, err
}

func (r *repository) ListByUserIDAndEntity(ctx context.Context, userID uint, entityID uint) ([]DonationWithUser, error) {
	var donations []DonationWithUser
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(donationSelectFields).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.user_id = ? AND d.entity_id = ?", userID, entityID).
		Order("d.created_at DESC").
		Find(&donations).Error

	return donations, err
}

func (r *repository) ListWithFilters(ctx context.Context, filters DonationFilters) ([]DonationWithUser, int, error) {
	var donations []DonationWithUser
	var total int64

	query := r.db.WithContext(ctx).
		Table("donations d").
		Select(donationSelectFields).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id")

	// Role-based filtering
	if filters.EntityID != 0 && filters.UserID != 0 {
		query = query.Where("d.entity_id = ? AND d.user_id = ?", filters.EntityID, filters.UserID)
	} else if filters.EntityID != 0 {
		query = query.Where("d.entity_id = ?", filters.EntityID)
	} else if filters.UserID != 0 {
		query = query.Where("d.user_id = ?", filters.UserID)
	}

	query = r.applyFilters(query, filters)

	// Count query
	countQuery := r.db.WithContext(ctx).
		Table("donations d").
		Joins("LEFT JOIN users u ON d.user_id = u.id")

	if filters.EntityID != 0 && filters.UserID != 0 {
		countQuery = countQuery.Where("d.entity_id = ? AND d.user_id = ?", filters.EntityID, filters.UserID)
	} else if filters.EntityID != 0 {
		countQuery = countQuery.Where("d.entity_id = ?", filters.EntityID)
	} else if filters.UserID != 0 {
		countQuery = countQuery.Where("d.user_id = ?", filters.UserID)
	}

	countQuery = r.applyFilters(countQuery, filters)
	countQuery.Count(&total)

	if filters.Page > 0 && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset).Limit(filters.Limit)
	}

	err := query.Order("d.created_at DESC").Find(&donations).Error
	return donations, int(total), err
}

func (r *repository) applyFilters(query *gorm.DB, filters DonationFilters) *gorm.DB {
	if filters.Status != "" && filters.Status != "all" {
		query = query.Where("LOWER(d.status) = LOWER(?)", filters.Status)
	}
	if filters.Type != "" && filters.Type != "all" {
		query = query.Where("LOWER(d.donation_type) = LOWER(?)", filters.Type)
	}
	if filters.Method != "" && filters.Method != "all" {
		query = query.Where("LOWER(d.method) = LOWER(?)", filters.Method)
	}
	if filters.From != nil {
		query = query.Where("d.created_at >= ?", filters.From)
	}
	if filters.To != nil {
		query = query.Where("d.created_at <= ?", filters.To)
	}
	if filters.MinAmount != nil {
		query = query.Where("d.amount >= ?", *filters.MinAmount)
	}
	if filters.MaxAmount != nil {
		query = query.Where("d.amount <= ?", *filters.MaxAmount)
	}
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where(`
			COALESCE(NULLIF(u.full_name, ''), u.email, 'Anonymous') ILIKE ? OR 
			u.email ILIKE ? OR 
			d.payment_id ILIKE ? OR 
			d.order_id ILIKE ?
		`, searchTerm, searchTerm, searchTerm, searchTerm)
	}
	return query
}

// ==============================
// Analytics Queries
// ==============================

func (r *repository) GetTotalStats(ctx context.Context, entityID uint) (*StatsResult, error) {
	var result StatsResult
	err := r.db.WithContext(ctx).
		Table("donations").
		Select(`
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN amount ELSE 0 END), 0) as amount,
			COUNT(*) as count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN 1 ELSE 0 END), 0) as completed_count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'failed' THEN 1 ELSE 0 END), 0) as failed_count
		`).
		Where("entity_id = ?", entityID).
		Scan(&result).Error
	return &result, err
}

func (r *repository) GetStatsInDateRange(ctx context.Context, entityID uint, from, to time.Time) (*StatsResult, error) {
	var result StatsResult
	err := r.db.WithContext(ctx).
		Table("donations").
		Select(`
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN amount ELSE 0 END), 0) as amount,
			COUNT(*) as count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN 1 ELSE 0 END), 0) as completed_count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'failed' THEN 1 ELSE 0 END), 0) as failed_count
		`).
		Where("entity_id = ? AND created_at >= ? AND created_at <= ?", entityID, from, to).
		Scan(&result).Error
	return &result, err
}

func (r *repository) GetUniqueDonorCount(ctx context.Context, entityID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("donations").
		Where("entity_id = ? AND LOWER(status) = 'success'", entityID).
		Distinct("user_id").
		Count(&count).Error
	return int(count), err
}

func (r *repository) GetTopDonors(ctx context.Context, entityID uint, limit int) ([]TopDonor, error) {
	var donors []TopDonor
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(`
			COALESCE(NULLIF(u.full_name, ''), u.email, 'Anonymous') as name,
			COALESCE(u.email, '') as email,
			SUM(d.amount) as total_amount,
			COUNT(d.id) as donation_count
		`).
		Joins("JOIN users u ON d.user_id = u.id").
		Where("d.entity_id = ? AND LOWER(d.status) = 'success'", entityID).
		Group("u.id, u.full_name, u.email").
		Order("total_amount DESC").
		Limit(limit).
		Scan(&donors).Error
	return donors, err
}

func (r *repository) GetDonationTrends(ctx context.Context, entityID uint, days int) ([]TrendData, error) {
	var trends []TrendData
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Table("donations").
		Select(`
			DATE(created_at) as date,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN amount ELSE 0 END), 0) as amount,
			COUNT(*) as count
		`).
		Where("entity_id = ? AND created_at >= ? AND created_at <= ?", entityID, startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends).Error
	return trends, err
}

func (r *repository) GetDonationsByType(ctx context.Context, entityID uint) ([]TypeData, error) {
	var typeData []TypeData
	err := r.db.WithContext(ctx).
		Table("donations").
		Select(`
			donation_type as type,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN amount ELSE 0 END), 0) as amount,
			COUNT(*) as count
		`).
		Where("entity_id = ?", entityID).
		Group("donation_type").
		Order("amount DESC").
		Scan(&typeData).Error
	return typeData, err
}

func (r *repository) GetDonationsByMethod(ctx context.Context, entityID uint) ([]MethodData, error) {
	var methodData []MethodData
	err := r.db.WithContext(ctx).
		Table("donations").
		Select(`
			method,
			COALESCE(SUM(CASE WHEN LOWER(status) = 'success' THEN amount ELSE 0 END), 0) as amount,
			COUNT(*) as count
		`).
		Where("entity_id = ? AND LOWER(status) = 'success'", entityID).
		Group("method").
		Order("amount DESC").
		Scan(&methodData).Error
	return methodData, err
}

// ==============================
// Recent Donations
// ==============================

const recentDonationSelect = `
	d.amount, d.donation_type, d.method, d.status,
	COALESCE(d.donated_at, d.created_at) as donated_at,
	COALESCE(NULLIF(u.full_name, ''), u.email, 'Anonymous') as user_name,
	COALESCE(e.name, '') as entity_name
`

func (r *repository) GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error) {
	var recent []RecentDonation
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(recentDonationSelect).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.user_id = ?", userID).
		Order("COALESCE(d.donated_at, d.created_at) DESC").
		Limit(limit).
		Scan(&recent).Error
	return recent, err
}

func (r *repository) GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error) {
	var recent []RecentDonation
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(recentDonationSelect).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.user_id = ? AND d.entity_id = ?", userID, entityID).
		Order("COALESCE(d.donated_at, d.created_at) DESC").
		Limit(limit).
		Scan(&recent).Error
	return recent, err
}

func (r *repository) GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int) ([]RecentDonation, error) {
	var recent []RecentDonation
	err := r.db.WithContext(ctx).
		Table("donations d").
		Select(recentDonationSelect).
		Joins("LEFT JOIN users u ON d.user_id = u.id").
		Joins("LEFT JOIN entities e ON d.entity_id = e.id").
		Where("d.entity_id = ?", entityID).
		Order("COALESCE(d.donated_at, d.created_at) DESC").
		Limit(limit).
		Scan(&recent).Error
	return recent, err
}
