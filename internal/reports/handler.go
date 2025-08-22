package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
)

// Handler holds service & repo (repo used for entity lookups here)
type Handler struct {
	service  ReportService
	repo     ReportRepository
	auditSvc auditlog.Service
}

// NewHandler creates a new reports handler
func NewHandler(svc ReportService, repo ReportRepository, auditSvc auditlog.Service) *Handler {
	return &Handler{
		service:  svc,
		repo:     repo,
		auditSvc: auditSvc,
	}
}

// GetActivities handles requests for the activities report
func (h *Handler) GetActivities(c *gin.Context) {
	// get logged-in user (AuthMiddleware already ran)
	userVal, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user object"})
		return
	}

	// Get IP address from context (set by AuditMiddleware)
	ip := middleware.GetIPFromContext(c)

	// reports/handler.go - in GetActivities method
	entityParam := c.Param("id") // instead of "entity_id"
	// either "all" or numeric id
	reportType := c.Query("type")
	if reportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type query param required: events|sevas|bookings"})
		return
	}
	dateRange := c.Query("date_range")
	if dateRange == "" {
		dateRange = DateRangeWeekly
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	format := c.Query("format") // excel, csv, pdf -> if empty return JSON

	// compute start & end
	start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := ActivitiesReportRequest{
		EntityID:  entityParam,
		Type:      reportType,
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Format:    format,
	}

	// resolve entity IDs: if "all" -> fetch user's entities; else validate single entity ownership
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
			return
		}
		if len(ids) == 0 {
			c.JSON(http.StatusOK, gin.H{"data": ReportData{}})
			return
		}
		for _, id := range ids {
			entityIDs = append(entityIDs, fmt.Sprint(id))
		}
	} else {
		// parse numeric entity id
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		// verify user owns this entity (fetch user's entities and check presence)
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate entity"})
			return
		}
		owned := false
		for _, id := range ids {
			if id == uint(eid) {
				owned = true
				break
			}
		}
		if !owned {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
			return
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	req.EntityIDs = entityIDs

	// If no format -> return JSON preview
	if req.Format == "" {
		data, err := h.service.GetActivities(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Log report view (optional - for JSON preview)
		details := map[string]interface{}{
			"report_type": req.Type,
			"format":     "json_preview",
			"entity_ids": req.EntityIDs,
			"date_range": req.DateRange,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "TEMPLE_ACTIVITIES_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Else export file (format present)
	bytes, fname, mime, err := h.service.ExportActivities(c.Request.Context(), req, &user.ID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

func (h *Handler) GetTempleRegisteredReport(c *gin.Context) {
	// get logged-in user (AuthMiddleware already ran)
	userVal, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user object"})
		return
	}

	// Get IP address from context (set by AuditMiddleware)
	ip := middleware.GetIPFromContext(c)

	entityParam := c.Param("id") // "all" or specific entity id

	dateRange := c.Query("date_range")
	if dateRange == "" {
		dateRange = DateRangeWeekly
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	status := c.Query("status") // approve|rejected|pending
	format := c.Query("format")

	start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Resolve entity IDs same way as in GetActivities
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
			return
		}
		if len(ids) == 0 {
			c.JSON(http.StatusOK, gin.H{"data": []TempleRegisteredReportRow{}})
			return
		}
		for _, id := range ids {
			entityIDs = append(entityIDs, fmt.Sprint(id))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate entity"})
			return
		}
		owned := false
		for _, id := range ids {
			if id == uint(eid) {
				owned = true
				break
			}
		}
		if !owned {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
			return
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	req := TempleRegisteredReportRequest{
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Status:    status,
		Format:    format,
		EntityID:  entityParam,
	}
	
	// The 'format' query parameter determines the report type for the exporter.
	var reportType string
	switch format {
	case "excel":
		reportType = ReportTypeTempleRegisteredExcel
	case "pdf":
		reportType = ReportTypeTempleRegisteredPDF
	case "csv": // Explicitly handle csv for clarity
		reportType = ReportTypeTempleRegistered
	default:
		// If no format is specified, return JSON preview
		data, err := h.service.GetTempleRegisteredReport(req, entityIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Log report view (optional - for JSON preview)
		details := map[string]interface{}{
			"report_type": "temple_registered",
			"format":     "json_preview",
			"entity_ids": entityIDs,
			"status":     status,
			"date_range": dateRange,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "TEMPLE_REGISTER_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Export file (format is present)
	bytes, fname, mime, err := h.service.ExportTempleRegisteredReport(c.Request.Context(), req, entityIDs, reportType, &user.ID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

func (h *Handler) GetDevoteeBirthdaysReport(c *gin.Context) {
	// get logged-in user (AuthMiddleware already ran)
	userVal, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user object"})
		return
	}

	// Get IP address from context (set by AuditMiddleware)
	ip := middleware.GetIPFromContext(c)

	entityParam := c.Param("id") // "all" or specific entity id

	dateRange := c.Query("date_range")
	if dateRange == "" {
		dateRange = DateRangeWeekly
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	format := c.Query("format")

	start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Resolve entity IDs same way as in other handlers
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
			return
		}
		if len(ids) == 0 {
			c.JSON(http.StatusOK, gin.H{"data": []DevoteeBirthdayReportRow{}})
			return
		}
		for _, id := range ids {
			entityIDs = append(entityIDs, fmt.Sprint(id))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate entity"})
			return
		}
		owned := false
		for _, id := range ids {
			if id == uint(eid) {
				owned = true
				break
			}
		}
		if !owned {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
			return
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	req := DevoteeBirthdaysReportRequest{
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Format:    format,
		EntityID:  entityParam,
	}
	
	// The 'format' query parameter determines the report type for the exporter.
	var reportType string
	switch format {
	case "excel":
		reportType = ReportTypeDevoteeBirthdaysExcel
	case "pdf":
		reportType = ReportTypeDevoteeBirthdaysPDF
	case "csv": // Explicitly handle csv for clarity
		reportType = ReportTypeDevoteeBirthdays
	default:
		// If no format is specified, return JSON preview
		data, err := h.service.GetDevoteeBirthdaysReport(req, entityIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Log report view (optional - for JSON preview)
		details := map[string]interface{}{
			"report_type": "devotee_birthdays",
			"format":     "json_preview",
			"entity_ids": entityIDs,
			"date_range": dateRange,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "DEVOTEE_BIRTHDAYS_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Export file (format is present)
	bytes, fname, mime, err := h.service.ExportDevoteeBirthdaysReport(c.Request.Context(), req, entityIDs, reportType, &user.ID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}