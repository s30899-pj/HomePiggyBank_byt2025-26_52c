package reports

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
)

func GenerateReportPDF(report store.Report) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(0, 10, "Expense Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(0, 8, fmt.Sprintf("Period: %s - %s",
		report.PeriodStart.Format("02.01.2006"),
		report.PeriodEnd.Format("02.01.2006"),
	))
	pdf.Ln(8)

	pdf.Cell(0, 8, fmt.Sprintf("Total expenses: %.2f", report.TotalExpenses))
	pdf.Ln(8)

	pdf.Cell(0, 8, fmt.Sprintf("Payment status: %s", report.PaymentStatus))
	pdf.Ln(8)

	dir := "./files/reports"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, report.FileName)

	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
