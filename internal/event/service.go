package event

import (
	"errors"
	"time"
)

// Service wraps business logic for temple events
type Service struct {
	Repo *Repository
}

// NewService initializes a new Service
func NewService(r *Repository) *Service {
	return &Service{Repo: r}
}

// ===========================
// 🎯 Create Event
func (s *Service) CreateEvent(req *CreateEventRequest, createdBy uint, entityID uint) error {
	// 🔄 Parse EventDate
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return errors.New("invalid event_date format. Use YYYY-MM-DD")
	}

	// 🔄 Parse EventTime (optional)
	var eventTimePtr *time.Time
	if req.EventTime != "" {
		parsedTime, err := time.Parse("15:04", req.EventTime)
		if err != nil {
			return errors.New("invalid event_time format. Use HH:MM in 24-hour format")
		}
		normalizedTime := time.Date(0, 1, 1, parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
		eventTimePtr = &normalizedTime
	}

	// 🛡 Handle optional IsActive safely
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	event := &Event{
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
		EventTime:   eventTimePtr,
		Location:    req.Location,
		EventType:   req.EventType,
		IsActive:    isActive,
		CreatedBy:   createdBy,
		EntityID:    entityID,
	}

	return s.Repo.CreateEvent(event)
}

// ===========================
// 🔍 Get Event by ID
func (s *Service) GetEventByID(id uint) (*Event, error) {
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	count, _ := s.Repo.CountRSVPs(event.ID)
	event.RSVPCount = count

	return event, nil
}

// ===========================
// 📆 Get Upcoming Events
func (s *Service) GetUpcomingEvents(entityID uint) ([]Event, error) {
	return s.Repo.GetUpcomingEvents(entityID)
}

// ===========================
// 📄 List Events with Pagination
func (s *Service) ListEventsByEntity(entityID uint, limit, offset int, search string) ([]Event, error) {
	events, err := s.Repo.ListEventsByEntity(entityID, limit, offset, search)
	if err != nil {
		return nil, err
	}

	for i := range events {
		count, _ := s.Repo.CountRSVPs(events[i].ID)
		events[i].RSVPCount = count
	}

	return events, nil
}

// ===========================
// 📊 Dashboard Stats
func (s *Service) GetEventStats(entityID uint) (*EventStatsResponse, error) {
	return s.Repo.GetEventStats(entityID)
}

// ===========================
// 🛠 Update Event (with ownership check)
func (s *Service) UpdateEvent(id uint, req *UpdateEventRequest, entityID uint) error {
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event.EntityID != entityID {
		return errors.New("unauthorized: cannot update this event")
	}

	// 🔄 Update fields
	event.Title = req.Title
	event.Description = req.Description

	// 🔄 Parse EventDate
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return errors.New("invalid event_date format. Use YYYY-MM-DD")
	}
	event.EventDate = eventDate

	// 🔄 Parse EventTime
	if req.EventTime != "" {
		parsedTime, err := time.Parse("15:04", req.EventTime)
		if err != nil {
			return errors.New("invalid event_time format. Use HH:MM in 24-hour format")
		}
		normalizedTime := time.Date(0, 1, 1, parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
		event.EventTime = &normalizedTime
	} else {
		event.EventTime = nil
	}

	event.Location = req.Location
	event.EventType = req.EventType

	if req.IsActive != nil {
		event.IsActive = *req.IsActive
	}

	return s.Repo.UpdateEvent(id, entityID, req)
}

// ===========================
// ❌ Delete Event (with ownership check)
func (s *Service) DeleteEvent(id uint, entityID uint) error {
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event.EntityID != entityID {
		return errors.New("unauthorized: cannot delete this event")
	}

	return s.Repo.DeleteEvent(id, entityID)
}
