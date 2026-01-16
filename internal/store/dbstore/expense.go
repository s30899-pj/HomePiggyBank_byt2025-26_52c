package dbstore

import (
	"errors"
	"time"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type ExpenseStore struct {
	db *gorm.DB
}

type NewExpenseStoreParams struct {
	DB *gorm.DB
}

func NewExpenseStore(params NewExpenseStoreParams) *ExpenseStore {
	return &ExpenseStore{
		db: params.DB,
	}
}

func (s *ExpenseStore) CreateExpense(name string, amount float64, category store.ExpenseCategory, createdOn time.Time, householdID, createdByID uint) (uint, error) {
	expense := store.Expense{
		Name:        name,
		Amount:      amount,
		Category:    category,
		CreatedOn:   createdOn,
		HouseholdID: householdID,
		CreatedByID: createdByID,
	}

	err := s.db.Create(&expense).Error
	if err != nil {
		return 0, err
	}

	return expense.ID, nil
}

func (s *ExpenseStore) NameExists(name string) (bool, error) {
	var expense store.Expense
	err := s.db.Select("id").Where("name = ?", name).First(&expense).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, err
}

func (s *ExpenseStore) GetExpensesByHouseholdID(householdID uint) ([]store.Expense, error) {
	var expenses []store.Expense
	err := s.db.
		Where("household_id = ?", householdID).
		Preload("CreatedBy").
		Find(&expenses).Error
	return expenses, err
}
