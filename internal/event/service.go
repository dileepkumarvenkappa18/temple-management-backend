package event

type Service struct {
	Repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{Repo: r}
}

// CreateEvent creates a new temple event
func (s *Service) CreateEvent(e *Event) error {
	return s.Repo.CreateEvent(e)
}

// GetEventByID retrieves a single event by its ID
func (s *Service) GetEventByID(id uint) (*Event, error) {
	return s.Repo.GetEventByID(id)
}

// GetUpcomingEvents fetches upcoming events for an entity
func (s *Service) GetUpcomingEvents(entityID uint) ([]Event, error) {
	return s.Repo.GetUpcomingEvents(entityID)
}

// ListEventsByEntity returns all events for a specific temple
func (s *Service) ListEventsByEntity(entityID uint) ([]Event, error) {
	return s.Repo.ListEventsByEntity(entityID)
}

// UpdateEvent updates an existing event
func (s *Service) UpdateEvent(e *Event) error {
	return s.Repo.UpdateEvent(e)
}

// DeleteEvent deletes an event by ID
func (s *Service) DeleteEvent(id uint) error {
	return s.Repo.DeleteEvent(id)
}
