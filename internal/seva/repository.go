<<<<<<< HEAD

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
	// Enhanced seva listing with filters
	GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error)

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	// Booking core
	BookSeva(ctx context.Context, booking *SevaBooking) error
	ListBookingsByUserID(ctx context.Context, userID uint) ([]SevaBooking, error)
	ListBookingsByEntityID(ctx context.Context, entityID uint) ([]SevaBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uint, newStatus string) error

<<<<<<< HEAD
	// Booking limits
	CountBookingsForSlot(ctx context.Context, sevaID uint, date time.Time, slot string) (int64, error)

	// Composite list with Seva + User info
	ListBookingsWithDetails(ctx context.Context, entityID uint) ([]DetailedBooking, error)

	GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)
	SearchBookingsWithFilters(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error)
	CountBookingsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error)

	ListPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)
=======
	// 🔄 Availability schedule
	// CreateAvailability(ctx context.Context, slot *SevaAvailability) error
	// GetAvailabilityBySevaID(ctx context.Context, sevaID uint) ([]SevaAvailability, error)
	// DeleteAvailabilityBySevaID(ctx context.Context, sevaID uint) error

	// 🔄 Booking limits
	CountBookingsForSlot(ctx context.Context, sevaID uint, date time.Time, slot string) (int64, error)

	// 🔄 Composite list with Seva + User info
	ListBookingsWithDetails(ctx context.Context, entityID uint) ([]DetailedBooking, error)

		GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error)                          // 🆕
	SearchBookingsWithFilters(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) // 🆕
	CountBookingsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error)                // 🆕

	ListPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error)


>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
}

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

<<<<<<< HEAD
func (r *repository) ListSevasByEntityID(ctx context.Context, entityID uint) ([]Seva, error) {
	var sevas []Seva
	err := r.db.WithContext(ctx).
		Where("entity_id = ?", entityID).
=======

func (r *repository) ListSevasByEntityID(ctx context.Context, entityID uint) ([]Seva, error) {
	var sevas []Seva
	err := r.db.WithContext(ctx).
		Where("entity_id = ? AND is_active = true", entityID).
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		Find(&sevas).Error
	return sevas, err
}

<<<<<<< HEAD
=======

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *repository) ListPaginatedSevas(ctx context.Context, entityID uint, sevaType string, search string, limit int, offset int) ([]Seva, error) {
	var sevas []Seva

	query := r.db.WithContext(ctx).
		Model(&Seva{}).
<<<<<<< HEAD
		Where("entity_id = ? AND status IN (?)", entityID, []string{"upcoming", "ongoing"}) // Only show bookable sevas to devotees
=======
		Where("entity_id = ? AND is_active = true", entityID)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6

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

<<<<<<< HEAD
// Enhanced seva listing with filters for temple admin
func (r *repository) GetSevasWithFilters(ctx context.Context, entityID uint, sevaType, search, status string, limit, offset int) ([]Seva, int64, error) {
=======


func (r *repository) GetSevasWithFilters(ctx context.Context, entityID uint, search, sevaType string, page, limit int) ([]Seva, int64, error) {
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	var sevas []Seva
	var total int64

	query := r.db.WithContext(ctx).
		Model(&Seva{}).
<<<<<<< HEAD
		Where("entity_id = ?", entityID)

	// Apply filters
=======
		Where("entity_id = ? AND is_active = true", entityID)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	if sevaType != "" {
		query = query.Where("seva_type = ?", sevaType)
	}

<<<<<<< HEAD
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total before pagination
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	// Apply pagination and ordering
	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
=======
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	}

	err := query.Find(&sevas).Error
	return sevas, total, err
}

<<<<<<< HEAD
=======


>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *repository) UpdateSeva(ctx context.Context, seva *Seva) error {
	return r.db.WithContext(ctx).Save(seva).Error
}

<<<<<<< HEAD
// Permanent delete - removes the record completely
func (r *repository) DeleteSeva(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&Seva{}, id).Error
=======
func (r *repository) DeleteSeva(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Seva{}, id).Error
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
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
=======

// -----------------------------------------
// 🔄 Seva Availability
// -----------------------------------------
// func (r *repository) CreateAvailability(ctx context.Context, slot *SevaAvailability) error {
// 	return r.db.WithContext(ctx).Create(slot).Error
// }

// func (r *repository) GetAvailabilityBySevaID(ctx context.Context, sevaID uint) ([]SevaAvailability, error) {
// 	var slots []SevaAvailability
// 	err := r.db.WithContext(ctx).
// 		Where("seva_id = ?", sevaID).
// 		Order("date, time_slot").
// 		Find(&slots).Error
// 	return slots, err
// }

// func (r *repository) DeleteAvailabilityBySevaID(ctx context.Context, sevaID uint) error {
// 	return r.db.WithContext(ctx).
// 		Where("seva_id = ?", sevaID).
// 		Delete(&SevaAvailability{}).Error
// }

// -----------------------------------------
// 🔄 Booking Limit Checker
// -----------------------------------------
func (r *repository) CountBookingsForSlot(ctx context.Context, sevaID uint, date time.Time, slot string) (int64, error) {
	var count int64
	// err := r.db.WithContext(ctx).
	// 	Model(&SevaBooking{}).
	// 	Where("seva_id = ? AND booking_date = ? AND TO_CHAR(booking_time, 'HH24:MI') = ?", sevaID, date.Format("2006-01-02"), slot).
	// 	Count(&count).Error
	return count, nil
}


// -----------------------------------------
// 🔄 Detailed Booking Listing
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
// -----------------------------------------
type DetailedBooking struct {
	SevaBooking
	SevaName     string `json:"seva_name"`
	SevaType     string `json:"seva_type"`
	DevoteeName  string `json:"devotee_name"`
	DevoteePhone string `json:"devotee_phone"`
}

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

<<<<<<< HEAD
// View Booking by ID (for view modal)
=======

// 🆕 View Booking by ID (for view modal)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *repository) GetBookingByID(ctx context.Context, bookingID uint) (*SevaBooking, error) {
	var booking SevaBooking
	err := r.db.WithContext(ctx).
		Where("id = ?", bookingID).
		First(&booking).Error
	return &booking, err
}

<<<<<<< HEAD
// Search + Filter + Paginate Seva Bookings
=======
// 🆕 Search + Filter + Paginate Seva Bookings
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *repository) SearchBookingsWithFilters(ctx context.Context, filter BookingFilter) ([]DetailedBooking, int64, error) {
	var results []DetailedBooking
	var total int64

	query := r.db.WithContext(ctx).
		Table("seva_bookings AS b").
		Select("b.*, s.name AS seva_name, s.seva_type, u.full_name AS devotee_name, u.phone AS devotee_phone").
		Joins("JOIN sevas s ON s.id = b.seva_id").
		Joins("JOIN users u ON u.id = b.user_id").
		Where("b.entity_id = ?", filter.EntityID)

<<<<<<< HEAD
	// Apply filters
=======
	// ✅ Apply filters
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
<<<<<<< HEAD
		query = query.Where("b.booking_time BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	// Count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort
=======
		// ⚠️ Replace booking_date with created_at or booking_time
		query = query.Where("b.booking_time BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	// ✅ Count before pagination
	query.Count(&total)

	// ✅ Sort
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	sortBy := "b.booking_time"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

<<<<<<< HEAD
	// Pagination
=======
	// ✅ Pagination
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}

	err := query.Scan(&results).Error
	return results, total, err
}

<<<<<<< HEAD
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
=======

// 🆕 Get Counts by Status
func (r *repository) CountBookingsByStatus(ctx context.Context, entityID uint) (BookingStatusCounts, error) {
	var counts BookingStatusCounts
	counts.Total = 0

	query := r.db.WithContext(ctx).
		Model(&SevaBooking{}).
		Where("entity_id = ?", entityID)

	// Total
	query.Count(&counts.Total)

	// Approved
	query.Where("status = ?", "approved").Count(&counts.Approved)

	// Pending
	query.Where("status = ?", "pending").Count(&counts.Pending)

	// Rejected
	query.Where("status = ?", "rejected").Count(&counts.Rejected)

	return counts, nil
}

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
