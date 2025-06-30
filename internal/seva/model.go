package seva

import (
	"time"
)

// Seva represents a ritual/service offered by a temple
type Seva struct {
	ID                   uint      `gorm:"primaryKey" json:"id" example:"1"`
	EntityID             uint      `gorm:"not null" json:"entity_id" example:"1"`                      // Temple ID
	Name                 string    `gorm:"type:varchar(255);not null" json:"name" example:"Abhishekam"` // Seva Name
	Description          string    `gorm:"type:text" json:"description" example:"Milk abhishekam to deity"`
	Price                float64   `gorm:"type:decimal(10,2);default:0" json:"price" example:"500"`
	Duration             int       `json:"duration" example:"30"` // Duration in minutes
	AvailabilitySchedule string    `gorm:"type:text" json:"availability_schedule" example:"{\"Monday\": [\"9:00\", \"10:00\"]}"`
	MaxBookingsPerDay    int       `json:"max_bookings_per_day" example:"10"`
	IsActive             bool      `gorm:"default:true" json:"is_active" example:"true"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// SevaBooking represents a seva booked by a devotee
type SevaBooking struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SevaID          uint      `gorm:"not null" json:"seva_id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	EntityID        uint      `gorm:"not null" json:"entity_id"`
	BookingDate     time.Time `gorm:"type:date;not null" json:"booking_date"`
	BookingTime     time.Time `json:"booking_time"` // ✅ MATCHES `timestamp`
	Status          string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	SpecialRequests string    `gorm:"type:text" json:"special_requests"`
	AmountPaid      float64   `gorm:"type:decimal(10,2)" json:"amount_paid"`
	PaymentStatus   string    `gorm:"type:varchar(20);default:'pending'" json:"payment_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}


