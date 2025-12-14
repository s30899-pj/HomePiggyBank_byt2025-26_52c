package dbstore

import (
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

type NewUserStoreParams struct {
	DB *gorm.DB
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db: params.DB,
	}
}

func (s *UserStore) CreateUser(username string, email string, password string) error {
	return s.db.Create(&store.User{
		Username: username,
		Email:    email,
		Password: password,
	}).Error
}

func (s *UserStore) GetUser(email string) (*store.User, error) {
	var user store.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}
