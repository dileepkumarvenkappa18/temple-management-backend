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
	repo     Repository
	authRepo auth.Repository
	auditSvc auditlog.Service
	email    Channel
	sms      Channel
	whatsapp Channel
}

// ✅ Updated constructor to accept audit service
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

// ✅ COMPLETELY REWRITTEN with batch processing and async support
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

	// Create notification log entry first
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

	fmt.Printf("📨 Starting notification send: channel=%s, recipients=%d\n", channel, len(recipients))

	// ✅ Process notifications based on channel with batch support
	var sendErr error
	batchSize := 50 // Send 50 at a time to avoid overwhelming services
	
	switch channel {
	case "email":
		sendErr = s.sendEmailInBatches(recipients, subject, body, batchSize)
	case "sms":
		sendErr = s.sendSMSInBatches(recipients, subject, body, batchSize)
	case "whatsapp":
		sendErr = s.sendWhatsAppInBatches(recipients, subject, body, batchSize)
	default:
		sendErr = fmt.Errorf("unsupported channel: %s", channel)
	}

	// Update notification log status
	if sendErr != nil {
		errMsg := sendErr.Error()
		log.Status = "failed"
		log.Error = &errMsg
		fmt.Printf("❌ Notification send failed: %v\n", sendErr)
	} else {
		log.Status = "sent"
		fmt.Printf("✅ Notification sent successfully to %d recipients\n", len(recipients))
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

// ✅ NEW: Helper function to send emails in batches
func (s *service) sendEmailInBatches(recipients []string, subject, body string, batchSize int) error {
	totalRecipients := len(recipients)
	var lastErr error
	successCount := 0
	failedCount := 0
	
	fmt.Printf("📧 Sending emails in batches of %d (total: %d)\n", batchSize, totalRecipients)
	
	for i := 0; i < totalRecipients; i += batchSize {
		end := i + batchSize
		if end > totalRecipients {
			end = totalRecipients
		}
		
		batch := recipients[i:end]
		batchNum := (i / batchSize) + 1
		totalBatches := (totalRecipients + batchSize - 1) / batchSize
		
		fmt.Printf("📤 Processing batch %d/%d: sending to %d recipients\n", 
			batchNum, totalBatches, len(batch))
		
		if err := s.email.Send(batch, subject, body); err != nil {
			fmt.Printf("❌ Batch %d/%d failed: %v\n", batchNum, totalBatches, err)
			lastErr = err
			failedCount += len(batch)
			// Continue with next batch instead of stopping
		} else {
			successCount += len(batch)
			fmt.Printf("✅ Batch %d/%d sent successfully\n", batchNum, totalBatches)
		}
		
		// Small delay between batches to avoid rate limiting
		if end < totalRecipients {
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	fmt.Printf("📊 Email send complete: %d succeeded, %d failed out of %d total\n", 
		successCount, failedCount, totalRecipients)
	
	// If some emails succeeded, consider it a partial success
	if successCount > 0 && failedCount > 0 {
		return fmt.Errorf("partial success: %d/%d emails sent, last error: %v", 
			successCount, totalRecipients, lastErr)
	}
	
	// If all failed, return the last error
	if failedCount == totalRecipients && lastErr != nil {
		return fmt.Errorf("all batches failed: %v", lastErr)
	}
	
	return nil
}

// ✅ NEW: Helper function to send SMS in batches
func (s *service) sendSMSInBatches(recipients []string, subject, body string, batchSize int) error {
	totalRecipients := len(recipients)
	var lastErr error
	successCount := 0
	failedCount := 0
	
	fmt.Printf("📱 Sending SMS in batches of %d (total: %d)\n", batchSize, totalRecipients)
	
	for i := 0; i < totalRecipients; i += batchSize {
		end := i + batchSize
		if end > totalRecipients {
			end = totalRecipients
		}
		
		batch := recipients[i:end]
		batchNum := (i / batchSize) + 1
		totalBatches := (totalRecipients + batchSize - 1) / batchSize
		
		fmt.Printf("📤 Processing SMS batch %d/%d: sending to %d recipients\n", 
			batchNum, totalBatches, len(batch))
		
		if err := s.sms.Send(batch, subject, body); err != nil {
			fmt.Printf("❌ SMS Batch %d/%d failed: %v\n", batchNum, totalBatches, err)
			lastErr = err
			failedCount += len(batch)
		} else {
			successCount += len(batch)
			fmt.Printf("✅ SMS Batch %d/%d sent successfully\n", batchNum, totalBatches)
		}
		
		if end < totalRecipients {
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	fmt.Printf("📊 SMS send complete: %d succeeded, %d failed out of %d total\n", 
		successCount, failedCount, totalRecipients)
	
	if successCount > 0 && failedCount > 0 {
		return fmt.Errorf("partial success: %d/%d SMS sent, last error: %v", 
			successCount, totalRecipients, lastErr)
	}
	
	if failedCount == totalRecipients && lastErr != nil {
		return fmt.Errorf("all SMS batches failed: %v", lastErr)
	}
	
	return nil
}

// ✅ NEW: Helper function to send WhatsApp in batches
func (s *service) sendWhatsAppInBatches(recipients []string, subject, body string, batchSize int) error {
	totalRecipients := len(recipients)
	var lastErr error
	successCount := 0
	failedCount := 0
	
	fmt.Printf("💬 Sending WhatsApp in batches of %d (total: %d)\n", batchSize, totalRecipients)
	
	for i := 0; i < totalRecipients; i += batchSize {
		end := i + batchSize
		if end > totalRecipients {
			end = totalRecipients
		}
		
		batch := recipients[i:end]
		batchNum := (i / batchSize) + 1
		totalBatches := (totalRecipients + batchSize - 1) / batchSize
		
		fmt.Printf("📤 Processing WhatsApp batch %d/%d: sending to %d recipients\n", 
			batchNum, totalBatches, len(batch))
		
		if err := s.whatsapp.Send(batch, subject, body); err != nil {
			fmt.Printf("❌ WhatsApp Batch %d/%d failed: %v\n", batchNum, totalBatches, err)
			lastErr = err
			failedCount += len(batch)
		} else {
			successCount += len(batch)
			fmt.Printf("✅ WhatsApp Batch %d/%d sent successfully\n", batchNum, totalBatches)
		}
		
		if end < totalRecipients {
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	fmt.Printf("📊 WhatsApp send complete: %d succeeded, %d failed out of %d total\n", 
		successCount, failedCount, totalRecipients)
	
	if successCount > 0 && failedCount > 0 {
		return fmt.Errorf("partial success: %d/%d WhatsApp sent, last error: %v", 
			successCount, totalRecipients, lastErr)
	}
	
	if failedCount == totalRecipients && lastErr != nil {
		return fmt.Errorf("all WhatsApp batches failed: %v", lastErr)
	}
	
	return nil
}

// CreateInAppNotification stores a bell notification for a specific user
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

	// Publish to Redis channel for realtime SSE subscribers
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

// CreateInAppForEntityRoles sends an in-app notification to all users of given roles within an entity
func (s *service) CreateInAppForEntityRoles(ctx context.Context, entityID uint, roleNames []string, title, message, category string) error {
	// Collect unique user IDs
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
	// Create notifications for each
	for uid := range unique {
		if err := s.CreateInAppNotification(ctx, uid, entityID, title, message, category); err != nil {
			// continue on error to avoid blocking others
			fmt.Printf("in-app fanout error for user %d: %v\n", uid, err)
		}
	}
	return nil
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