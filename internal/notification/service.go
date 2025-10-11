package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"gorm.io/datatypes"
)

// Service interface - updated with IP parameters for audit logging
type Service interface {
	CreateTemplate(ctx context.Context, template *NotificationTemplate, ip string) error
	UpdateTemplate(ctx context.Context, template *NotificationTemplate, ip string) error
	GetTemplates(ctx context.Context, entityID uint) ([]NotificationTemplate, error)
	GetTemplateByID(ctx context.Context, id uint, entityID uint) (*NotificationTemplate, error)
	DeleteTemplate(ctx context.Context, id uint, entityID uint, userID uint, ip string) error
	SendNotification(ctx context.Context, senderID, entityID uint, templateID *uint, channel, subject, body string, recipients []string, ip string) error
	GetNotificationsByUser(ctx context.Context, userID uint) ([]NotificationLog, error)
	GetEmailsByAudience(entityID uint, audience string) ([]string, error)
}

type service struct {
	repo       Repository
	authRepo   auth.Repository
	auditSvc   auditlog.Service // ✅ NEW: Audit service
	email      Channel
	sms        Channel
	whatsapp   Channel
}

// ✅ Updated constructor to accept audit service
func NewService(repo Repository, authRepo auth.Repository, cfg *config.Config, auditSvc auditlog.Service) Service {
	return &service{
		repo:       repo,
		authRepo:   authRepo,
		auditSvc:   auditSvc, // ✅ injected audit service
		email:      NewEmailSender(cfg),
		sms:        NewSMSChannel(),
		whatsapp:   NewWhatsAppChannel(),
	}
}

// ✅ Updated with audit logging
func (s *service) CreateTemplate(ctx context.Context, t *NotificationTemplate, ip string) error {
	err := s.repo.CreateTemplate(ctx, t)
	
	// ✅ Audit log the template creation
	status := "success"
	if err != nil {
		status = "failure"
	}
	
	details := map[string]interface{}{
		"template_name": t.Name,
		"category":      t.Category,
	}
	
	auditErr := s.auditSvc.LogAction(ctx, &t.UserID, &t.EntityID, "TEMPLATE_CREATED", details, ip, status)
	if auditErr != nil {
		// Log audit error but don't fail the main operation
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}
	
	return err
}

// ✅ Updated with audit logging
func (s *service) UpdateTemplate(ctx context.Context, t *NotificationTemplate, ip string) error {
	err := s.repo.UpdateTemplate(ctx, t)
	
	// ✅ Audit log the template update
	status := "success"
	if err != nil {
		status = "failure"
	}
	
	details := map[string]interface{}{
		"template_id":   t.ID,
		"template_name": t.Name,
		"category":      t.Category,
	}
	
	auditErr := s.auditSvc.LogAction(ctx, &t.UserID, &t.EntityID, "TEMPLATE_UPDATED", details, ip, status)
	if auditErr != nil {
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}
	
	return err
}

func (s *service) GetTemplates(ctx context.Context, entityID uint) ([]NotificationTemplate, error) {
	return s.repo.GetTemplatesByEntity(ctx, entityID)
}

func (s *service) GetTemplateByID(ctx context.Context, id uint, entityID uint) (*NotificationTemplate, error) {
	return s.repo.GetTemplateByID(ctx, id, entityID)
}

// ✅ Updated with audit logging
func (s *service) DeleteTemplate(ctx context.Context, id uint, entityID uint, userID uint, ip string) error {
	// Get template details before deletion for audit
	template, getErr := s.repo.GetTemplateByID(ctx, id, entityID)
	templateName := "unknown"
	if getErr == nil && template != nil {
		templateName = template.Name
	}
	
	err := s.repo.DeleteTemplate(ctx, id, entityID)
	
	// ✅ Audit log the template deletion
	status := "success"
	if err != nil {
		status = "failure"
	}
	
	details := map[string]interface{}{
		"template_id":   id,
		"template_name": templateName,
	}
	
	auditErr := s.auditSvc.LogAction(ctx, &userID, &entityID, "TEMPLATE_DELETED", details, ip, status)
	if auditErr != nil {
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}
	
	return err
}

// ✅ Updated with audit logging
func (s *service) SendNotification(
	ctx context.Context,
	senderID, entityID uint,
	templateID *uint,
	channel, subject, body string,
	recipients []string,
	ip string,
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

	// Update notification log status
	if sendErr != nil {
		errMsg := sendErr.Error()
		log.Status = "failed"
		log.Error = &errMsg
	} else {
		log.Status = "sent"
	}

	log.UpdatedAt = time.Now()
	updateErr := s.repo.UpdateNotificationLog(ctx, log)

	// ✅ Audit log the notification send action
	auditAction := ""
	switch channel {
	case "email":
		auditAction = "EMAIL_SENT"
	case "sms":
		auditAction = "SMS_SENT"
	case "whatsapp":
		auditAction = "WHATSAPP_SENT"
	default:
		auditAction = "NOTIFICATION_SENT"
	}

	status := "success"
	if sendErr != nil {
		status = "failure"
	}

	details := map[string]interface{}{
		"channel":          channel,
		"recipients_count": len(recipients),
		"template_id":      templateID,
		"subject":          subject,
	}

	auditErr := s.auditSvc.LogAction(ctx, &senderID, &entityID, auditAction, details, ip, status)
	if auditErr != nil {
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}

	// Return the original send error if any, otherwise return update error
	if sendErr != nil {
		return sendErr
	}
	return updateErr
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