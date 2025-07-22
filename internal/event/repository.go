package event

import (
	"time"

	"gorm.io/gorm"
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
// 🔍 Get Event By ID
func (r *Repository) GetEventByID(id uint) (*Event, error) {
	var e Event
	err := r.DB.First(&e, id).Error
	if err != nil {
		return nil, err
	}

	var count int64
	err = r.DB.Table("rsvps").Where("event_id = ?", id).Count(&count).Error
	if err != nil {
		return nil, err
	}

	e.RSVPCount = int(count)
	return &e, nil
}

// ===========================
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

	for i := range events {
		var count int64
		r.DB.Table("rsvps").Where("event_id = ?", events[i].ID).Count(&count)
		events[i].RSVPCount = int(count)
	}

	return events, nil
}

// ===========================
// 🛠 Update Event
func (r *Repository) UpdateEvent(id uint, entityID uint, update *UpdateEventRequest) error {
	return r.DB.Model(&Event{}).
		Where("id = ? AND entity_id = ?", id, entityID).
		Updates(map[string]interface{}{
			"title":       update.Title,
			"description": update.Description,
			"event_date":  update.EventDate,
			"event_time":  update.EventTime,
			"location":    update.Location,
			"event_type":  update.EventType,
		}).Error
}

// ===========================
// ❌ Delete Event
func (r *Repository) DeleteEvent(id uint, entityID uint) error {
	return r.DB.
		Where("id = ? AND entity_id = ?", id, entityID).
		Delete(&Event{}).Error
}

// ===========================
// 🔢 Count RSVPs for an Event
func (r *Repository) CountRSVPs(eventID uint) (int, error) {
	var count int64
	err := r.DB.Table("rsvps").Where("event_id = ?", eventID).Count(&count).Error
	return int(count), err
}

// ===========================
// 📊 Event Dashboard Stats
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

	// Total Events
	r.DB.Model(&Event{}).
		Where("entity_id = ?", entityID).
		Count(&total)

	// This Month's Events
	r.DB.Model(&Event{}).
		Where("entity_id = ? AND event_date >= ?", entityID, startOfMonth).
		Count(&thisMonth)

	// Upcoming Events
	r.DB.Model(&Event{}).
		Where("entity_id = ? AND event_date >= CURRENT_DATE", entityID).
		Count(&upcoming)

	// Total RSVPs
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
