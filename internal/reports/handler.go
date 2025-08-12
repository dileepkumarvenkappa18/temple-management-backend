package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// Handler holds service & repo (repo used for entity lookups here)
type Handler struct {
	service ReportService
	repo    ReportRepository
}

// NewHandler creates a new reports handler
func NewHandler(svc ReportService, repo ReportRepository) *Handler {
	return &Handler{service: svc, repo: repo}
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
	// Removed req.TenantID assignment, since not used in ActivitiesReportRequest anymore

	// If no format -> return JSON preview
	if req.Format == "" {
		data, err := h.service.GetActivities(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}

	// Else export file (format present)
	bytes, fname, mime, err := h.service.ExportActivities(req)
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
    
    // --- START OF UPDATED LOGIC ---
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
        c.JSON(http.StatusOK, data)
        return
    }

    // Export file (format is present)
    // The service method must be updated to accept the `reportType` parameter
    bytes, fname, mime, err := h.service.ExportTempleRegisteredReport(req, entityIDs, reportType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
    c.Data(http.StatusOK, mime, bytes)
    // --- END OF UPDATED LOGIC ---
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
        c.JSON(http.StatusOK, data)
        return
    }

    // Export file (format is present)
    bytes, fname, mime, err := h.service.ExportDevoteeBirthdaysReport(req, entityIDs, reportType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
    c.Data(http.StatusOK, mime, bytes)
}