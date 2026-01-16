package dbstore

import (
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type ExpenseShareStore struct {
	db *gorm.DB
}

type NewExpenseShareStoreParams struct {
	DB *gorm.DB
}

func NewExpenseShareStore(params NewExpenseShareStoreParams) *ExpenseShareStore {
	return &ExpenseShareStore{
		db: params.DB,
	}
}

func (s *ExpenseShareStore) CreateExpenseShare(expenseID uint, userID uint, amount float64) error {
	return s.db.Create(&store.ExpenseShare{
		ExpenseID: expenseID,
		UserID:    userID,
		Amount:    amount,
		Paid:      false,
	}).Error
}

func (s *ExpenseShareStore) GetExpenseShare(expenseID uint, userID uint) (store.ExpenseShare, error) {
	var share store.ExpenseShare
	err := s.db.
		Preload("Expense").
		Preload("Expense.Household").
		Preload("User").
		Where("expense_id = ? AND user_id = ?", expenseID, userID).
		First(&share).Error
	return share, err
}

func (s *ExpenseShareStore) GetExpensesByUserID(userID uint) ([]store.ExpenseShare, error) {
	var shares []store.ExpenseShare

	err := s.db.
		Preload("Expense").
		Preload("Expense.Household").
		Preload("User").
		Where("user_id = ?", userID).
		Order("expense_id desc").
		Find(&shares).Error

	return shares, err
}

func (s *ExpenseShareStore) UpdateExpenseShare(share store.ExpenseShare) error {
	return s.db.Save(&share).Error
}
