package seva

import (
    "time"
)

// ======================
// 🔹 Seva Core Model
// ======================
type Seva struct {
    ID                uint      `gorm:"primaryKey" json:"id"`
    EntityID          uint      `gorm:"not null" json:"entity_id"`
    Name              string    `gorm:"type:varchar(255);not null" json:"name"`
    SevaType          string    `gorm:"type:varchar(50);not null" json:"seva_type"` // e.g., Archana, Abhishekam
    Description       string    `gorm:"type:text" json:"description"`
    Price             float64   `gorm:"type:decimal(10,2);default:0" json:"price"`
    Date              string    `gorm:"type:varchar(20)" json:"date"`        // Format: dd-mm-yyyy
    StartTime         string    `gorm:"type:varchar(10)" json:"start_time"`  // Format: HH:mm
    EndTime           string    `gorm:"type:varchar(10)" json:"end_time"`    // Format: HH:mm
    Duration          int       `json:"duration"` // in minutes
    MaxBookingsPerDay int       `json:"max_bookings_per_day"`
    Status            string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // e.g., approved/pending/rejected
    IsActive          bool      `gorm:"default:true" json:"is_active"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

// ======================
// 🔹 Booking Model
// ======================
type SevaBooking struct {
    ID          uint      `gorm:"primaryKey" json:"id"` // PK
    SevaID      uint      `gorm:"not null" json:"seva_id"` // Which Seva is being booked
    UserID      uint      `gorm:"not null" json:"user_id"` // Who is booking (devotee)
    EntityID    uint      `gorm:"not null" json:"entity_id"` // Temple where the seva is hosted
    BookingTime time.Time `gorm:"not null" json:"booking_time"` // Auto-timestamp
    Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending / approved / rejected
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// DetailedBooking combines booking with seva and user details
type DetailedBooking struct {
    SevaBooking
    SevaName     string `json:"seva_name"`
    SevaType     string `json:"seva_type"`
    DevoteeName  string `json:"devotee_name"`
    DevoteePhone string `json:"devotee_phone"`
}

// ✅ For Filtered Search (Admin Dashboard) - UPDATED with UserID
type BookingFilter struct {
    EntityID   uint   `json:"entity_id"`
    UserID     uint   `json:"user_id,omitempty"`    // NEW: Filter by specific user (for GetMyBookings)
    Status     string `json:"status"`
    SevaType   string `json:"seva_type"`
    Search     string `json:"search"`
    StartDate  string `json:"start_date"`
    EndDate    string `json:"end_date"`
    SortBy     string `json:"sort_by"`
    SortOrder  string `json:"sort_order"`
    Limit      int    `json:"limit"`
    Offset     int    `json:"offset"`
}

// ✅ Booking Status Counts (Dashboard Card)
type BookingStatusCounts struct {
    Total    int64 `json:"total"`
    Approved int64 `json:"approved"`
    Pending  int64 `json:"pending"`
    Rejected int64 `json:"rejected"`
}