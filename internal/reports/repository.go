package reports

import (
	"time"

	"gorm.io/gorm"
)

// ReportRepository defines the database operations required by the reports service.
type ReportRepository interface {
	// GetEntitiesByTenant returns entity IDs created by the given tenant (temple admin user)
	GetEntitiesByTenant(userID uint) ([]uint, error)

	GetEvents(entityIDs []uint, start, end time.Time) ([]EventReportRow, error)
	GetSevas(entityIDs []uint, start, end time.Time) ([]SevaReportRow, error)
	GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error)
	GetTemplesRegistered(entityIDs []uint, start, end time.Time, status string) ([]TempleRegisteredReportRow, error)
	GetDevoteeBirthdays(entityIDs []uint, start, end time.Time) ([]DevoteeBirthdayReportRow, error)
	GetDonations(entityIDs []uint, start, end time.Time) ([]DonationReportRow, error)
	GetDevoteeList(entityIDs []uint, start, end time.Time, status string) ([]DevoteeListReportRow, error)
	GetDevoteeProfiles(entityIDs []uint, start, end time.Time, status string) ([]DevoteeProfileReportRow, error)
	GetAuditLogs(entityIDs []uint, start, end time.Time, actionTypes []string, status string) ([]AuditLogReportRow, error)
<<<<<<< HEAD
	GetApprovalStatus(entityIDs []uint, start, end time.Time, role, status string) ([]ApprovalStatusReportRow, error)
	GetUserDetails(entityIDs []uint, start, end time.Time, role, status string) ([]UserDetailsReportRow, error)
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
		Where("created_by = ?", userID).
		Scan(&ids).Error
	return ids, err
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
	return out, err
}

func (r *repository) GetSevaBookings(entityIDs []uint, start, end time.Time) ([]SevaBookingReportRow, error) {
	var out []SevaBookingReportRow
	if len(entityIDs) == 0 {
		return out, nil
	}

	err := r.db.Table("seva_bookings sb").
		Select("s.name as seva_name, s.seva_type, u.full_name as devotee_name, u.phone as devotee_phone, sb.booking_time, sb.status, sb.created_at, sb.updated_at").
		Joins("LEFT JOIN sevas s ON sb.seva_id = s.id").
		Joins("LEFT JOIN users u ON sb.user_id = u.id").
		Where("sb.entity_id IN ?", entityIDs).
		Where("sb.created_at BETWEEN ? AND ?", start, end).
		Order("sb.created_at DESC").
		Scan(&out).Error
	return out, err
}

// GetDonations fetches donation records for reporting
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

	query := r.db.Table("entities").
		Select("id, name, created_at, status").
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
		Where("u.role_id = ?", 3).
		Where("uem.status = ?", "active").
		Where("uem.entity_id IN ?", entityIDs)

	// Birthday range filtering logic
	startMonth := int(start.Month())
	startDay := start.Day()
	endMonth := int(end.Month())
	endDay := end.Day()

	if startMonth == endMonth {
		query = query.Where(
			"EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) BETWEEN ? AND ?",
			startMonth, startDay, endDay,
		)
	} else if startMonth < endMonth {
		query = query.Where(`
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) >= ?) OR
			(EXTRACT(MONTH FROM dp.dob) > ? AND EXTRACT(MONTH FROM dp.dob) < ?) OR
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) <= ?)
		`, startMonth, startDay, startMonth, endMonth, endMonth, endDay)
	} else {
		query = query.Where(`
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) >= ?) OR
			(EXTRACT(MONTH FROM dp.dob) > ?) OR
			(EXTRACT(MONTH FROM dp.dob) < ?) OR
			(EXTRACT(MONTH FROM dp.dob) = ? AND EXTRACT(DAY FROM dp.dob) <= ?)
		`, startMonth, startDay, startMonth, endMonth, endMonth, endDay)
	}

	err := query.Order("EXTRACT(MONTH FROM dp.dob), EXTRACT(DAY FROM dp.dob)").
		Scan(&rows).Error

	return rows, err
}

func (r *repository) GetDevoteeList(entityIDs []uint, start, end time.Time, status string) ([]DevoteeListReportRow, error) {
	var rows []DevoteeListReportRow
	if len(entityIDs) == 0 {
		return rows, nil
	}

	query := r.db.Table("users u").
		Select(`
            u.id as user_id,
            u.full_name as devotee_name,
            uem.joined_at,
            uem.status as devotee_status,
            u.created_at
        `).
		Joins("INNER JOIN user_entity_memberships uem ON u.id = uem.user_id").
		Where("uem.entity_id IN ?", entityIDs)

	if status != "" {
		query = query.Where("uem.status = ?", status)
	}

	query = query.Where("uem.joined_at BETWEEN ? AND ?", start, end).
		Order("uem.joined_at DESC")

	err := query.Scan(&rows).Error
	return rows, err
}

func (r *repository) GetDevoteeProfiles(entityIDs []uint, start, end time.Time, status string) ([]DevoteeProfileReportRow, error) {
	var rows []DevoteeProfileReportRow
	if len(entityIDs) == 0 {
		return rows, nil
	}

	query := r.db.Table("users u").
		Select(`
            u.id as user_id,
            u.full_name,
            dp.dob,
            dp.gender,
            CONCAT(
                COALESCE(dp.street_address, ''), ' ',
                COALESCE(dp.city, ''), ' ',
                COALESCE(dp.state, ''), ' ',
                COALESCE(dp.country, ''), ' ',
                COALESCE(dp.pincode, '')
            ) as full_address,
            COALESCE(dp.gotra, '') as gotra,
            COALESCE(dp.nakshatra, '') as nakshatra,
            COALESCE(dp.rashi, '') as rashi,
            COALESCE(dp.lagna, '') as lagna
        `).
		Joins("INNER JOIN user_entity_memberships uem ON u.id = uem.user_id").
		Joins("INNER JOIN devotee_profiles dp ON u.id = dp.user_id").
		Where("u.role_id = ?", 3).
		Where("uem.entity_id IN ?", entityIDs)

	if status != "" {
		query = query.Where("uem.status = ?", status)
	}

	query = query.Where("uem.joined_at BETWEEN ? AND ?", start, end).
		Order("u.full_name ASC")

	err := query.Scan(&rows).Error
	return rows, err
}

// GetAuditLogs fetches audit log entries for given entity IDs, date range, action types, and status
func (r *repository) GetAuditLogs(entityIDs []uint, start, end time.Time, actionTypes []string, status string) ([]AuditLogReportRow, error) {
	var rows []AuditLogReportRow
	if len(entityIDs) == 0 {
		return rows, nil
	}

	query := r.db.Table("audit_logs al").
		Select(`
            al.id,
            al.entity_id,
            e.name AS entity_name,
            al.user_id,
            u.full_name AS user_name,
            COALESCE(ur.role_name, '') AS user_role,
            al.action,
            al.status,
            al.ip_address,
            al.created_at AS timestamp,
<<<<<<< HEAD
            al.created_at,
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
            COALESCE(al.details::text, '') AS details
        `).
		Joins("LEFT JOIN users u ON al.user_id = u.id").
		Joins("LEFT JOIN entities e ON al.entity_id = e.id").
		Joins("LEFT JOIN user_roles ur ON u.role_id = ur.id").
		Where("al.entity_id IN ?", entityIDs).
		Where("al.created_at BETWEEN ? AND ?", start, end)

	if len(actionTypes) > 0 {
		query = query.Where("al.action IN ?", actionTypes)
	}

	if status != "" {
		query = query.Where("al.status = ?", status)
	}

	err := query.Order("al.created_at DESC").Scan(&rows).Error
	return rows, err
}
<<<<<<< HEAD
// GetApprovalStatus fetches approval status records for reporting
func (r *repository) GetApprovalStatus(entityIDs []uint, start, end time.Time, role, status string) ([]ApprovalStatusReportRow, error) {
	var rows []ApprovalStatusReportRow

	query := r.db.Table("users u").
		Select(`
            u.full_name as name,
            COALESCE(CAST(uem.entity_id AS CHAR), 'N/A') as id,  -- tenant_id replaced as id
            ur.role_name as role,
            uem.status,          -- real status from DB
            u.created_at,
            u.email
        `).
		Joins("LEFT JOIN user_entity_memberships uem ON u.id = uem.user_id").
		Joins("LEFT JOIN user_roles ur ON u.role_id = ur.id")

	// Filter by entity only if provided
	if len(entityIDs) > 0 {
		query = query.Where("uem.entity_id IN ?", entityIDs)
	}

	if !start.IsZero() && !end.IsZero() {
		query = query.Where("u.created_at BETWEEN ? AND ?", start, end)
	}

	if role != "" {
		query = query.Where("ur.role_name = ?", role)
	}

	if status != "" {
		query = query.Where("uem.status = ?", status)
	}

	err := query.Order("u.created_at DESC").Scan(&rows).Error

	// Replace null/empty status with "N/A" for users without membership
	for i := range rows {
		if rows[i].Status == "" {
			rows[i].Status = "N/A"
		}
	}

	return rows, err
}


// GetUserDetails fetches user detail records for reporting
func (r *repository) GetUserDetails(entityIDs []uint, start, end time.Time, role, status string) ([]UserDetailsReportRow, error) {
	var rows []UserDetailsReportRow

	query := r.db.Table("users u").
		Select(`
            u.id,
            u.full_name as name,
            COALESCE(e.name,'N/A') as entity_name,
            COALESCE(CAST(uem.entity_id AS CHAR), 'N/A') as tenant_id,
            u.email,
            ur.role_name as role,
            COALESCE(uem.status, 'Active') as status, -- show all statuses: Active/Inactive/Locked
            u.created_at
        `).
		Joins("LEFT JOIN user_entity_memberships uem ON u.id = uem.user_id").
		Joins("LEFT JOIN entities e ON uem.entity_id = e.id").
		Joins("LEFT JOIN user_roles ur ON u.role_id = ur.id")

	if len(entityIDs) > 0 {
		query = query.Where("uem.entity_id IN ?", entityIDs)
	}

	if !start.IsZero() && !end.IsZero() {
		query = query.Where("u.created_at BETWEEN ? AND ?", start, end)
	}

	if role != "" {
		query = query.Where("ur.role_name = ?", role)
	}

	if status != "" {
		query = query.Where("uem.status = ?", status)
	}

	err := query.Order("u.created_at DESC").Scan(&rows).Error
	return rows, err
}
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
