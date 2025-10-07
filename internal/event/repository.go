package event

import (
	"time"
<<<<<<< HEAD
	"gorm.io/gorm"
	"fmt"
=======

	"gorm.io/gorm"
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// ===========================
// 🎯 Create Event
func (r *Repository) CreateEvent(e *Event) error {
	return r.DB.Create(e).Error
}

// ===========================
<<<<<<< HEAD
// 🔍 Get Event By ID with entity validation and proper RSVP counting
=======
// 🔍 Get Event By ID
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *Repository) GetEventByID(id uint) (*Event, error) {
	var e Event
	err := r.DB.First(&e, id).Error
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// Get RSVP count for this specific event with entity validation
	var count int64
	err = r.DB.Table("rsvps").
		Joins("JOIN events ON events.id = rsvps.event_id").
		Where("rsvps.event_id = ? AND events.entity_id = ?", id, e.EntityID).
		Count(&count).Error
=======
	var count int64
	err = r.DB.Table("rsvps").Where("event_id = ?", id).Count(&count).Error
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	if err != nil {
		return nil, err
	}

	e.RSVPCount = int(count)
	return &e, nil
}

// ===========================
<<<<<<< HEAD
// 📆 Get Upcoming Events - FIXED to use proper entity filtering
func (r *Repository) GetUpcomingEvents(entityID uint) ([]Event, error) {
	var events []Event
	
	// Ensure we only get events from the specified entity
	err := r.DB.
		Where("entity_id = ? AND event_date >= CURRENT_DATE - INTERVAL '7 day' AND is_active = TRUE", entityID).
		Order("event_date ASC").
		Find(&events).Error
	
	if err != nil {
		return nil, err
	}

	// Add RSVP counts for each event - FIXED to ensure entity filtering
	for i := range events {
		var count int64
		r.DB.Table("rsvps").
			Joins("JOIN events ON events.id = rsvps.event_id").
			Where("rsvps.event_id = ? AND events.entity_id = ?", events[i].ID, entityID).
			Count(&count)
		events[i].RSVPCount = int(count)
	}

	return events, nil
}

// ===========================
// 📄 List Events With Pagination & Search - FIXED entity filtering
=======
// 📆 Get Upcoming Events
func (r *Repository) GetUpcomingEvents(entityID uint) ([]Event, error) {
	var events []Event
	err := r.DB.
		Where("entity_id = ? AND event_date >= CURRENT_DATE AND is_active = TRUE", entityID).
		Order("event_date ASC").
		Limit(5).
		Find(&events).Error
	return events, err
}

// ===========================
// 📄 List Events With Pagination & Search
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *Repository) ListEventsByEntity(entityID uint, limit, offset int, search string) ([]Event, error) {
	var events []Event

	query := r.DB.Where("entity_id = ?", entityID)

	if search != "" {
		ilike := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", ilike, ilike)
	}

	err := query.
		Order("event_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// FIXED: Ensure RSVP counts are only from events belonging to the specified entity
	for i := range events {
		var count int64
		r.DB.Table("rsvps").
			Joins("JOIN events ON events.id = rsvps.event_id").
			Where("rsvps.event_id = ? AND events.entity_id = ?", events[i].ID, entityID).
			Count(&count)
=======
	for i := range events {
		var count int64
		r.DB.Table("rsvps").Where("event_id = ?", events[i].ID).Count(&count)
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		events[i].RSVPCount = int(count)
	}

	return events, nil
}

// ===========================
// 🛠 Update Event
<<<<<<< HEAD
=======
// 🛠 Update Event
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func (r *Repository) UpdateEvent(e *Event) error {
	return r.DB.Save(e).Error
}

<<<<<<< HEAD
=======

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
// ===========================
// ❌ Delete Event
func (r *Repository) DeleteEvent(id uint, entityID uint) error {
	return r.DB.
		Where("id = ? AND entity_id = ?", id, entityID).
		Delete(&Event{}).Error
}

// ===========================
<<<<<<< HEAD
// 🔢 Count RSVPs for an Event - FIXED to validate entity ownership
func (r *Repository) CountRSVPs(eventID uint) (int, error) {
	// First get the event to determine its entity_id
	var event Event
	err := r.DB.First(&event, eventID).Error
	if err != nil {
		return 0, err
	}

	// Now count RSVPs with proper entity validation
	var count int64
	err = r.DB.Table("rsvps").
		Joins("JOIN events ON events.id = rsvps.event_id").
		Where("rsvps.event_id = ? AND events.entity_id = ?", eventID, event.EntityID).
		Count(&count).Error
=======
// 🔢 Count RSVPs for an Event
func (r *Repository) CountRSVPs(eventID uint) (int, error) {
	var count int64
	err := r.DB.Table("rsvps").Where("event_id = ?", eventID).Count(&count).Error
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	return int(count), err
}

// ===========================
<<<<<<< HEAD
// 🔢 Count RSVPs for an Event by Entity - Enhanced method for entity-specific counting
func (r *Repository) CountRSVPsByEntity(eventID uint, entityID uint) (int, error) {
	var count int64
	err := r.DB.Table("rsvps").
		Joins("JOIN events ON events.id = rsvps.event_id").
		Where("rsvps.event_id = ? AND events.entity_id = ?", eventID, entityID).
		Count(&count).Error
	return int(count), err
}

// ===========================
// 📊 Event Dashboard Stats - FIXED to use proper entity filtering
=======
// 📊 Event Dashboard Stats
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
type EventStatsResponse struct {
	TotalEvents     int `json:"total_events"`
	ThisMonthEvents int `json:"this_month_events"`
	UpcomingEvents  int `json:"upcoming_events"`
	TotalRSVPs      int `json:"total_rsvps"`
}

func (r *Repository) GetEventStats(entityID uint) (*EventStatsResponse, error) {
	var stats EventStatsResponse
	var total, thisMonth, upcoming, totalRSVPs int64

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

<<<<<<< HEAD
	// Total Events for this entity only
=======
	// Total Events
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	r.DB.Model(&Event{}).
		Where("entity_id = ?", entityID).
		Count(&total)

<<<<<<< HEAD
	// This Month's Events for this entity only
=======
	// This Month's Events
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	r.DB.Model(&Event{}).
		Where("entity_id = ? AND event_date >= ?", entityID, startOfMonth).
		Count(&thisMonth)

<<<<<<< HEAD
	// Upcoming Events for this entity only
=======
	// Upcoming Events
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	r.DB.Model(&Event{}).
		Where("entity_id = ? AND event_date >= CURRENT_DATE", entityID).
		Count(&upcoming)

<<<<<<< HEAD
	// FIXED: Total RSVPs for events belonging to this entity only
=======
	// Total RSVPs
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	r.DB.Table("rsvps").
		Joins("JOIN events ON events.id = rsvps.event_id").
		Where("events.entity_id = ?", entityID).
		Count(&totalRSVPs)

	stats.TotalEvents = int(total)
	stats.ThisMonthEvents = int(thisMonth)
	stats.UpcomingEvents = int(upcoming)
	stats.TotalRSVPs = int(totalRSVPs)

	return &stats, nil
}
<<<<<<< HEAD

// ===========================
// 🔢 NEW: Get Total RSVP Count by Entity - Additional helper method
func (r *Repository) GetTotalRSVPsByEntity(entityID uint) (int, error) {
	var count int64
	err := r.DB.Table("rsvps").
		Joins("JOIN events ON events.id = rsvps.event_id").
		Where("events.entity_id = ?", entityID).
		Count(&count).Error
	return int (count), err
}

// ===========================
// 🔢 NEW: Get Event Count by Entity - Additional helper method
func (r *Repository) GetEventCountByEntity(entityID uint) (int, error) {
	var count int64
	err := r.DB.Model(&Event{}).
		Where("entity_id = ?", entityID).
		Count(&count).Error
    fmt.Println("count ()=",count)
	return int(count), err
}
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
