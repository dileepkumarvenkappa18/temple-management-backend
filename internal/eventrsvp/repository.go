package eventrsvp

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// CreateRSVP inserts a new RSVP
func (r *Repository) CreateRSVP(rsvp *RSVP) error {
	return r.DB.Create(rsvp).Error
}

// GetMyRSVPs lists all RSVPs by a user
func (r *Repository) GetMyRSVPs(userID uint) ([]RSVP, error) {
	var rsvps []RSVP
	err := r.DB.Where("user_id = ?", userID).Order("rsvp_date desc").Find(&rsvps).Error
	return rsvps, err
}

// GetRSVPsByEvent lists all RSVPs for a specific event (admin view)
func (r *Repository) GetRSVPsByEvent(eventID uint) ([]RSVP, error) {
	var rsvps []RSVP
	err := r.DB.Where("event_id = ?", eventID).Find(&rsvps).Error
	return rsvps, err
}

// UpdateRSVP allows updating the RSVP status or notes
func (r *Repository) UpdateRSVPStatus(eventID, userID uint, status, notes string) error {
	result := r.DB.Model(&RSVP{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Updates(map[string]interface{}{
			"status": status,
			"notes":  notes,
		})

	if result.RowsAffected == 0 {
		return errors.New("no RSVP found to update")
	}
	return result.Error
}