package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/utils"
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

	// In-app notifications
	CreateInAppNotification(ctx context.Context, userID, entityID uint, title, message, category string) error
	ListInAppByUser(ctx context.Context, userID uint, entityID *uint, limit int) ([]InAppNotification, error)
	MarkInAppAsRead(ctx context.Context, id uint, userID uint) error

	// Fan-out helpers
	CreateInAppForEntityRoles(ctx context.Context, entityID uint, roleNames []string, title, message, category string) error
}

type service struct {
	repo       Repository
	authRepo   auth.Repository
	auditSvc   auditlog.Service // ✅ Audit service
	email      Channel
	sms        Channel
	whatsapp   Channel
}

// Constructor
func NewService(repo Repository, authRepo auth.Repository, cfg *config.Config, auditSvc auditlog.Service) Service {
	return &service{
		repo:       repo,
		authRepo:   authRepo,
		auditSvc:   auditSvc,
		email:      NewEmailSender(cfg),
		sms:        NewSMSChannel(),
		whatsapp:   NewWhatsAppChannel(),
	}
}

// ================= Template Management =================
func (s *service) CreateTemplate(ctx context.Context, t *NotificationTemplate, ip string) error {
	err := s.repo.CreateTemplate(ctx, t)
	status := "success"
	if err != nil {
		status = "failure"
	}
	details := map[string]interface{}{
		"template_name": t.Name,
		"category":      t.Category,
	}
	if auditErr := s.auditSvc.LogAction(ctx, &t.UserID, &t.EntityID, "TEMPLATE_CREATED", details, ip, status); auditErr != nil {
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}
	return err
}

func (s *service) UpdateTemplate(ctx context.Context, t *NotificationTemplate, ip string) error {
	err := s.repo.UpdateTemplate(ctx, t)
	status := "success"
	if err != nil {
		status = "failure"
	}
	details := map[string]interface{}{
		"template_id":   t.ID,
		"template_name": t.Name,
		"category":      t.Category,
	}
	if auditErr := s.auditSvc.LogAction(ctx, &t.UserID, &t.EntityID, "TEMPLATE_UPDATED", details, ip, status); auditErr != nil {
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

func (s *service) DeleteTemplate(ctx context.Context, id uint, entityID uint, userID uint, ip string) error {
	template, getErr := s.repo.GetTemplateByID(ctx, id, entityID)
	templateName := "unknown"
	if getErr == nil && template != nil {
		templateName = template.Name
	}
	err := s.repo.DeleteTemplate(ctx, id, entityID)
	status := "success"
	if err != nil {
		status = "failure"
	}
	details := map[string]interface{}{
		"template_id":   id,
		"template_name": templateName,
	}
	if auditErr := s.auditSvc.LogAction(ctx, &userID, &entityID, "TEMPLATE_DELETED", details, ip, status); auditErr != nil {
		fmt.Printf("❌ Audit log error: %v\n", auditErr)
	}
	return err
}

// ================= Notification Sending =================
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

	// ===== Asynchronous sending =====
	switch channel {
	case "email":
		utils.SendBulkEmailsAsync(recipients, subject, body)
	case "sms":
		go func() {
			if err := s.sms.Send(recipients, subject, body); err != nil {
				fmt.Printf("❌ SMS send error: %v\n", err)
			}
		}()
	case "whatsapp":
		go func() {
			if err := s.whatsapp.Send(recipients, subject, body); err != nil {
				fmt.Printf("❌ WhatsApp send error: %v\n", err)
			}
		}()
	default:
		return fmt.Errorf("unsupported channel: %s", channel)
	}

	// Mark log as sent immediately
	log.Status = "sent"
	log.UpdatedAt = time.Now()
	updateErr := s.repo.UpdateNotificationLog(ctx, log)

	// ===== Audit logging =====
	auditAction := map[string]string{
		"email":    "EMAIL_SENT",
		"sms":      "SMS_SENT",
		"whatsapp": "WHATSAPP_SENT",
	}[channel]
	if auditAction == "" {
		auditAction = "NOTIFICATION_SENT"
	}

	details := map[string]interface{}{
		"channel":          channel,
		"recipients_count": len(recipients),
		"template_id":      templateID,
		"subject":          subject,
	}

	if err := s.auditSvc.LogAction(ctx, &senderID, &entityID, auditAction, details, ip, "success"); err != nil {
		fmt.Printf("❌ Audit log error: %v\n", err)
	}

	return updateErr
}

// ================= In-App Notifications =================
func (s *service) CreateInAppNotification(ctx context.Context, userID, entityID uint, title, message, category string) error {
	item := &InAppNotification{
		UserID:    userID,
		EntityID:  entityID,
		Title:     title,
		Message:   message,
		Category:  category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateInApp(ctx, item); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"id":         item.ID,
		"user_id":    item.UserID,
		"entity_id":  item.EntityID,
		"title":      item.Title,
		"message":    item.Message,
		"category":   item.Category,
		"is_read":    item.IsRead,
		"created_at": item.CreatedAt,
	})
	channel := fmt.Sprintf("notifications:user:%d", userID)
	_ = utils.RedisClient.Publish(utils.Ctx, channel, string(payload)).Err()
	return nil
}

func (s *service) ListInAppByUser(ctx context.Context, userID uint, entityID *uint, limit int) ([]InAppNotification, error) {
	return s.repo.ListInAppByUser(ctx, userID, entityID, limit)
}

func (s *service) MarkInAppAsRead(ctx context.Context, id uint, userID uint) error {
	return s.repo.MarkInAppAsRead(ctx, id, userID)
}

func (s *service) CreateInAppForEntityRoles(ctx context.Context, entityID uint, roleNames []string, title, message, category string) error {
	unique := make(map[uint]struct{})
	for _, role := range roleNames {
		ids, err := s.authRepo.GetUserIDsByRole(role, entityID)
		if err != nil {
			return err
		}
		for _, id := range ids {
			unique[id] = struct{}{}
		}
	}
	for uid := range unique {
		if err := s.CreateInAppNotification(ctx, uid, entityID, title, message, category); err != nil {
			fmt.Printf("in-app fanout error for user %d: %v\n", uid, err)
		}
	}
	return nil
}

// ================= Fetch Notifications =================
func (s *service) GetNotificationsByUser(ctx context.Context, userID uint) ([]NotificationLog, error) {
	return s.repo.GetNotificationsByUser(ctx, userID)
}

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
