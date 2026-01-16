package reports

import (
	"net/http"
	"path/filepath"
	"time"

	templBasic "github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type GetReportsHandler struct {
	reportStore store.ReportStore
}

type GetReportsHandlerParams struct {
	ReportStore store.ReportStore
}

func NewGetReportsHandler(params GetReportsHandlerParams) *GetReportsHandler {
	return &GetReportsHandler{
		reportStore: params.ReportStore,
	}
}

func (h *GetReportsHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	reports, err := h.reportStore.GetReportsByUser(user.ID)
	if err != nil {
		http.Error(w, "Failed to load reports", http.StatusInternalServerError)
		return
	}

	isHX := r.Header.Get("HX-Request") == "true"

	c := templ.Reports(isHX, reports)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Reports | Home Piggy Bank", true, user)
	}

	err = out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type GetReportHandler struct {
	reportStore store.ReportStore
}

type GetReportHandlerParams struct {
	ReportStore store.ReportStore
}

func NewGetReportHandler(params GetReportHandlerParams) *GetReportHandler {
	return &GetReportHandler{
		reportStore: params.ReportStore,
	}
}

func (h *GetReportHandler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fileName := chi.URLParam(r, "file")

	report, err := h.reportStore.GetReportByFileName(fileName)
	if err != nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	if report.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	path := filepath.Join("./files/reports", fileName)
	http.ServeFile(w, r, path)
}

type PostReportHandler struct {
	reportStore store.ReportStore
}

type PostReportHandlerParams struct {
	ReportStore store.ReportStore
}

func NewPostReportsHandler(params PostReportHandlerParams) *PostReportHandler {
	return &PostReportHandler{
		reportStore: params.ReportStore,
	}
}

func (h *PostReportHandler) PostGenerateReport(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	fromStr := r.FormValue("period_start")
	toStr := r.FormValue("period_end")
	paymentStatus := r.FormValue("payment_status")

	switch paymentStatus {
	case "all", "paid", "unpaid":
	default:
		paymentStatus = "all"
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	if from.After(to) {
		http.Error(w, "Start date must be before end date", http.StatusBadRequest)
		return
	}

	to = to.AddDate(0, 0, 1).Add(-time.Nanosecond)

	report, err := h.reportStore.CreateReport(user.ID, from, to, paymentStatus)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	_, err = GenerateReportPDF(report)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/reports")
	w.WriteHeader(http.StatusOK)
}
