package dbstore

import (
	"fmt"
	"time"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type ReportStore struct {
	db *gorm.DB
}

type NewReportStoreParams struct {
	DB *gorm.DB
}

func NewReportStore(params NewReportStoreParams) *ReportStore {
	return &ReportStore{
		db: params.DB,
	}
}

func (s *ReportStore) CreateReport(userID uint, from, to time.Time, paymentStatus string) (store.Report, error) {
	var total float64

	query := s.db.Model(&store.ExpenseShare{}).
		Select("SUM(expense_shares.amount)").
		Joins("JOIN expenses ON expenses.id = expense_shares.expense_id").
		Where(
			"expense_shares.user_id = ? AND expenses.created_on BETWEEN ? AND ?",
			userID, from, to,
		)

	switch paymentStatus {
	case "paid":
		query = query.Where("expense_shares.paid = ?", true)
	case "unpaid":
		query = query.Where("expense_shares.paid = ?", false)
	}

	if err := query.Scan(&total).Error; err != nil {
		return store.Report{}, err
	}

	fileName := fmt.Sprintf("report_%d_%d.pdf", userID, time.Now().Unix())

	report := store.Report{
		UserID:         userID,
		PeriodStart:    from,
		PeriodEnd:      to,
		TotalExpenses:  total,
		PaymentStatus:  paymentStatus,
		GenerationDate: time.Now(),
		FileName:       fileName,
	}

	if err := s.db.Create(&report).Error; err != nil {
		return store.Report{}, err
	}

	return report, nil
}

func (s *ReportStore) GetReportsByUser(userID uint) ([]store.Report, error) {
	var reports []store.Report
	err := s.db.Where("user_id = ?", userID).Order("generation_date desc").Find(&reports).Error
	return reports, err
}

func (s *ReportStore) GetReportByFileName(fileName string) (store.Report, error) {
	var report store.Report
	err := s.db.Where("file_name = ?", fileName).First(&report).Error
	return report, err
}
