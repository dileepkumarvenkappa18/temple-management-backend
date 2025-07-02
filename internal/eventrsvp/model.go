package eventrsvp

import "time"

type RSVP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   uint      `gorm:"not null;index:idx_event_user,unique" json:"event_id"`
	UserID    uint      `gorm:"not null;index:idx_event_user,unique" json:"user_id"`
	Status    string    `gorm:"type:varchar(20);default:'attending'" json:"status"` // attending, maybe, not_attending
	Notes     string    `gorm:"type:text"`
	RSVPDate  time.Time `gorm:"autoCreateTime" json:"rsvp_date"`
}