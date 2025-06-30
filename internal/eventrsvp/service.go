package eventrsvp

import "github.com/sharath018/temple-management-backend/internal/event"

type Service struct {
	Repo         *Repository
	EventService *event.Service
}

// Updated constructor to accept eventService
func NewService(repo *Repository, eventService *event.Service) *Service {
	return &Service{
		Repo:         repo,
		EventService: eventService,
	}
}

func (s *Service) CreateRSVP(rsvp *RSVP) error {
	return s.Repo.CreateRSVP(rsvp)
}

func (s *Service) GetMyRSVPs(userID uint) ([]RSVP, error) {
	return s.Repo.GetMyRSVPs(userID)
}

func (s *Service) GetRSVPsByEvent(eventID uint) ([]RSVP, error) {
	return s.Repo.GetRSVPsByEvent(eventID)
}

func (s *Service) UpdateRSVPStatus(eventID, userID uint, status, notes string) error {
	return s.Repo.UpdateRSVPStatus(eventID, userID, status, notes)
}
