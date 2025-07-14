package event

import (
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

	// Count RSVPs for this event
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
		Where("entity_id = ? AND event_date >= CURRENT_DATE", entityID).
		Order("event_date ASC").
		Limit(5).
		Find(&events).Error
	return events, err
}


// ===========================
// 📄 List Events With Pagination + Filters
func (r *Repository) ListEventsByEntity(entityID uint, limit, offset int, category, search, status string) ([]Event, error) {
	var events []Event

	query := r.DB.Where("entity_id = ?", entityID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		ilike := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", ilike, ilike)
	}

	if status != "" {
		if status == "active" {
			query = query.Where("is_active = TRUE")
		} else if status == "inactive" {
			query = query.Where("is_active = FALSE")
		}
	}

	err := query.
		Order("event_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	// ✅ Load RSVP count for each event
	for i := range events {
		var count int64
		r.DB.Table("rsvps").
			Where("event_id = ?", events[i].ID).
			Count(&count)
		events[i].RSVPCount = int(count)
	}

	return events, nil
}




// ===========================
// 🛠 Update Event (Safe Fields)
func (r *Repository) UpdateEvent(id uint, entityID uint, update *UpdateEventRequest) error {
	return r.DB.Model(&Event{}).
		Where("id = ? AND entity_id = ?", id, entityID).
		Updates(map[string]interface{}{
			"title":                 update.Title,
			"description":           update.Description,
			"event_date":            update.EventDate,
			"event_time":            update.EventTime,
			"end_date":              update.EndDate,
			"end_time":              update.EndTime,
			"location":              update.Location,
			"image_url":             update.ImageURL,
			"max_attendees":         update.MaxAttendees,
			"registration_required": update.RegistrationRequired,
			"is_active":             update.IsActive,
		}).Error
}

// ===========================
// ❌ Delete Event (Safe)
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

