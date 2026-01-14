package dbstore

import (
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type MembershipStore struct {
	db *gorm.DB
}

type NewMembershipStoreParams struct {
	DB *gorm.DB
}

func NewMembershipStore(params NewMembershipStoreParams) *MembershipStore {
	return &MembershipStore{
		db: params.DB,
	}
}

func (s *MembershipStore) CreateMembership(userID uint, householdID uint, role string) error {
	return s.db.Create(&store.Membership{
		UserID:      userID,
		HouseholdID: householdID,
		Role:        role,
	}).Error
}
