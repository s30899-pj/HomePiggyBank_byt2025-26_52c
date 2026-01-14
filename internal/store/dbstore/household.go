package dbstore

import (
	"errors"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type HouseholdStore struct {
	db *gorm.DB
}

type NewHouseholdStoreParams struct {
	DB *gorm.DB
}

func NewHouseholdStore(params NewHouseholdStoreParams) *HouseholdStore {
	return &HouseholdStore{
		db: params.DB,
	}
}

func (s *HouseholdStore) CreateHousehold(name string, description string, createdByID uint) (uint, error) {
	household := store.Household{
		Name:        name,
		Description: description,
		CreatedByID: createdByID,
	}

	err := s.db.Create(&household).Error
	if err != nil {
		return 0, err
	}

	return household.ID, nil
}

func (s *HouseholdStore) GetHouseholdsByUserID(userID uint) ([]store.Household, error) {
	var households []store.Household
	err := s.db.Joins("JOIN memberships ON memberships.household_id = households.id").
		Where("memberships.user_id = ?", userID).
		Preload("Memberships").
		Preload("CreatedBy").
		Find(&households).Error

	if err != nil {
		return nil, err
	}

	return households, err
}

func (s *HouseholdStore) NameExists(name string) (bool, error) {
	var household store.Household
	err := s.db.Select("id").Where("name = ?", name).First(&household).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, err
}
