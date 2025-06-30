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

// CreateEvent inserts a new event
func (r *Repository) CreateEvent(e *Event) error {
	return r.DB.Create(e).Error
}

// GetEventByID fetches a specific event by its ID
func (r *Repository) GetEventByID(id uint) (*Event, error) {
	var e Event
	err := r.DB.First(&e, id).Error
	return &e, err
}

// ✅ FIXED: GetUpcomingEvents - fetch events from today onward (date only comparison)
func (r *Repository) GetUpcomingEvents(entityID uint) ([]Event, error) {
	var events []Event
	err := r.DB.
		Where("entity_id = ? AND event_date >= CURRENT_DATE", entityID).
		Order("event_date ASC").
		Limit(5).
		Find(&events).Error
	return events, err
}

// ListEventsByEntity returns all events for a temple (with optional active filter)
func (r *Repository) ListEventsByEntity(entityID uint) ([]Event, error) {
	var events []Event
	err := r.DB.
		Where("entity_id = ?", entityID).
		Order("event_date ASC").
		Find(&events).Error
	return events, err
}

// UpdateEvent updates the event with new fields
func (r *Repository) UpdateEvent(e *Event) error {
	return r.DB.Model(&Event{}).Where("id = ?", e.ID).Updates(e).Error
}

// DeleteEvent removes the event by ID
func (r *Repository) DeleteEvent(id uint) error {
	return r.DB.Delete(&Event{}, id).Error
}
