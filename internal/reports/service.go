package reports

import (
	"context"
	"fmt"

	"github.com/sharath018/temple-management-backend/internal/auditlog"
)

// ReportService performs business logic and coordinates repo + exporter.
type ReportService interface {
	GetActivities(req ActivitiesReportRequest) (ReportData, error)
	ExportActivities(ctx context.Context, req ActivitiesReportRequest, userID *uint, ip string) ([]byte, string, string, error)
	GetTempleRegisteredReport(req TempleRegisteredReportRequest, entityIDs []string) ([]TempleRegisteredReportRow, error)
	ExportTempleRegisteredReport(ctx context.Context, req TempleRegisteredReportRequest, entityIDs []string, reportType string, userID *uint, ip string) ([]byte, string, string, error)
	GetDevoteeBirthdaysReport(req DevoteeBirthdaysReportRequest, entityIDs []string) ([]DevoteeBirthdayReportRow, error)
	ExportDevoteeBirthdaysReport(ctx context.Context, req DevoteeBirthdaysReportRequest, entityIDs []string, reportType string, userID *uint, ip string) ([]byte, string, string, error)
}

type reportService struct {
	repo     ReportRepository
	exporter ReportExporter
	auditSvc auditlog.Service
}

func NewReportService(repo ReportRepository, exporter ReportExporter, auditSvc auditlog.Service) ReportService {
	return &reportService{
		repo:     repo,
		exporter: exporter,
		auditSvc: auditSvc,
	}
}

func (s *reportService) GetActivities(req ActivitiesReportRequest) (ReportData, error) {
	// validate request
	if req.Type != ReportTypeEvents && req.Type != ReportTypeSevas && req.Type != ReportTypeBookings {
		return ReportData{}, fmt.Errorf("invalid report type: %s", req.Type)
	}
	start := req.StartDate
	end := req.EndDate

	var data ReportData
	var err error
	switch req.Type {
	case ReportTypeEvents:
		data.Events, err = s.repo.GetEvents(convertUintSlice(req.EntityIDs), start, end)
	case ReportTypeSevas:
		data.Sevas, err = s.repo.GetSevas(convertUintSlice(req.EntityIDs), start, end)
	case ReportTypeBookings:
		data.Bookings, err = s.repo.GetSevaBookings(convertUintSlice(req.EntityIDs), start, end)
	}
	return data, err
}

func (s *reportService) ExportActivities(ctx context.Context, req ActivitiesReportRequest, userID *uint, ip string) ([]byte, string, string, error) {
	// fetch data first
	data, err := s.GetActivities(req)
	if err != nil {
		// Log failed export attempt
		details := map[string]interface{}{
			"report_type": req.Type,
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_ACTIVITIES_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// use exporter with format
	bytes, filename, mimeType, err := s.exporter.Export(req.Type, req.Format, data)
	if err != nil {
		// Log failed export
		details := map[string]interface{}{
			"report_type": req.Type,
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_ACTIVITIES_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// Log successful export
	details := map[string]interface{}{
		"report_type": req.Type,
		"format":     req.Format,
		"filename":   filename,
		"entity_ids": req.EntityIDs,
		"date_range": req.DateRange,
	}
	s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_ACTIVITIES_REPORT_DOWNLOADED", details, ip, "success")

	return bytes, filename, mimeType, nil
}

// convert slice of string ids to []uint
func convertUintSlice(strs []string) []uint {
	out := make([]uint, 0, len(strs))
	for _, s := range strs {
		var id uint
		_, err := fmt.Sscan(s, &id)
		if err == nil {
			out = append(out, id)
		}
	}
	return out
}

func (s *reportService) GetTempleRegisteredReport(req TempleRegisteredReportRequest, entityIDs []string) ([]TempleRegisteredReportRow, error) {
	status := req.Status
	start := req.StartDate
	end := req.EndDate
	
	return s.repo.GetTemplesRegistered(convertUintSlice(entityIDs), start, end, status)
}

func (s *reportService) ExportTempleRegisteredReport(ctx context.Context, req TempleRegisteredReportRequest, entityIDs []string, reportType string, userID *uint, ip string) ([]byte, string, string, error) {
	rows, err := s.GetTempleRegisteredReport(req, entityIDs)
	if err != nil {
		// Log failed export attempt
		details := map[string]interface{}{
			"report_type": "temple_registered",
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_REGISTER_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// Pass rows to ReportData.TemplesRegistered
	data := ReportData{
		TemplesRegistered: rows,
	}

	// Pass the reportType parameter to the exporter
	bytes, filename, mimeType, err := s.exporter.Export(reportType, req.Format, data)
	if err != nil {
		// Log failed export
		details := map[string]interface{}{
			"report_type": "temple_registered",
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_REGISTER_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// Log successful export
	details := map[string]interface{}{
		"report_type": "temple_registered",
		"format":     req.Format,
		"filename":   filename,
		"entity_ids": entityIDs,
		"status":     req.Status,
		"date_range": req.DateRange,
		"record_count": len(rows),
	}
	s.auditSvc.LogAction(ctx, userID, nil, "TEMPLE_REGISTER_REPORT_DOWNLOADED", details, ip, "success")

	return bytes, filename, mimeType, nil
}

func (s *reportService) GetDevoteeBirthdaysReport(req DevoteeBirthdaysReportRequest, entityIDs []string) ([]DevoteeBirthdayReportRow, error) {
	start := req.StartDate
	end := req.EndDate
	
	return s.repo.GetDevoteeBirthdays(convertUintSlice(entityIDs), start, end)
}

func (s *reportService) ExportDevoteeBirthdaysReport(ctx context.Context, req DevoteeBirthdaysReportRequest, entityIDs []string, reportType string, userID *uint, ip string) ([]byte, string, string, error) {
	rows, err := s.GetDevoteeBirthdaysReport(req, entityIDs)
	if err != nil {
		// Log failed export attempt
		details := map[string]interface{}{
			"report_type": "devotee_birthdays",
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "DEVOTEE_BIRTHDAYS_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// Pass rows to ReportData.DevoteeBirthdays
	data := ReportData{
		DevoteeBirthdays: rows,
	}

	// Pass the reportType parameter to the exporter
	bytes, filename, mimeType, err := s.exporter.Export(reportType, req.Format, data)
	if err != nil {
		// Log failed export
		details := map[string]interface{}{
			"report_type": "devotee_birthdays",
			"format":     req.Format,
			"error":      err.Error(),
		}
		s.auditSvc.LogAction(ctx, userID, nil, "DEVOTEE_BIRTHDAYS_REPORT_DOWNLOAD_FAILED", details, ip, "failure")
		return nil, "", "", err
	}

	// Log successful export
	details := map[string]interface{}{
		"report_type": "devotee_birthdays",
		"format":     req.Format,
		"filename":   filename,
		"entity_ids": entityIDs,
		"date_range": req.DateRange,
		"record_count": len(rows),
	}
	s.auditSvc.LogAction(ctx, userID, nil, "DEVOTEE_BIRTHDAYS_REPORT_DOWNLOADED", details, ip, "success")

	return bytes, filename, mimeType, nil
}