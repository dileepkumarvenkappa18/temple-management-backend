package reports

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TempleadminInfo represents basic templeadmin information for superadmin selection
type TempleadminInfo struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

// ReportRepository defines the database operations required by the reports service.
type ReportRepository interface {
	// Existing methods
	GetEntitiesByTenant(userID uint) ([]uint, error)
	GetEvents(entityIDs []uint, start, end time.Time) ([]EventReportRow, error)
	GetSevas(entityIDs []uint, start, end time.Time) ([]SevaReportRow, error)
	GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error)
	GetTemplesRegistered(entityIDs []uint, start, end time.Time, status string) ([]TempleRegisteredReportRow, error)
	GetDevoteeBirthdays(entityIDs []uint, start, end time.Time) ([]DevoteeBirthdayReportRow, error)
	
	// New methods for superadmin support
	GetEntitiesByMultipleTenants(userIDs []uint) ([]uint, error)
	ValidateEntityOwnershipByMultipleTenants(entityID uint, userIDs []uint) (bool, error)
	GetAllTempleadmins() ([]TempleadminInfo, error)

	CountEventsByEntities(entityIDs []uint) (int64, error)
    CountSevasByEntities(entityIDs []uint) (int64, error)
    CountBookingsByEntities(entityIDs []uint) (int64, error)
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
	err := r.db.Table("entities").
		Select("id").
		Where("created_by = ? AND status = ?", userID, "approved").
		Scan(&ids).Error
	
	fmt.Printf("DEBUG: GetEntitiesByTenant - UserID: %d, Found entities: %v\n", userID, ids)
	return ids, err
}

// GetEntitiesByMultipleTenants returns entity IDs created by multiple templeadmins
func (r *repository) GetEntitiesByMultipleTenants(userIDs []uint) ([]uint, error) {
	var ids []uint
	if len(userIDs) == 0 {
		return ids, nil
	}
	
	err := r.db.Table("entities").
		Select("id").
		Where("created_by IN ? AND status = ?", userIDs, "approved").
		Scan(&ids).Error
	
	fmt.Printf("DEBUG: GetEntitiesByMultipleTenants - UserIDs: %v, Found entities: %v\n", userIDs, ids)
	return ids, err
}

// ValidateEntityOwnershipByMultipleTenants checks if an entity is owned by any of the given tenants
func (r *repository) ValidateEntityOwnershipByMultipleTenants(entityID uint, userIDs []uint) (bool, error) {
	if len(userIDs) == 0 {
		return false, nil
	}
	
	var count int64
	err := r.db.Table("entities").
		Where("id = ? AND created_by IN ?", entityID, userIDs).
		Count(&count).Error
	
	fmt.Printf("DEBUG: ValidateEntityOwnershipByMultipleTenants - EntityID: %d, UserIDs: %v, Found: %t\n", entityID, userIDs, count > 0)
	return count > 0, err
}

// GetAllTempleadmins returns list of all templeadmins for superadmin selection
func (r *repository) GetAllTempleadmins() ([]TempleadminInfo, error) {
	var templeadmins []TempleadminInfo
	
	err := r.db.Table("users u").
		Select("u.id, u.full_name, u.email, u.phone, u.status").
		Joins("INNER JOIN user_roles ur ON u.role_id = ur.id").
		Where("ur.role_name = ?", "templeadmin").
		Where("u.status = ?", "active").
		Order("u.full_name ASC").
		Scan(&templeadmins).Error
	
	fmt.Printf("DEBUG: GetAllTempleadmins - Found %d templeadmins\n", len(templeadmins))
	return templeadmins, err
}

func (r *repository) GetEvents(entityIDs []uint, start, end time.Time) ([]EventReportRow, error) {
	var out []EventReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}
	
	err := r.db.Table("events").
		Select("title, description, event_type, event_date, event_time, location, created_by, created_at, updated_at, is_active").
		Where("entity_id IN ?", entityIDs).
		Where("event_date BETWEEN ? AND ?", start, end).
		Order("event_date DESC").
		Scan(&out).Error
	
	fmt.Printf("DEBUG: GetEvents - EntityIDs: %v, DateRange: %v to %v, Found: %d events\n", entityIDs, start.Format("2006-01-02"), end.Format("2006-01-02"), len(out))
	return out, err
}

func (r *repository) GetSevas(entityIDs []uint, start, end time.Time) ([]SevaReportRow, error) {
	var out []SevaReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}
	
	err := r.db.Table("sevas").
		Select("name, seva_type, description, price, created_at as date, start_time, end_time, duration, max_bookings_per_day, status, is_active, created_at, updated_at").
		Where("entity_id IN ?", entityIDs).
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Scan(&out).Error
	
	fmt.Printf("DEBUG: GetSevas - EntityIDs: %v, DateRange: %v to %v, Found: %d sevas\n", entityIDs, start.Format("2006-01-02"), end.Format("2006-01-02"), len(out))
	return out, err
}

func (r *repository) GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error) {
	var out []SevaBookingReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}

	err := r.db.Table("seva_bookings sb").
		Select("s.name as seva_name, s.seva_type, u.full_name as devotee_name, u.phone as devotee_phone, sb.booking_time as booking_time, sb.status, sb.created_at, sb.updated_at").
		Joins("LEFT JOIN sevas s ON sb.seva_id = s.id").
		Joins("LEFT JOIN users u ON sb.user_id = u.id").
		Where("sb.entity_id IN ?", entityIDs).
		Where("sb.created_at BETWEEN ? AND ?", start, end).
		Order("sb.created_at DESC").
		Scan(&out).Error
	
	fmt.Printf("DEBUG: GetSevaBookings - EntityIDs: %v, DateRange: %v to %v, Found: %d bookings\n", entityIDs, start.Format("2006-01-02"), end.Format("2006-01-02"), len(out))
	return out, err
}

func (r *repository) GetTemplesRegistered(entityIDs []uint, start, end time.Time, status string) ([]TempleRegisteredReportRow, error) {
	var rows []TempleRegisteredReportRow
	if len(entityIDs) == 0 {
		return rows, nil
	}
	
	query := r.db.Table("entities").
		Select("id, name, created_at, status").
		Where("id IN ?", entityIDs).
		Where("created_at BETWEEN ? AND ?", start, end)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("created_at DESC").Scan(&rows).Error
	
	fmt.Printf("DEBUG: GetTemplesRegistered - EntityIDs: %v, DateRange: %v to %v, Status: %s, Found: %d temples\n", entityIDs, start.Format("2006-01-02"), end.Format("2006-01-02"), status, len(rows))
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

	fmt.Printf("DEBUG: GetDevoteeBirthdays - EntityIDs: %v, DateRange: %v to %v, Found: %d birthdays\n", entityIDs, start.Format("2006-01-02"), end.Format("2006-01-02"), len(rows))
	return rows, err
}

// in reports/repository.go
func (r *repository) CountEventsByEntities(entityIDs []uint) (int64, error) {
	var count int64
	err := r.db.Table("events").Where("entity_id IN ?", entityIDs).Count(&count).Error
	return count, err
}

func (r *repository) CountSevasByEntities(entityIDs []uint) (int64, error) {
	var count int64
	err := r.db.Table("sevas").Where("entity_id IN ?", entityIDs).Count(&count).Error
	return count, err
}

func (r *repository) CountBookingsByEntities(entityIDs []uint) (int64, error) {
	var count int64
	err := r.db.Table("seva_bookings").Where("entity_id IN ?", entityIDs).Count(&count).Error
	return count, err
}
