package reports

import "time"

// ===== Report Types =====
const (
	ReportTypeEvents   = "events"
	ReportTypeSevas    = "sevas"
	ReportTypeBookings = "bookings"
	ReportTypeTempleRegistered = "temple_registered"
	 // ... existing report types
    ReportTypeTempleRegisteredPDF   = "temples_registered_pdf"
    ReportTypeTempleRegisteredExcel = "temples_registered_excel"
)

// ===== Date Range Presets =====
const (
	DateRangeDaily   = "daily"
	DateRangeWeekly  = "weekly"
	DateRangeMonthly = "monthly"
	DateRangeYearly  = "yearly"
	DateRangeCustom  = "custom"
)

// ===== Export Formats =====
const (
	FormatExcel = "excel"
	FormatCSV   = "csv"
	FormatPDF   = "pdf"
)

// ===== Request Struct =====
type ActivitiesReportRequest struct {
	EntityID   string    `form:"entity_id" json:"entity_id"`     // "all" or specific UUID
	Type       string    `form:"type" json:"type"`               // events, sevas, bookings
	DateRange  string    `form:"date_range" json:"date_range"`   // daily, weekly, monthly, yearly, custom
	StartDate  time.Time `form:"start_date" json:"start_date"`   // required if date_range=custom
	EndDate    time.Time `form:"end_date" json:"end_date"`       // required if date_range=custom
	Format     string    `form:"format" json:"format"`           // excel, csv, pdf
	TenantID   string    `json:"-"`                              // extracted from auth context
	EntityIDs  []string  `json:"-"`                              // resolved entity IDs from DB
}

// ===== Events Report Row =====
type EventReportRow struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	EventDate   time.Time `json:"event_date"`
	EventTime   string    `json:"event_time"`
	Location    string    `json:"location"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
}

// ===== Sevas Report Row =====
type SevaReportRow struct {
	Name               string    `json:"name"`
	SevaType           string    `json:"seva_type"`
	Description        string    `json:"description"`
	Price              float64   `json:"price"`
	Date               time.Time `json:"date"`
	StartTime          string    `json:"start_time"`
	EndTime            string    `json:"end_time"`
	Duration           string    `json:"duration"`
	MaxBookingsPerDay  int       `json:"max_bookings_per_day"`
	Status             string    `json:"status"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ===== Seva Bookings Report Row =====
type SevaBookingReportRow struct {
	SevaName     string    `json:"seva_name"`
	SevaType     string    `json:"seva_type"`
	DevoteeName  string    `json:"devotee_name"`
	DevoteePhone string    `json:"devotee_phone"`
	BookingTime  time.Time `json:"booking_time"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ===== Generic Report Response =====
type ReportData struct {
	Events   []EventReportRow        `json:"events,omitempty"`
	Sevas    []SevaReportRow         `json:"sevas,omitempty"`
	Bookings []SevaBookingReportRow  `json:"bookings,omitempty"`
	TemplesRegistered []TempleRegisteredReportRow
}

type TempleRegisteredReportRequest struct {
    DateRange  string    `form:"date_range" json:"date_range"`
    StartDate  time.Time `form:"start_date" json:"start_date"`
    EndDate    time.Time `form:"end_date" json:"end_date"`
    Status     string    `form:"status" json:"status"`
    Format     string    `form:"format" json:"format"`
    EntityID   string    `form:"entity_id" json:"entity_id"` // use this for filtering temples
}

type TempleRegisteredReportRow struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    Status    string    `json:"status"`
}

