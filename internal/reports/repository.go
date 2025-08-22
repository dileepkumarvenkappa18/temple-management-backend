package reports

import (
	"time"

	"gorm.io/gorm"
)

// ReportRepository defines the database operations required by the reports service.
type ReportRepository interface {
	// GetEntitiesByTempleAdmin returns entity IDs created by the given templeadmin user
	GetEntitiesByTenant(userID uint) ([]uint, error)
    GetEvents(entityIDs []uint, start, end time.Time) ([]EventReportRow, error)
    GetSevas(entityIDs []uint, start, end time.Time) ([]SevaReportRow, error)
    GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error)
    GetTemplesRegistered(entityIDs []uint, start, end time.Time, status string) ([]TempleRegisteredReportRow, error)
    GetDevoteeBirthdays(entityIDs []uint, start, end time.Time) ([]DevoteeBirthdayReportRow, error)
    GetDonations(entityIDs []uint, start, end time.Time) ([]DonationReportRow, error) // New method for donations
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ReportRepository {
	return &repository{db: db}
}

func (r *repository) GetEntitiesByTenant(userID uint) ([]uint, error) {
	var ids []uint
	// Table "entities" has created_by which stores templeadmin user ID
	err := r.db.Table("entities").Select("id").Where("created_by = ?", userID).Scan(&ids).Error
	return ids, err
}

func (r *repository) GetEvents(entityIDs []uint, start, end time.Time) ([]EventReportRow, error) {
	var out []EventReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}
	
	// Fixed: Select actual event_type column instead of empty string
	err := r.db.Table("events").
		Select("title, description, event_type, event_date, event_time, location, created_by, created_at, updated_at, is_active").
		Where("entity_id IN ?", entityIDs).
		Where("event_date BETWEEN ? AND ?", start, end).
		Order("event_date DESC").
		Scan(&out).Error
	return out, err
}

func (r *repository) GetSevas(entityIDs []uint, start, end time.Time) ([]SevaReportRow, error) {
	var out []SevaReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}
	
	// Fixed: Use created_at as date since the date column appears to be empty
	// Also handle the empty date column by using created_at instead
	err := r.db.Table("sevas").
		Select("name, seva_type, description, price, created_at as date, start_time, end_time, duration, max_bookings_per_day, status, is_active, created_at, updated_at").
		Where("entity_id IN ?", entityIDs).
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Scan(&out).Error
	return out, err
}

func (r *repository) GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error) {
	var out []SevaBookingReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}

	// Fixed: Get seva_type from sevas table
	err := r.db.Table("seva_bookings sb").
		Select("s.name as seva_name, s.seva_type, u.full_name as devotee_name, u.phone as devotee_phone, sb.booking_time as booking_time, sb.status, sb.created_at, sb.updated_at").
		Joins("LEFT JOIN sevas s ON sb.seva_id = s.id").
		Joins("LEFT JOIN users u ON sb.user_id = u.id").
		Where("sb.entity_id IN ?", entityIDs).
		Where("sb.created_at BETWEEN ? AND ?", start, end).
		Order("sb.created_at DESC").
		Scan(&out).Error
	return out, err
}

// New method to get donations for reporting
// Add this method to your repository struct
func (r *repository) GetDonations(entityIDs []uint, start, end time.Time) ([]DonationReportRow, error) {
    var out []DonationReportRow
    if len(entityIDs) == 0 {
        return out, nil
    }

    err := r.db.Table("donations d").
        Select(`
            d.id, 
            COALESCE(NULLIF(u.full_name, ''), u.email, 'Anonymous') as donor_name, 
            COALESCE(u.email, '') as donor_email,
            d.amount, 
            d.donation_type, 
            d.method as payment_method, 
            d.status, 
            COALESCE(d.donated_at, d.created_at) as donation_date,
            d.order_id,
            d.payment_id,
            d.created_at, 
            d.updated_at
        `).
        Joins("LEFT JOIN users u ON d.user_id = u.id").
        Where("d.entity_id IN ?", entityIDs).
        Where("d.created_at BETWEEN ? AND ?", start, end).
        Order("d.created_at DESC").
        Scan(&out).Error
    return out, err
}

func (r *repository) GetTemplesRegistered(entityIDs []uint, start, end time.Time, status string) ([]TempleRegisteredReportRow, error) {
    var rows []TempleRegisteredReportRow
    if len(entityIDs) == 0 {
        return rows, nil
    }
    query := r.db.Table("entities").Select("id, name, created_at, status").
        Where("id IN ?", entityIDs).
        Where("created_at BETWEEN ? AND ?", start, end)
    if status != "" {
        query = query.Where("status = ?", status)
    }
    err := query.Order("created_at DESC").Scan(&rows).Error
    return rows, err
}

func (r *repository) GetDevoteeBirthdays(entityIDs []uint, start, end time.Time) ([]DevoteeBirthdayReportRow, error) {
	var rows []DevoteeBirthdayReportRow
	if len(entityIDs) == 0 {
		return rows, nil
	}

	// Build the base query with all necessary joins
	query := r.db.Table("users u").
		Select(`
			u.full_name,
			dp.dob as date_of_birth,
			dp.gender,
			u.phone,
			u.email,
			e.name as temple_name,
			uem.joined_at as member_since
		`).
		Joins("INNER JOIN user_entity_memberships uem ON u.id = uem.user_id").
		Joins("INNER JOIN entities e ON uem.entity_id = e.id").
		Joins("INNER JOIN devotee_profiles dp ON u.id = dp.user_id").
		Where("u.role_id = ?", 3). // devotee role
		Where("uem.status = ?", "active").
		Where("uem.entity_id IN ?", entityIDs)

	// Filter by birthday date range
	// For birthdays, we need to check if the birthday (month-day) falls within the date range
	// We'll extract month and day from both the DOB and the range dates
	startMonth := int(start.Month())
	startDay := start.Day()
	endMonth := int(end.Month())
	endDay := end.Day()

	// Handle different scenarios for date range filtering
	if startMonth == endMonth {
		// Same month - simple day range
		query = query.Where(
			"EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) BETWEEN ? AND ?",
			startMonth, startDay, endDay,
		)
	} else if startMonth < endMonth {
		// Range within same year (e.g., March to May)
		query = query.Where(`
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) >= ?) OR
			(EXTRACT(MONTH FROM dp.dob) > ? AND EXTRACT(MONTH FROM dp.dob) < ?) OR
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) <= ?)
		`, startMonth, startDay, startMonth, endMonth, endMonth, endDay)
	} else {
		// Range crosses year boundary (e.g., December to February)
		query = query.Where(`
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) >= ?) OR
			(EXTRACT(MONTH FROM dp.dob) > ?) OR
			(EXTRACT(MONTH FROM dp.dob) < ?) OR
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) <= ?)
		`, startMonth, startDay, startMonth, endMonth, endMonth, endDay)
	}

	// Order by month and day for better readability
	err := query.Order("EXTRACT(MONTH FROM dp.dob), EXTRACT(DAY FROM dp.dob)").
		Scan(&rows).Error

	return rows, err
}