package reports

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ReportExporter defines the interface for exporting reports in different formats
type ReportExporter interface {
	Export(reportType, format string, data ReportData) ([]byte, string, string, error)
}

type reportExporter struct{}

func NewReportExporter() ReportExporter {
	return &reportExporter{}
}

// Update the existing Export method in reportExporter struct
// Find this function and add the case for donations
func (e *reportExporter) Export(reportType, format string, data ReportData) ([]byte, string, string, error) {
    timestamp := time.Now().Format("20060102_150405")
    
    switch reportType {
    case ReportTypeEvents:
        return e.exportEventsByFormat(format, timestamp, data.Events)
    case ReportTypeSevas:
        return e.exportSevasByFormat(format, timestamp, data.Sevas)
    case ReportTypeBookings:
        return e.exportBookingsByFormat(format, timestamp, data.Bookings)
    case ReportTypeDonations: // Add this case
        return e.exportDonationsByFormat(format, timestamp, data.Donations)
    case ReportTypeTempleRegistered:
        return e.exportTemplesRegistered(data.TemplesRegistered)
    case ReportTypeTempleRegisteredPDF:
        return e.exportTemplesRegisteredPDF(data.TemplesRegistered)
    case ReportTypeTempleRegisteredExcel:
        return e.exportTemplesRegisteredExcel(data.TemplesRegistered)
    case ReportTypeDevoteeBirthdays:
        return e.exportDevoteeBirthdays(data.DevoteeBirthdays)
    case ReportTypeDevoteeBirthdaysPDF:
        return e.exportDevoteeBirthdaysPDF(data.DevoteeBirthdays)
    case ReportTypeDevoteeBirthdaysExcel:
        return e.exportDevoteeBirthdaysExcel(data.DevoteeBirthdays)
    default:
        return nil, "", "", fmt.Errorf("unsupported report type: %s", reportType)
    }
}

// Export Events by format
func (e *reportExporter) exportEventsByFormat(format, timestamp string, events []EventReportRow) ([]byte, string, string, error) {
	switch format {
	case FormatExcel:
		data, err := e.exportEventsExcel(events)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("events_report_%s.xlsx", timestamp)
		return data, filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		
	case FormatCSV:
		data, err := e.exportEventsCSV(events)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("events_report_%s.csv", timestamp)
		return data, filename, "text/csv", nil
		
	case FormatPDF:
		data, err := e.exportEventsPDF(events)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("events_report_%s.pdf", timestamp)
		return data, filename, "application/pdf", nil
		
	default:
		return nil, "", "", fmt.Errorf("unsupported format for events: %s", format)
	}
}

// Export Sevas by format
func (e *reportExporter) exportSevasByFormat(format, timestamp string, sevas []SevaReportRow) ([]byte, string, string, error) {
	switch format {
	case FormatExcel:
		data, err := e.exportSevasExcel(sevas)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("sevas_report_%s.xlsx", timestamp)
		return data, filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		
	case FormatCSV:
		data, err := e.exportSevasCSV(sevas)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("sevas_report_%s.csv", timestamp)
		return data, filename, "text/csv", nil
		
	case FormatPDF:
		data, err := e.exportSevasPDF(sevas)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("sevas_report_%s.pdf", timestamp)
		return data, filename, "application/pdf", nil
		
	default:
		return nil, "", "", fmt.Errorf("unsupported format for sevas: %s", format)
	}
}

// Export Bookings by format
func (e *reportExporter) exportBookingsByFormat(format, timestamp string, bookings []SevaBookingReportRow) ([]byte, string, string, error) {
	switch format {
	case FormatExcel:
		data, err := e.exportBookingsExcel(bookings)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("bookings_report_%s.xlsx", timestamp)
		return data, filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		
	case FormatCSV:
		data, err := e.exportBookingsCSV(bookings)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("bookings_report_%s.csv", timestamp)
		return data, filename, "text/csv", nil
		
	case FormatPDF:
		data, err := e.exportBookingsPDF(bookings)
		if err != nil {
			return nil, "", "", err
		}
		filename := fmt.Sprintf("bookings_report_%s.pdf", timestamp)
		return data, filename, "application/pdf", nil
		
	default:
		return nil, "", "", fmt.Errorf("unsupported format for bookings: %s", format)
	}
}

func (e *reportExporter) exportEventsExcel(events []EventReportRow) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Events"
	f.SetSheetName("Sheet1", sheetName)
	
	// Headers
	headers := []string{"Title", "Description", "Event Type", "Event Date", "Event Time", "Location", "Created By", "Created At", "Updated At", "Is Active"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	// Data
	for i, event := range events {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), event.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), event.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), event.EventType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), event.EventDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), event.EventTime)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), event.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), event.CreatedBy)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), event.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), event.UpdatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), event.IsActive)
	}
	
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *reportExporter) exportEventsCSV(events []EventReportRow) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// Headers
	headers := []string{"Title", "Description", "Event Type", "Event Date", "Event Time", "Location", "Created By", "Created At", "Updated At", "Is Active"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	// Data
	for _, event := range events {
		record := []string{
			event.Title,
			event.Description,
			event.EventType,
			event.EventDate.Format("2006-01-02"),
			event.EventTime,
			event.Location,
			strconv.FormatUint(uint64(event.CreatedBy), 10),
			event.CreatedAt.Format("2006-01-02 15:04:05"),
			event.UpdatedAt.Format("2006-01-02 15:04:05"),
			strconv.FormatBool(event.IsActive),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	
	// Important: Flush before getting bytes
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (e *reportExporter) exportEventsPDF(events []EventReportRow) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Events Report")
	pdf.Ln(20)
	
	pdf.SetFont("Arial", "B", 10)
	// Define column widths
	widths := []float64{40, 30, 25, 20, 30, 25, 15}
	headers := []string{"Title", "Event Type", "Date", "Time", "Location", "Created At", "Active"}
	
	// Print headers with borders
	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	
	// Print data rows with borders
	pdf.SetFont("Arial", "", 8)
	for _, event := range events {
		pdf.CellFormat(widths[0], 6, event.Title, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 6, event.EventType, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 6, event.EventDate.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[3], 6, event.EventTime, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 6, event.Location, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 6, event.CreatedAt.Format("01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 6, strconv.FormatBool(event.IsActive), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *reportExporter) exportSevasExcel(sevas []SevaReportRow) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sevas"
	f.SetSheetName("Sheet1", sheetName)
	
	headers := []string{"Name", "Seva Type", "Description", "Price", "Date", "Start Time", "End Time", "Duration", "Max Bookings", "Status", "Is Active", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	for i, seva := range sevas {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), seva.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), seva.SevaType)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), seva.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), seva.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), seva.Date.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), seva.StartTime)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), seva.EndTime)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), seva.Duration)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), seva.MaxBookingsPerDay)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), seva.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), seva.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), seva.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), seva.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *reportExporter) exportSevasCSV(sevas []SevaReportRow) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	headers := []string{"Name", "Seva Type", "Description", "Price", "Date", "Start Time", "End Time", "Duration", "Max Bookings", "Status", "Is Active", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	for _, seva := range sevas {
		record := []string{
			seva.Name,
			seva.SevaType,
			seva.Description,
			fmt.Sprintf("%.2f", seva.Price),
			seva.Date.Format("2006-01-02"),
			seva.StartTime,
			seva.EndTime,
			strconv.Itoa(seva.Duration),
			strconv.Itoa(seva.MaxBookingsPerDay),
			seva.Status,
			strconv.FormatBool(seva.IsActive),
			seva.CreatedAt.Format("2006-01-02 15:04:05"),
			seva.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	
	// Important: Flush before getting bytes
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (e *reportExporter) exportSevasPDF(sevas []SevaReportRow) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Sevas Report")
	pdf.Ln(20)
	
	pdf.SetFont("Arial", "B", 10)
	// Define column widths
	widths := []float64{40, 20, 20, 25, 25, 15, 20, 15}
	headers := []string{"Name", "Type", "Price", "Start Time", "End Time", "Duration", "Status", "Active"}
	
	// Print headers with borders
	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	
	// Print data rows with borders
	pdf.SetFont("Arial", "", 8)
	for _, seva := range sevas {
		pdf.CellFormat(widths[0], 6, seva.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 6, seva.SevaType, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 6, fmt.Sprintf("%.2f", seva.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[3], 6, seva.StartTime, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 6, seva.EndTime, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 6, strconv.Itoa(seva.Duration), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 6, seva.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[7], 6, strconv.FormatBool(seva.IsActive), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *reportExporter) exportBookingsExcel(bookings []SevaBookingReportRow) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Bookings"
	f.SetSheetName("Sheet1", sheetName)
	
	headers := []string{"Seva Name", "Seva Type", "Devotee Name", "Devotee Phone", "Booking Time", "Status", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	for i, booking := range bookings {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), booking.SevaName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), booking.SevaType)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), booking.DevoteeName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), booking.DevoteePhone)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), booking.BookingTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), booking.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), booking.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), booking.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *reportExporter) exportBookingsCSV(bookings []SevaBookingReportRow) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	headers := []string{"Seva Name", "Seva Type", "Devotee Name", "Devotee Phone", "Booking Time", "Status", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	for _, booking := range bookings {
		record := []string{
			booking.SevaName,
			booking.SevaType,
			booking.DevoteeName,
			booking.DevoteePhone,
			booking.BookingTime.Format("2006-01-02 15:04:05"),
			booking.Status,
			booking.CreatedAt.Format("2006-01-02 15:04:05"),
			booking.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	
	// Important: Flush before getting bytes
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (e *reportExporter) exportBookingsPDF(bookings []SevaBookingReportRow) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Seva Bookings Report")
	pdf.Ln(20)
	
	pdf.SetFont("Arial", "B", 10)
	// Define column widths
	widths := []float64{40, 25, 35, 30, 35, 20}
	headers := []string{"Seva Name", "Seva Type", "Devotee Name", "Phone", "Booking Time", "Status"}
	
	// Print headers with borders
	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	
	// Print data rows with borders
	pdf.SetFont("Arial", "", 8)
	for _, booking := range bookings {
		pdf.CellFormat(widths[0], 6, booking.SevaName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 6, booking.SevaType, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 6, booking.DevoteeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 6, booking.DevoteePhone, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 6, booking.BookingTime.Format("01-02 15:04"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 6, booking.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// exportTemplesRegistered exports temples registered as CSV.
func (e *reportExporter) exportTemplesRegistered(rows []TempleRegisteredReportRow) ([]byte, string, string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	headers := []string{"id", "name", "created_at", "status"}
	if err := w.Write(headers); err != nil {
		return nil, "", "", err
	}

	for _, r := range rows {
		record := []string{
			fmt.Sprint(r.ID),
			r.Name,
			r.CreatedAt.Format("2006-01-02 15:04:05"),
			r.Status,
		}
		if err := w.Write(record); err != nil {
			return nil, "", "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "temples_registered_report.csv", "text/csv", nil
}

// exportTemplesRegisteredExcel exports temples registered as Excel.
func (e *reportExporter) exportTemplesRegisteredExcel(rows []TempleRegisteredReportRow) ([]byte, string, string, error) {
    f := excelize.NewFile()
    sheet := "Temples Registered"
    index, err := f.NewSheet(sheet)
    if err != nil {
        return nil, "", "", err
    }
    f.DeleteSheet("Sheet1")
    f.SetActiveSheet(index)

    headers := []string{"id", "name", "created_at", "status"}
    for i, h := range headers {
        cell, err := excelize.CoordinatesToCellName(i+1, 1)
        if err != nil {
            return nil, "", "", err
        }
        f.SetCellValue(sheet, cell, h)
    }

    for rIdx, r := range rows {
        row := rIdx + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.ID)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Name)
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.CreatedAt.Format("2006-01-02 15:04:05"))
        f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.Status)
    }

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, "", "", err
    }

    return buf.Bytes(), "temples_registered_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

// exportTemplesRegisteredPDF exports temples registered as PDF.
func (e *reportExporter) exportTemplesRegisteredPDF(rows []TempleRegisteredReportRow) ([]byte, string, string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(40, 10, "Temples Registered Report")
    pdf.Ln(10)

    pdf.SetFont("Arial", "B", 10)
    headers := []string{"ID", "Name", "Created At", "Status"}
    widths := []float64{20, 80, 50, 40}

    // Print headers
    for i, h := range headers {
        pdf.CellFormat(widths[i], 7, h, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    // Print data rows
    pdf.SetFont("Arial", "", 10)
    for _, r := range rows {
        pdf.CellFormat(widths[0], 6, fmt.Sprint(r.ID), "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[1], 6, r.Name, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[2], 6, r.CreatedAt.Format("2006-01-02 15:04:05"), "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[3], 6, r.Status, "1", 0, "C", false, 0, "")
        pdf.Ln(-1)
    }

    var buf bytes.Buffer
    if err := pdf.Output(&buf); err != nil {
        return nil, "", "", err
    }

    return buf.Bytes(), "temples_registered_report.pdf", "application/pdf", nil
}

// exportDevoteeBirthdays exports devotee birthdays as CSV.
func (e *reportExporter) exportDevoteeBirthdays(rows []DevoteeBirthdayReportRow) ([]byte, string, string, error) {
    buf := &bytes.Buffer{}
    w := csv.NewWriter(buf)

    headers := []string{"full_name", "date_of_birth", "gender", "phone", "email", "temple_name", "member_since"}
    if err := w.Write(headers); err != nil {
        return nil, "", "", err
    }

    for _, r := range rows {
        record := []string{
            r.FullName,
            r.DateOfBirth.Format("2006-01-02"),
            r.Gender,
            r.Phone,
            r.Email,
            r.TempleName,
            r.MemberSince.Format("2006-01-02 15:04:05"),
        }
        if err := w.Write(record); err != nil {
            return nil, "", "", err
        }
    }

    w.Flush()
    if err := w.Error(); err != nil {
        return nil, "", "", err
    }

    return buf.Bytes(), "devotee_birthdays_report.csv", "text/csv", nil
}

// exportDevoteeBirthdaysExcel exports devotee birthdays as Excel.
func (e *reportExporter) exportDevoteeBirthdaysExcel(rows []DevoteeBirthdayReportRow) ([]byte, string, string, error) {
    f := excelize.NewFile()
    sheet := "Devotee Birthdays"
    index, err := f.NewSheet(sheet)
    if err != nil {
        return nil, "", "", err
    }
    f.DeleteSheet("Sheet1")
    f.SetActiveSheet(index)

    headers := []string{"Full Name", "Date of Birth", "Gender", "Phone", "Email", "Temple Name", "Member Since"}
    for i, h := range headers {
        cell, err := excelize.CoordinatesToCellName(i+1, 1)
        if err != nil {
            return nil, "", "", err
        }
        f.SetCellValue(sheet, cell, h)
    }

    for rIdx, r := range rows {
        row := rIdx + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.FullName)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.DateOfBirth.Format("2006-01-02"))
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.Gender)
        f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.Phone)
        f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.Email)
        f.SetCellValue(sheet, fmt.Sprintf("F%d", row), r.TempleName)
        f.SetCellValue(sheet, fmt.Sprintf("G%d", row), r.MemberSince.Format("2006-01-02 15:04:05"))
    }

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, "", "", err
    }

    return buf.Bytes(), "devotee_birthdays_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

// exportDevoteeBirthdaysPDF exports devotee birthdays as PDF.
func (e *reportExporter) exportDevoteeBirthdaysPDF(rows []DevoteeBirthdayReportRow) ([]byte, string, string, error) {
    pdf := gofpdf.New("L", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(40, 10, "Devotee Birthdays Report")
    pdf.Ln(10)

    pdf.SetFont("Arial", "B", 10)
    headers := []string{"Full Name", "Date of Birth", "Gender", "Phone", "Email", "Temple", "Member Since"}
    widths := []float64{35, 22, 12, 22, 40, 30, 22}

    // Print headers
    for i, h := range headers {
        pdf.CellFormat(widths[i], 7, h, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    // Print data rows
    pdf.SetFont("Arial", "", 8)
    for _, r := range rows {
        pdf.CellFormat(widths[0], 6, r.FullName, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[1], 6, r.DateOfBirth.Format("2006-01-02"), "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[2], 6, r.Gender, "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[3], 6, r.Phone, "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[4], 6, r.Email, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[5], 6, r.TempleName, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[6], 6, r.MemberSince.Format("2006-01-02"), "1", 0, "C", false, 0, "")
        pdf.Ln(-1)
    }

    var buf bytes.Buffer
    if err := pdf.Output(&buf); err != nil {
        return nil, "", "", err
    }

    return buf.Bytes(), "devotee_birthdays_report.pdf", "application/pdf", nil
}


// Add these functions to your existing exporter.go file
// Export Donations by format
// Add to your exporter.go file (NOT inside another function):

// Export Donations by format
func (e *reportExporter) exportDonationsByFormat(format, timestamp string, donations []DonationReportRow) ([]byte, string, string, error) {
    switch format {
    case FormatExcel:
        data, err := e.exportDonationsExcel(donations)
        if err != nil {
            return nil, "", "", err
        }
        filename := fmt.Sprintf("donations_report_%s.xlsx", timestamp)
        return data, filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
        
    case FormatCSV:
        data, err := e.exportDonationsCSV(donations)
        if err != nil {
            return nil, "", "", err
        }
        filename := fmt.Sprintf("donations_report_%s.csv", timestamp)
        return data, filename, "text/csv", nil
        
    case FormatPDF:
        data, err := e.exportDonationsPDF(donations)
        if err != nil {
            return nil, "", "", err
        }
        filename := fmt.Sprintf("donations_report_%s.pdf", timestamp)
        return data, filename, "application/pdf", nil
        
    default:
        return nil, "", "", fmt.Errorf("unsupported format for donations: %s", format)
    }
}

func (e *reportExporter) exportDonationsExcel(donations []DonationReportRow) ([]byte, error) {
    f := excelize.NewFile()
    sheetName := "Donations"
    f.SetSheetName("Sheet1", sheetName)
    
    headers := []string{"ID", "Donor Name", "Donor Email", "Amount", "Donation Type", "Payment Method", "Status", "Donation Date", "Order ID", "Payment ID", "Created At", "Updated At"}
    for i, header := range headers {
        cell := fmt.Sprintf("%c1", 'A'+i)
        f.SetCellValue(sheetName, cell, header)
    }
    
    for i, donation := range donations {
        row := i + 2
        paymentID := ""
        if donation.PaymentID != nil {
            paymentID = *donation.PaymentID
        }
        
        f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), donation.ID)
        f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), donation.DonorName)
        f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), donation.DonorEmail)
        f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), donation.Amount)
        f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), donation.DonationType)
        f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), donation.PaymentMethod)
        f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), donation.Status)
        f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), donation.DonationDate.Format("2006-01-02 15:04:05"))
        f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), donation.OrderID)
        f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), paymentID)
        f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), donation.CreatedAt.Format("2006-01-02 15:04:05"))
        f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), donation.UpdatedAt.Format("2006-01-02 15:04:05"))
    }
    
    buf, err := f.WriteToBuffer()
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func (e *reportExporter) exportDonationsCSV(donations []DonationReportRow) ([]byte, error) {
    var buf bytes.Buffer
    writer := csv.NewWriter(&buf)
    
    headers := []string{"ID", "Donor Name", "Donor Email", "Amount", "Donation Type", "Payment Method", "Status", "Donation Date", "Order ID", "Payment ID", "Created At", "Updated At"}
    if err := writer.Write(headers); err != nil {
        return nil, err
    }
    
    for _, donation := range donations {
        paymentID := ""
        if donation.PaymentID != nil {
            paymentID = *donation.PaymentID
        }
        
        record := []string{
            strconv.FormatUint(uint64(donation.ID), 10),
            donation.DonorName,
            donation.DonorEmail,
            fmt.Sprintf("%.2f", donation.Amount),
            donation.DonationType,
            donation.PaymentMethod,
            donation.Status,
            donation.DonationDate.Format("2006-01-02 15:04:05"),
            donation.OrderID,
            paymentID,
            donation.CreatedAt.Format("2006-01-02 15:04:05"),
            donation.UpdatedAt.Format("2006-01-02 15:04:05"),
        }
        if err := writer.Write(record); err != nil {
            return nil, err
        }
    }
    
    // Important: Flush before getting bytes
    writer.Flush()
    if err := writer.Error(); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

func (e *reportExporter) exportDonationsPDF(donations []DonationReportRow) ([]byte, error) {
    pdf := gofpdf.New("L", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(0, 10, "Donations Report")
    pdf.Ln(20)
    
    pdf.SetFont("Arial", "B", 10)
    // Define column widths
    widths := []float64{35, 30, 20, 25, 25, 20, 25, 25}
    headers := []string{"Donor Name", "Donor Email", "Amount", "Type", "Method", "Status", "Donation Date", "Order ID"}
    
    // Print headers with borders
    for i, header := range headers {
        pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)
    
    // Print data rows with borders
    pdf.SetFont("Arial", "", 8)
    for _, donation := range donations {
        pdf.CellFormat(widths[0], 6, donation.DonorName, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[1], 6, donation.DonorEmail, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[2], 6, fmt.Sprintf("%.2f", donation.Amount), "1", 0, "R", false, 0, "")
        pdf.CellFormat(widths[3], 6, donation.DonationType, "1", 0, "L", false, 0, "")
        pdf.CellFormat(widths[4], 6, donation.PaymentMethod, "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[5], 6, donation.Status, "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[6], 6, donation.DonationDate.Format("2006-01-02"), "1", 0, "C", false, 0, "")
        pdf.CellFormat(widths[7], 6, donation.OrderID, "1", 0, "L", false, 0, "")
        pdf.Ln(-1)
    }
    
    var buf bytes.Buffer
    if err := pdf.Output(&buf); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}