package reports

import (
	"fmt"
)

// ReportService performs business logic and coordinates repo + exporter.
type ReportService interface {
	GetActivities(req ActivitiesReportRequest) (ReportData, error)
	ExportActivities(req ActivitiesReportRequest) ([]byte, string, string, error)
	GetTempleRegisteredReport(req TempleRegisteredReportRequest, entityIDs []string) ([]TempleRegisteredReportRow, error)
	ExportTempleRegisteredReport(req TempleRegisteredReportRequest, entityIDs []string, reportType string) ([]byte, string, string, error)
}

type reportService struct {
	repo     ReportRepository
	exporter ReportExporter
}

func NewReportService(repo ReportRepository, exporter ReportExporter) ReportService {
	return &reportService{repo: repo, exporter: exporter}
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

func (s *reportService) ExportActivities(req ActivitiesReportRequest) ([]byte, string, string, error) {
	// fetch data first
	data, err := s.GetActivities(req)
	if err != nil {
		return nil, "", "", err
	}

	// use exporter with format
	return s.exporter.Export(req.Type, req.Format, data)
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

func (s *reportService) ExportTempleRegisteredReport(req TempleRegisteredReportRequest, entityIDs []string, reportType string) ([]byte, string, string, error) {
    rows, err := s.GetTempleRegisteredReport(req, entityIDs)
    if err != nil {
        return nil, "", "", err
    }

    // Pass rows to ReportData.TemplesRegistered
    data := ReportData{
        TemplesRegistered: rows,
    }

    // Pass the reportType parameter to the exporter
    return s.exporter.Export(reportType, req.Format, data)
}