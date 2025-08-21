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

// resolveEntityIDs resolves entity IDs based on user role and parameters
func (h *Handler) resolveEntityIDs(c *gin.Context, user auth.User, entityParam string) ([]string, error) {
	// Check if user is superadmin
	if user.Role.RoleName == "super_admin" {
		return h.resolveSuperadminEntities(c, entityParam)
	}

	// Existing templeadmin logic
	return h.resolveTempleadminEntities(user, entityParam)
}

// resolveSuperadminEntities handles entity resolution for superadmin
func (h *Handler) resolveSuperadminEntities(c *gin.Context, entityParam string) ([]string, error) {
	var entityIDs []string

	if strings.ToLower(entityParam) == "all" {
		// Get selected templeadmin IDs from query parameter
		selectedTempleadminsParam := c.Query("selected_templeadmins")
		if selectedTempleadminsParam == "" {
			return nil, fmt.Errorf("selected_templeadmins query parameter is required for superadmin when using 'all'")
		}

		// Parse comma-separated templeadmin IDs
		templeadminIDStrs := strings.Split(selectedTempleadminsParam, ",")
		var templeadminIDs []uint
		for _, idStr := range templeadminIDStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid templeadmin ID: %s", idStr)
			}
			templeadminIDs = append(templeadminIDs, uint(id))
		}

		if len(templeadminIDs) == 0 {
			return nil, fmt.Errorf("at least one templeadmin must be selected")
		}

		// 🔧 FIX: Use GetEntitiesByMultipleTenants instead of GetEntitiesByTenant
		ids, err := h.repo.GetEntitiesByMultipleTenants(templeadminIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch entities for selected templeadmins: %v", err)
		}

		// Convert to string slice
		for _, id := range ids {
			entityIDs = append(entityIDs, fmt.Sprint(id))
		}

		// Add debug logging
		fmt.Printf("DEBUG: SuperAdmin selected templeadmins: %v, found entities: %v\n", templeadminIDs, entityIDs)
	} else {
		// Single entity - validate it belongs to one of the selected templeadmins
		eid, err := strconv.ParseUint(entityParam, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid entity_id path param: %s", entityParam)
		}

		selectedTempleadminsParam := c.Query("selected_templeadmins")
		if selectedTempleadminsParam == "" {
			return nil, fmt.Errorf("selected_templeadmins query parameter is required for superadmin")
		}

		templeadminIDStrs := strings.Split(selectedTempleadminsParam, ",")
		var templeadminIDs []uint
		for _, idStr := range templeadminIDStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid templeadmin ID: %s", idStr)
			}
			templeadminIDs = append(templeadminIDs, uint(id))
		}

		if len(templeadminIDs) == 0 {
			return nil, fmt.Errorf("at least one templeadmin must be selected")
		}

		// Verify entity belongs to one of the selected templeadmins
		owned, err := h.repo.ValidateEntityOwnershipByMultipleTenants(uint(eid), templeadminIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to validate entity ownership: %v", err)
		}
		if !owned {
			return nil, fmt.Errorf("entity %d not owned by selected templeadmins", eid)
		}

		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	return entityIDs, nil
}

// resolveTempleadminEntities handles entity resolution for templeadmin (existing logic)
func (h *Handler) resolveTempleadminEntities(user auth.User, entityParam string) ([]string, error) {
	var entityIDs []string

	if strings.ToLower(entityParam) == "all" {
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user entities: %v", err)
		}
		for _, id := range ids {
			entityIDs = append(entityIDs, fmt.Sprint(id))
		}
	} else {
		// parse numeric entity id
		eid, err := strconv.ParseUint(entityParam, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid entity_id path param: %s", entityParam)
		}
		// verify user owns this entity
		ids, err := h.repo.GetEntitiesByTenant(user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate entity: %v", err)
		}
		owned := false
		for _, id := range ids {
			if id == uint(eid) {
				owned = true
				break
			}
		}
		if !owned {
			return nil, fmt.Errorf("not authorized for entity %d", eid)
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	return entityIDs, nil
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

	entityParam := c.Param("id") // instead of "entity_id"
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

	// Resolve entity IDs based on user role
	entityIDs, err := h.resolveEntityIDs(c, user, entityParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(entityIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": ReportData{},
			"message": "No entities found for the selected criteria",
		})
		return
	}

	req := ActivitiesReportRequest{
		EntityID:  entityParam,
		Type:      reportType,
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Format:    format,
		EntityIDs: entityIDs,
	}

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
			"user_role":  user.Role.RoleName,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "TEMPLE_ACTIVITIES_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, gin.H{
			"data": data,
			"metadata": gin.H{
				"total_entities": len(entityIDs),
				"entity_ids":     entityIDs,
				"date_range":     dateRange,
				"start_date":     start.Format("2006-01-02"),
				"end_date":       end.Format("2006-01-02"),
			},
		})
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

	// Resolve entity IDs based on user role
	entityIDs, err := h.resolveEntityIDs(c, user, entityParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(entityIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": []TempleRegisteredReportRow{},
			"message": "No entities found for the selected criteria",
		})
		return
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
			"user_role":  user.Role.RoleName,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "TEMPLE_REGISTER_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, gin.H{
			"data": data,
			"metadata": gin.H{
				"total_entities": len(entityIDs),
				"entity_ids":     entityIDs,
				"status_filter":  status,
				"date_range":     dateRange,
				"start_date":     start.Format("2006-01-02"),
				"end_date":       end.Format("2006-01-02"),
			},
		})
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

	// Resolve entity IDs based on user role
	entityIDs, err := h.resolveEntityIDs(c, user, entityParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(entityIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": []DevoteeBirthdayReportRow{},
			"message": "No entities found for the selected criteria",
		})
		return
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
			"user_role":  user.Role.RoleName,
		}
		h.auditSvc.LogAction(c.Request.Context(), &user.ID, nil, "DEVOTEE_BIRTHDAYS_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, gin.H{
			"data": data,
			"metadata": gin.H{
				"total_entities": len(entityIDs),
				"entity_ids":     entityIDs,
				"date_range":     dateRange,
				"start_date":     start.Format("2006-01-02"),
				"end_date":       end.Format("2006-01-02"),
			},
		})
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

// GetTempleadminsList returns list of templeadmins for superadmin selection
func (h *Handler) GetTempleadminsList(c *gin.Context) {
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

	// Only superadmin can access this endpoint
	if user.Role.RoleName != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: superadmin only"})
		return
	}

	templeadmins, err := h.repo.GetAllTempleadmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch templeadmins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templeadmins})
}

// Add this method to your Handler struct for debugging
func (h *Handler) TestSuperadminData(c *gin.Context) {
	selectedTempleadmins := c.Query("selected_templeadmins")
	if selectedTempleadmins == "" {
		c.JSON(400, gin.H{"error": "selected_templeadmins required"})
		return
	}

	// Parse templeadmin IDs
	templeadminIDStrs := strings.Split(selectedTempleadmins, ",")
	var templeadminIDs []uint
	for _, idStr := range templeadminIDStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid templeadmin ID: %s", idStr)})
			return
		}
		templeadminIDs = append(templeadminIDs, uint(id))
	}

	// Get entities
	entities, err := h.repo.GetEntitiesByMultipleTenants(templeadminIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get counts
	var eventCount, sevaCount, bookingCount int64
if len(entities) > 0 {
    eventCount, _ = h.repo.CountEventsByEntities(entities)
    sevaCount, _ = h.repo.CountSevasByEntities(entities)
    bookingCount, _ = h.repo.CountBookingsByEntities(entities)
}


	c.JSON(200, gin.H{
		"templeadmin_ids": templeadminIDs,
		"entities_found": entities,
		"counts": gin.H{
			"entities": len(entities),
			"events":   eventCount,
			"sevas":    sevaCount,
			"bookings": bookingCount,
		},
	})
}