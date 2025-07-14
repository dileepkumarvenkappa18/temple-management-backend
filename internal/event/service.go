package event

import "errors"

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
	event := &Event{
		Title:                req.Title,
		Description:          req.Description,
		EventDate:            req.EventDate,
		EventTime:            req.EventTime,
		EndDate:              req.EndDate,
		EndTime:              req.EndTime,
		Location:             req.Location,
		ImageURL:             req.ImageURL,
		MaxAttendees:         req.MaxAttendees,
		RegistrationRequired: req.RegistrationRequired,
		Category:             req.Category, // ✅ added
		IsActive:             true,
		CreatedBy:            createdBy,
		EntityID:             entityID,
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
// 📄 List Events with Pagination + Filters
func (s *Service) ListEventsByEntity(entityID uint, limit, offset int, category, search, status string) ([]Event, error) {
	events, err := s.Repo.ListEventsByEntity(entityID, limit, offset, category, search, status)
	if err != nil {
		return nil, err
	}

	for i := range events {
		count, _ := s.Repo.CountRSVPs(events[i].ID) // ignore error for now
		events[i].RSVPCount = count
	}

	return events, nil
}




// ===========================
// 🛠 Update Event (with ownership check)
func (s *Service) UpdateEvent(id uint, req *UpdateEventRequest, entityID uint) error {
	// Ensure ownership
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event.EntityID != entityID {
		return errors.New("unauthorized: cannot update this event")
	}

	// ✅ update with category support
	event.Title = req.Title
	event.Description = req.Description
	event.EventDate = req.EventDate
	event.EventTime = req.EventTime
	event.EndDate = req.EndDate
	event.EndTime = req.EndTime
	event.Location = req.Location
	event.ImageURL = req.ImageURL
	event.MaxAttendees = req.MaxAttendees
	event.RegistrationRequired = req.RegistrationRequired
	event.Category = req.Category

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
