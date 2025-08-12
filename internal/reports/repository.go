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