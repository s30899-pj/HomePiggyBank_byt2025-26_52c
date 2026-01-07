package dbstore

import (
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/hash"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordHash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordHash: params.PasswordHash,
	}
}

func (s *UserStore) CreateUser(username string, email string, password string) error {
	hashedPassword, err := s.passwordHash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&store.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
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

func (s *UserStore) CheckUsername(username string) (bool, error) {
	var count int64
	err := s.db.Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, err
}
