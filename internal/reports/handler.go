package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
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
	// Get access context from middleware
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

	// Get IP address from context (set by AuditMiddleware)
	ip := middleware.GetIPFromContext(c)

	entityParam := c.Param("id") // either "all" or numeric id
	reportType := c.Query("type")
	if reportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type query param required: events|sevas|bookings|donations"})
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

	// resolve entity IDs based on access context
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		// For "all", get accessible entities based on user role
		if ctx.RoleName == "templeadmin" {
			// Templeadmin can access their own entities
			ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
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
			// standarduser/monitoringuser can only access their assigned entity
			accessibleEntityID := ctx.GetAccessibleEntityID()
			if accessibleEntityID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
				return
			}
			entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
		}
	} else {
		// parse numeric entity id
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		
		// verify user can access this specific entity
		if !h.canAccessEntity(ctx, uint(eid)) {
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
		h.auditSvc.LogAction(c.Request.Context(), &ctx.UserID, nil, "TEMPLE_ACTIVITIES_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Else export file (format present)
	bytes, fname, mime, err := h.service.ExportActivities(c.Request.Context(), req, &ctx.UserID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

func (h *Handler) GetTempleRegisteredReport(c *gin.Context) {
	// Get access context from middleware
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

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

	// Resolve entity IDs based on access context
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		if ctx.RoleName == "templeadmin" {
			ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
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
			accessibleEntityID := ctx.GetAccessibleEntityID()
			if accessibleEntityID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
				return
			}
			entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		
		if !h.canAccessEntity(ctx, uint(eid)) {
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
		h.auditSvc.LogAction(c.Request.Context(), &ctx.UserID, nil, "TEMPLE_REGISTER_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Export file (format is present)
	bytes, fname, mime, err := h.service.ExportTempleRegisteredReport(c.Request.Context(), req, entityIDs, reportType, &ctx.UserID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

func (h *Handler) GetDevoteeBirthdaysReport(c *gin.Context) {
	// Get access context from middleware
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

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

	// Resolve entity IDs based on access context
	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		if ctx.RoleName == "templeadmin" {
			ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
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
			accessibleEntityID := ctx.GetAccessibleEntityID()
			if accessibleEntityID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
				return
			}
			entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		
		if !h.canAccessEntity(ctx, uint(eid)) {
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
		h.auditSvc.LogAction(c.Request.Context(), &ctx.UserID, nil, "DEVOTEE_BIRTHDAYS_REPORT_VIEWED", details, ip, "success")
		
		c.JSON(http.StatusOK, data)
		return
	}

	// Export file (format is present)
	bytes, fname, mime, err := h.service.ExportDevoteeBirthdaysReport(c.Request.Context(), req, entityIDs, reportType, &ctx.UserID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

// GetDevoteeListReport handles requests for devotee list report
func (h *Handler) GetDevoteeListReport(c *gin.Context) {
	// Get access context from middleware
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

	ip := middleware.GetIPFromContext(c)

	entityParam := c.Param("id") // "all" or entity id

	dateRange := c.Query("date_range")
	if dateRange == "" {
		dateRange = DateRangeWeekly
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	status := c.Query("status") // active|inactive|blocked etc
	format := c.Query("format")

	start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		if ctx.RoleName == "templeadmin" {
			ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
				return
			}
			if len(ids) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": []DevoteeListReportRow{}})
				return
			}
			for _, id := range ids {
				entityIDs = append(entityIDs, fmt.Sprint(id))
			}
		} else {
			accessibleEntityID := ctx.GetAccessibleEntityID()
			if accessibleEntityID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
				return
			}
			entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		
		if !h.canAccessEntity(ctx, uint(eid)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
			return
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	req := DevoteeListReportRequest{
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Status:    status,
		Format:    format,
		EntityID:  entityParam,
	}

	var reportType string
	switch format {
	case "excel":
		reportType = ReportTypeDevoteeListExcel
	case "pdf":
		reportType = ReportTypeDevoteeListPDF
	case "csv":
		reportType = ReportTypeDevoteeListCSV
	default:
		data, err := h.service.GetDevoteeListReport(req, entityIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		details := map[string]interface{}{
			"report_type": "devotee_list",
			"format":     "json_preview",
			"entity_ids": entityIDs,
			"status":     status,
			"date_range": dateRange,
		}
		h.auditSvc.LogAction(c.Request.Context(), &ctx.UserID, nil, "DEVOTEE_LIST_REPORT_VIEWED", details, ip, "success")
		c.JSON(http.StatusOK, data)
		return
	}

	bytes, fname, mime, err := h.service.ExportDevoteeListReport(c.Request.Context(), req, entityIDs, reportType, &ctx.UserID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

// GetDevoteeProfileReport handles requests for devotee profile report
func (h *Handler) GetDevoteeProfileReport(c *gin.Context) {
	// Get access context from middleware
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

	ip := middleware.GetIPFromContext(c)

	entityParam := c.Param("id") // "all" or entity id

	dateRange := c.Query("date_range")
	if dateRange == "" {
		dateRange = DateRangeWeekly
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	status := c.Query("status") // active|inactive|blocked etc
	format := c.Query("format")

	start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entityIDs []string
	if strings.ToLower(entityParam) == "all" {
		if ctx.RoleName == "templeadmin" {
			ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
				return
			}
			if len(ids) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": []DevoteeProfileReportRow{}})
				return
			}
			for _, id := range ids {
				entityIDs = append(entityIDs, fmt.Sprint(id))
			}
		} else {
			accessibleEntityID := ctx.GetAccessibleEntityID()
			if accessibleEntityID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
				return
			}
			entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
		}
	} else {
		eid, err := strconv.ParseUint(entityParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
			return
		}
		
		if !h.canAccessEntity(ctx, uint(eid)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
			return
		}
		entityIDs = append(entityIDs, fmt.Sprint(eid))
	}

	req := DevoteeProfileReportRequest{
		DateRange: dateRange,
		StartDate: start,
		EndDate:   end,
		Status:    status,
		Format:    format,
		EntityID:  entityParam,
	}

	var reportType string
	switch format {
	case "excel":
		reportType = ReportTypeDevoteeProfileExcel
	case "pdf":
		reportType = ReportTypeDevoteeProfilePDF
	case "csv":
		reportType = ReportTypeDevoteeProfileCSV
	default:
		data, err := h.service.GetDevoteeProfileReport(req, entityIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		details := map[string]interface{}{
			"report_type": "devotee_profile",
			"format":     "json_preview",
			"entity_ids": entityIDs,
			"status":     status,
			"date_range": dateRange,
		}
		h.auditSvc.LogAction(c.Request.Context(), &ctx.UserID, nil, "DEVOTEE_PROFILE_REPORT_VIEWED", details, ip, "success")
		c.JSON(http.StatusOK, data)
		return
	}

	bytes, fname, mime, err := h.service.ExportDevoteeProfileReport(c.Request.Context(), req, entityIDs, reportType, &ctx.UserID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	c.Data(http.StatusOK, mime, bytes)
}

// canAccessEntity checks if the user can access a specific entity
func (h *Handler) canAccessEntity(ctx middleware.AccessContext, entityID uint) bool {
	if ctx.RoleName == "templeadmin" {
		// Templeadmin can access entities they created
		ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
		if err != nil {
			return false
		}
		for _, id := range ids {
			if id == entityID {
				return true
			}
		}
		return false
	} else {
		// standarduser/monitoringuser can only access their assigned entity
		accessibleEntityID := ctx.GetAccessibleEntityID()
		return accessibleEntityID != nil && *accessibleEntityID == entityID
	}
}

// GetAuditLogsReport handles requests for audit logs report
func (h *Handler) GetAuditLogsReport(c *gin.Context) {
    // 1️⃣ Get access context from middleware
    accessContext, exists := c.Get("access_context")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
        return
    }
    ctx := accessContext.(middleware.AccessContext)

    // 2️⃣ Get IP address from context
    ip := middleware.GetIPFromContext(c)

    // 3️⃣ Get request query/path params
    entityParam := c.Param("id") // "all" or specific entity id
    action := c.Query("action")
    status := c.Query("status")
    dateRange := c.Query("date_range")
    if dateRange == "" {
        dateRange = DateRangeWeekly // default weekly
    }
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    format := c.Query("format") // json preview, csv, excel, pdf

    // 4️⃣ Determine start and end dates
    start, end, err := GetDateRange(dateRange, startDateStr, endDateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 5️⃣ Resolve entity IDs based on access context
    var entityIDs []string
    if strings.ToLower(entityParam) == "all" {
        if ctx.RoleName == "templeadmin" {
            ids, err := h.repo.GetEntitiesByTenant(ctx.UserID)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user entities"})
                return
            }
            if len(ids) == 0 {
                c.JSON(http.StatusOK, gin.H{"data": []AuditLogReportRow{}})
                return
            }
            for _, id := range ids {
                entityIDs = append(entityIDs, fmt.Sprint(id))
            }
        } else {
            accessibleEntityID := ctx.GetAccessibleEntityID()
            if accessibleEntityID == nil {
                c.JSON(http.StatusForbidden, gin.H{"error": "no accessible entity"})
                return
            }
            entityIDs = append(entityIDs, fmt.Sprint(*accessibleEntityID))
        }
    } else {
        eid, err := strconv.ParseUint(entityParam, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id path param"})
            return
        }

        if !h.canAccessEntity(ctx, uint(eid)) {
            c.JSON(http.StatusForbidden, gin.H{"error": "not authorized for this entity"})
            return
        }
        entityIDs = append(entityIDs, fmt.Sprint(eid))
    }

    // 6️⃣ Build the request struct
    req := AuditLogReportRequest{
        EntityID:  entityParam,
        Action:    action,
        Status:    status,
        DateRange: dateRange,
        StartDate: start,
        EndDate:   end,
        Format:    format,
    }

    // 7️⃣ Handle JSON preview (return rows with new fields)
    if format == "" {
        data, err := h.service.GetAuditLogsReport(req, entityIDs)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        details := map[string]interface{}{
            "report_type": "audit_logs",
            "format":      "json_preview",
           // "entity_ids":  entityIDs,
            "action":      action,
            "status":      status,
            "date_range":  dateRange,
            "ip_address":  ip, // ✅ log IP for audit trail
        }
        h.auditSvc.LogAction(
            c.Request.Context(),
            &ctx.UserID,
            nil,
            "AUDIT_LOGS_REPORT_VIEWED",
            details,
            ip,
            "success",
        )
        c.JSON(http.StatusOK, data)
        return
    }

    // 8️⃣ Export file logic
    var reportType string
    switch format {
    case "excel":
        reportType = ReportTypeAuditLogsExcel
    case "pdf":
        reportType = ReportTypeAuditLogsPDF
    case "csv":
        reportType = ReportTypeAuditLogsCSV
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported export format"})
        return
    }

    bytes, fname, mime, err := h.service.ExportAuditLogsReport(
        c.Request.Context(),
        req,
        entityIDs,
        reportType,
        &ctx.UserID,
        ip,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 9️⃣ Send the file
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
    c.Data(http.StatusOK, mime, bytes)
}
