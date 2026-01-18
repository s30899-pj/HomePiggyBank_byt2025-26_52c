package mock

import (
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) CreateUser(username, email, password string) error {
	args := m.Called(username, email, password)
	return args.Error(0)
}

func (m *UserStoreMock) GetUser(email string) (*store.User, error) {
	args := m.Called(email)
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *UserStoreMock) GetUserByUsername(username string) (*store.User, error) {
	args := m.Called(username)
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *UserStoreMock) GetAllUsers() ([]store.User, error) {
	args := m.Called()
	return args.Get(0).([]store.User), args.Error(1)
}

func (m *UserStoreMock) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *UserStoreMock) UsernameExists(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

type SessionStoreMock struct {
	mock.Mock
}

func (m *SessionStoreMock) CreateSession(session *store.Session) (*store.Session, error) {
	args := m.Called(session)
	return args.Get(0).(*store.Session), args.Error(1)
}

func (m *SessionStoreMock) GetUserFromSession(sessionID string, userID string) (*store.User, error) {
	args := m.Called(sessionID, userID)
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *SessionStoreMock) DeleteSession(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
