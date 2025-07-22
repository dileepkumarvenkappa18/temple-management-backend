package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"gorm.io/datatypes"
)

// Service interface
type Service interface {
	CreateTemplate(ctx context.Context, template *NotificationTemplate) error
	UpdateTemplate(ctx context.Context, template *NotificationTemplate) error
	GetTemplates(ctx context.Context, entityID uint) ([]NotificationTemplate, error)
	GetTemplateByID(ctx context.Context, id uint, entityID uint) (*NotificationTemplate, error)
	DeleteTemplate(ctx context.Context, id uint, entityID uint) error
	SendNotification(ctx context.Context, senderID, entityID uint, templateID *uint, channel, subject, body string, recipients []string) error
	GetNotificationsByUser(ctx context.Context, userID uint) ([]NotificationLog, error)
	GetEmailsByAudience(entityID uint, audience string) ([]string, error)
}

type service struct {
	repo     Repository
	authRepo auth.Repository // ✅ Injected auth repo
	email    Channel
	sms      Channel
	whatsapp Channel
}

// ✅ Updated constructor
func NewService(repo Repository, authRepo auth.Repository, cfg *config.Config) Service {
	return &service{
		repo:     repo,
		authRepo: authRepo,                  // ✅ assigned here
		email:    NewEmailSender(cfg),
		sms:      NewSMSChannel(),
		whatsapp: NewWhatsAppChannel(),
	}
}

func (s *service) CreateTemplate(ctx context.Context, t *NotificationTemplate) error {
	return s.repo.CreateTemplate(ctx, t)
}

func (s *service) UpdateTemplate(ctx context.Context, t *NotificationTemplate) error {
	return s.repo.UpdateTemplate(ctx, t)
}

func (s *service) GetTemplates(ctx context.Context, entityID uint) ([]NotificationTemplate, error) {
	return s.repo.GetTemplatesByEntity(ctx, entityID)
}

func (s *service) GetTemplateByID(ctx context.Context, id uint, entityID uint) (*NotificationTemplate, error) {
	return s.repo.GetTemplateByID(ctx, id, entityID)
}

func (s *service) DeleteTemplate(ctx context.Context, id uint, entityID uint) error {
	return s.repo.DeleteTemplate(ctx, id, entityID)
}

func (s *service) SendNotification(
	ctx context.Context,
	senderID, entityID uint,
	templateID *uint,
	channel, subject, body string,
	recipients []string,
) error {
	if len(recipients) == 0 {
		return errors.New("no recipients specified")
	}

	recipientsJSON, _ := json.Marshal(recipients)
	log := &NotificationLog{
		UserID:     senderID,
		EntityID:   entityID,
		TemplateID: templateID,
		Channel:    channel,
		Subject:    subject,
		Body:       body,
		Recipients: datatypes.JSON(recipientsJSON),
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateNotificationLog(ctx, log); err != nil {
		return err
	}

	var sendErr error
	switch channel {
	case "email":
		sendErr = s.email.Send(recipients, subject, body)
	case "sms":
		sendErr = s.sms.Send(recipients, subject, body)
	case "whatsapp":
		sendErr = s.whatsapp.Send(recipients, subject, body)
	default:
		sendErr = fmt.Errorf("unsupported channel: %s", channel)
	}

	if sendErr != nil {
		errMsg := sendErr.Error()
		log.Status = "failed"
		log.Error = &errMsg
	} else {
		log.Status = "sent"
	}

	log.UpdatedAt = time.Now()
	return s.repo.UpdateNotificationLog(ctx, log)
}

func (s *service) GetNotificationsByUser(ctx context.Context, userID uint) ([]NotificationLog, error) {
	return s.repo.GetNotificationsByUser(ctx, userID)
}

// ✅ Get Emails by audience using authRepo
func (s *service) GetEmailsByAudience(entityID uint, audience string) ([]string, error) {
	switch audience {
	case "devotees":
		return s.authRepo.GetUserEmailsByRole("devotee", entityID)
	case "volunteers":
		return s.authRepo.GetUserEmailsByRole("volunteer", entityID)
	case "all":
		devotees, err1 := s.authRepo.GetUserEmailsByRole("devotee", entityID)
		volunteers, err2 := s.authRepo.GetUserEmailsByRole("volunteer", entityID)

		if err1 != nil && err2 != nil {
			return nil, fmt.Errorf("failed to fetch both audiences: %v | %v", err1, err2)
		}
		if err1 != nil {
			return volunteers, nil
		}
		if err2 != nil {
			return devotees, nil
		}

		return append(devotees, volunteers...), nil
	default:
		return nil, fmt.Errorf("invalid audience: %s", audience)
	}
}
