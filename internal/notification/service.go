package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	repo     Repository
	authRepo auth.Repository
	auditSvc auditlog.Service
	email    Channel
	sms      Channel
	whatsapp Channel
}

// Constructor
func NewService(repo Repository, authRepo auth.Repository, cfg *config.Config, auditSvc auditlog.Service) Service {
	return &service{
		repo:     repo,
		authRepo: authRepo,
		auditSvc: auditSvc,
		email:    NewEmailSender(cfg),
		sms:      NewSMSChannel(),
		whatsapp: NewWhatsAppChannel(),
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
	// ✅ CRITICAL FIX: Return user-friendly error when no recipients
	if len(recipients) == 0 {
		log.Printf("⚠️ WARNING: No recipients found for entity %d", entityID)
		
		// Create log entry to track the attempt
		recipientsJSON, _ := json.Marshal([]string{})
		errorMsg := "No recipients found"
		logEntry := &NotificationLog{
			UserID:     senderID,
			EntityID:   entityID,
			TemplateID: templateID,
			Channel:    channel,
			Subject:    subject,
			Body:       body,
			Recipients: datatypes.JSON(recipientsJSON),
			Status:     "failed",
			Error:      &errorMsg,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		// Save failed attempt
		if err := s.repo.CreateNotificationLog(ctx, logEntry); err != nil {
			log.Printf("❌ Failed to create notification log: %v", err)
		}
		
		// Audit logging
		auditAction := map[string]string{
			"email":    "EMAIL_FAILED",
			"sms":      "SMS_FAILED",
			"whatsapp": "WHATSAPP_FAILED",
		}[channel]
		if auditAction == "" {
			auditAction = "NOTIFICATION_FAILED"
		}
		
		details := map[string]interface{}{
			"channel":          channel,
			"recipients_count": 0,
			"template_id":      templateID,
			"subject":          subject,
			"reason":           "no recipients found",
		}
		
		_ = s.auditSvc.LogAction(ctx, &senderID, &entityID, auditAction, details, ip, "failure")
		
		// Return specific error that frontend can display
		return errors.New("no recipients found - please add devotees or volunteers to this temple first")
	}

	// Validate SMTP configuration for email channel
	if channel == "email" {
		if err := s.validateEmailConfig(); err != nil {
			log.Printf("❌ Email configuration error: %v", err)
			return fmt.Errorf("email service not configured: %v", err)
		}
	}

	recipientsJSON, _ := json.Marshal(recipients)
	logEntry := &NotificationLog{
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

	if err := s.repo.CreateNotificationLog(ctx, logEntry); err != nil {
		return fmt.Errorf("failed to create notification log: %w", err)
	}

	// ✅ Enhanced error handling for each channel
	var sendErr error
	switch channel {
	case "email":
		// Send emails asynchronously but capture errors
		go func() {
			log.Printf("📧 Starting email send to %d recipients", len(recipients))
			if err := s.email.Send(recipients, subject, body); err != nil {
				log.Printf("❌ Email send error: %v", err)
				// Update log status to failed
				errMsg := err.Error()
				logEntry.Status = "failed"
				logEntry.Error = &errMsg
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			} else {
				log.Printf("✅ Email sent successfully to %d recipients", len(recipients))
				// Update log status to sent
				logEntry.Status = "sent"
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			}
		}()

	case "sms":
		go func() {
			log.Printf("📱 Starting SMS send to %d recipients", len(recipients))
			if err := s.sms.Send(recipients, subject, body); err != nil {
				log.Printf("❌ SMS send error: %v", err)
				errMsg := err.Error()
				logEntry.Status = "failed"
				logEntry.Error = &errMsg
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			} else {
				log.Printf("✅ SMS sent successfully")
				logEntry.Status = "sent"
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			}
		}()

	case "whatsapp":
		go func() {
			log.Printf("💬 Starting WhatsApp send to %d recipients", len(recipients))
			if err := s.whatsapp.Send(recipients, subject, body); err != nil {
				log.Printf("❌ WhatsApp send error: %v", err)
				errMsg := err.Error()
				logEntry.Status = "failed"
				logEntry.Error = &errMsg
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			} else {
				log.Printf("✅ WhatsApp sent successfully")
				logEntry.Status = "sent"
				logEntry.UpdatedAt = time.Now()
				_ = s.repo.UpdateNotificationLog(context.Background(), logEntry)
			}
		}()

	default:
		return fmt.Errorf("unsupported channel: %s", channel)
	}

	// ✅ Return success immediately (async sending)
	// Mark as "processing" initially
	logEntry.Status = "processing"
	logEntry.UpdatedAt = time.Now()
	if err := s.repo.UpdateNotificationLog(ctx, logEntry); err != nil {
		log.Printf("⚠️ Failed to update log status: %v", err)
	}

	// ===== Audit logging =====
	auditAction := map[string]string{
		"email":    "EMAIL_QUEUED",
		"sms":      "SMS_QUEUED",
		"whatsapp": "WHATSAPP_QUEUED",
	}[channel]
	if auditAction == "" {
		auditAction = "NOTIFICATION_QUEUED"
	}

	details := map[string]interface{}{
		"channel":          channel,
		"recipients_count": len(recipients),
		"template_id":      templateID,
		"subject":          subject,
	}

	if err := s.auditSvc.LogAction(ctx, &senderID, &entityID, auditAction, details, ip, "success"); err != nil {
		log.Printf("❌ Audit log error: %v", err)
	}

	return sendErr
}

// ✅ Add email configuration validator
func (s *service) validateEmailConfig() error {
	emailSender, ok := s.email.(*EmailSender)
	if !ok {
		return errors.New("email sender not properly initialized")
	}
	
	if emailSender.Host == "" {
		return errors.New("SMTP host not configured")
	}
	if emailSender.Port == "" {
		return errors.New("SMTP port not configured")
	}
	if emailSender.Username == "" {
		return errors.New("SMTP username not configured")
	}
	if emailSender.Password == "" {
		return errors.New("SMTP password not configured")
	}
	if emailSender.FromAddr == "" {
		return errors.New("SMTP from address not configured")
	}
	
	return nil
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
	log.Printf("🔍 GetEmailsByAudience called with entityID=%d, audience=%s", entityID, audience)

	switch audience {
	case "devotees":
		emails, err := s.authRepo.GetUserEmailsByRole("devotee", entityID)
		log.Printf("📧 Devotees query result: %d emails, error: %v", len(emails), err)
		return emails, err
	case "volunteers":
		emails, err := s.authRepo.GetUserEmailsByRole("volunteer", entityID)
		log.Printf("📧 Volunteers query result: %d emails, error: %v", len(emails), err)
		return emails, err
	case "all":
		log.Printf("🔍 Fetching devotees for entity %d...", entityID)
		devotees, err1 := s.authRepo.GetUserEmailsByRole("devotee", entityID)
		log.Printf("📧 Devotees result: %d emails, error: %v", len(devotees), err1)

		log.Printf("🔍 Fetching volunteers for entity %d...", entityID)
		volunteers, err2 := s.authRepo.GetUserEmailsByRole("volunteer", entityID)
		log.Printf("📧 Volunteers result: %d emails, error: %v", len(volunteers), err2)

		if err1 != nil && err2 != nil {
			return nil, fmt.Errorf("failed to fetch both audiences: %v | %v", err1, err2)
		}
		if err1 != nil {
			log.Printf("⚠️ Only returning volunteers (devotees failed)")
			return volunteers, nil
		}
		if err2 != nil {
			log.Printf("⚠️ Only returning devotees (volunteers failed)")
			return devotees, nil
		}

		combined := append(devotees, volunteers...)
		log.Printf("✅ Combined result: %d total emails", len(combined))
		return combined, nil
	default:
		return nil, fmt.Errorf("invalid audience: %s", audience)
	}
}