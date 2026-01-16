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
	}).Error
}
