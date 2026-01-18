package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	storemock "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/mock"
	"github.com/stretchr/testify/require"
)

func TestAddUserToContext_NoCookie(t *testing.T) {
	sessionStore := &storemock.SessionStoreMock{}

	middleware := NewAuthMiddleware(sessionStore, "session")

	handler := middleware.AddUserToContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		require.Nil(t, user)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestAddUserToContext_InvalidSession(t *testing.T) {
	sessionStore := &storemock.SessionStoreMock{}
	sessionStore.
		On("GetUserFromSession", "invalid", "1").
		Return((*store.User)(nil), http.ErrNoCookie) // dowolny error

	middleware := NewAuthMiddleware(sessionStore, "session")

	handler := middleware.AddUserToContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		require.Nil(t, user)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	cookieValue := base64.StdEncoding.EncodeToString([]byte("invalid:1"))
	req.AddCookie(&http.Cookie{Name: "session", Value: cookieValue})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	sessionStore.AssertExpectations(t)
}

func TestAddUserToContext_Success(t *testing.T) {
	sessionStore := &storemock.SessionStoreMock{}

	expectedUser := &store.User{ID: 1, Username: "test"}

	sessionStore.
		On("GetUserFromSession", "session-id", "1").
		Return(expectedUser, nil)

	middleware := NewAuthMiddleware(sessionStore, "session")

	handler := middleware.AddUserToContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		require.NotNil(t, user)
		require.Equal(t, expectedUser.ID, user.ID)
		require.Equal(t, expectedUser.Username, user.Username)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	cookieValue := base64.StdEncoding.EncodeToString([]byte("session-id:1"))
	req.AddCookie(&http.Cookie{Name: "session", Value: cookieValue})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	sessionStore.AssertExpectations(t)
}
