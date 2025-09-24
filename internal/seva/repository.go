package seva

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Seva core
	CreateSeva(ctx context.Context, seva *Seva) error
	GetSevaByID(ctx context.Context, id uint) (*Seva, error)
	ListSevasByEntityID(ctx context.Context, entityID uint) ([]Seva, error)
	UpdateSeva(ctx context.Context, seva *Seva) error
	DeleteSeva(ctx context.Context, id uint) error

	// Enhanced seva listing with filters
	GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error)

	// Booking core
	BookSeva(ctx context.Context, booking *SevaBooking) error
	ListBookingsByUserID(ctx context.Context, userID uint) ([]SevaBooking, error)
	ListBookingsByEntityID(ctx context.Context, entityID uint) ([]SevaBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error

	// Booking limits
	CountBookingsForSlot(ctx context.Context, sevaID uint, date time.Time, slot string) (int64, error)

	// Composite list with Seva + User info
	ListBookingsWithDetails(ctx context.Context, entityID uint) ([]DetailedBooking, error)

	GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)
	SearchBookingsWithFilters(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error)
	CountBookingsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error)

	ListPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)
}

// BookingFilter and DetailedBooking structs are now defined in model.go

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// -----------------------------------------
// Seva Core
// -----------------------------------------
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
	err := r.db.WithContext(ctx).
		Where("entity_id = ?", entityID).
		Find(&sevas).Error
	return sevas, err
}

func (r *repository) ListPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error) {
	var sevas []Seva

	query := r.db.WithContext(ctx).
		Model(&Seva{}).
		Where("entity_id = ? AND status IN (?)", entityID, []string{"upcoming", "ongoing"}) // Only show bookable sevas to devotees

	if sevaType != "" {
		query = query.Where("seva_type = ?", sevaType)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&sevas).Error
	return sevas, err
}

// Enhanced seva listing with filters for temple admin
func (r *repository) GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error) {
	var sevas []Seva
	var total int64

	query := r.db.WithContext(ctx).
		Model(&Seva{}).
		Where("entity_id = ?", entityID)

	// Apply filters
	if sevaType != "" {
		query = query.Where("seva_type = ?", sevaType)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&sevas).Error
	return sevas, total, err
}

func (r *repository) UpdateSeva(ctx context.Context, seva *Seva) error {
	return r.db.WithContext(ctx).Save(seva).Error
}

// Permanent delete - removes the record completely
func (r *repository) DeleteSeva(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&Seva{}, id).Error
}

// -----------------------------------------
// Booking Core
// -----------------------------------------
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

func (r *repository) UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error {
	return r.db.WithContext(ctx).
		Model(&SevaBooking{}).
		Where("id = ?", bookingID).
		Update("status", newStatus).Error
}

// -----------------------------------------
// Booking Limit Checker
// -----------------------------------------
func (r *repository) CountBookingsForSlot(ctx context.Context, sevaID uint, date time.Time, slot string) (int64, error) {
	var count int64
	// This can be implemented based on your specific slot booking requirements
	// For now, returning 0 as placeholder
	return count, nil
}

// -----------------------------------------
// Detailed Booking Listing
// -----------------------------------------
// DetailedBooking struct is now defined in model.go

func (r *repository) ListBookingsWithDetails(ctx context.Context, entityID uint) ([]DetailedBooking, error) {
	var results []DetailedBooking
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			b.*, 
			s.name AS seva_name, 
			s.seva_type, 
			u.full_name AS devotee_name, 
			u.phone AS devotee_phone
		FROM seva_bookings b
		JOIN sevas s ON s.id = b.seva_id
		JOIN users u ON u.id = b.user_id
		WHERE b.entity_id = ?
		ORDER BY b.booking_time DESC
	`, entityID).Scan(&results).Error

	return results, err
}

// View Booking by ID (for view modal)
func (r *repository) GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error) {
	var booking SevaBooking
	err := r.db.WithContext(ctx).
		Where("id = ?", bookingID).
		First(&booking).Error
	return &booking, err
}

// UPDATED: Search + Filter + Paginate Seva Bookings with UserID support
func (r *repository) SearchBookingsWithFilters(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) {
	var results []DetailedBooking
	var total int64

	query := r.db.WithContext(ctx).
		Table("seva_bookings AS b").
		Select("b.*, s.name AS seva_name, s.seva_type, u.full_name AS devotee_name, u.phone AS devotee_phone").
		Joins("JOIN sevas s ON s.id = b.seva_id").
		Joins("JOIN users u ON u.id = b.user_id").
		Where("b.entity_id = ?", filter.EntityID)

	// Apply filters
	if filter.UserID != 0 {  // NEW: Filter by user ID for GetMyBookings
		query = query.Where("b.user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		query = query.Where("b.status = ?", filter.Status)
	}
	if filter.SevaType != "" {
		query = query.Where("s.seva_type = ?", filter.SevaType)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("s.name LIKE ? OR u.full_name LIKE ?", searchTerm, searchTerm)
	}
	if filter.StartDate != "" && filter.EndDate != "" {
		query = query.Where("b.booking_time BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	// Count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort
	sortBy := "b.booking_time"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}

	err := query.Scan(&results).Error
	return results, total, err
}

// Get Counts by Status
func (r *repository) CountBookingsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error) {
	var counts BookingStatusCounts
	
	// Initialize counts to zero
	counts.Total = 0
	counts.Approved = 0
	counts.Pending = 0
	counts.Rejected = 0

	// Use raw SQL for better performance
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			status,
			COUNT(*) as count
		FROM seva_bookings 
		WHERE entity_id = ? 
		GROUP BY status
	`, entityID).Rows()
	
	if err != nil {
		return counts, err
	}
	defer rows.Close()

	// Map status counts
	statusCounts := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		statusCounts[status] = count
		counts.Total += count
	}

	// Assign specific status counts
	counts.Approved = statusCounts["approved"]
	counts.Pending = statusCounts["pending"]
	counts.Rejected = statusCounts["rejected"]

	return counts, nil
}