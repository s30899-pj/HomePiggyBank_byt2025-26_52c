package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	hashmock "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/hash/mock"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	storemock "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPostLogin_Success(t *testing.T) {
	userStore := &storemock.UserStoreMock{}
	sessionStore := &storemock.SessionStoreMock{}
	passwordHash := &hashmock.PasswordHashMock{}

	user := &store.User{
		ID:       1,
		Email:    "test@test.com",
		Password: "hashed",
	}

	userStore.On("GetUser", "test@test.com").Return(user, nil)
	passwordHash.On("ComparePasswordAndHash", "secret", "hashed").Return(true, nil)

	sessionStore.On("CreateSession", mock.Anything).Return(
		&store.Session{SessionID: "abc", UserID: 1}, nil,
	)

	handler := NewPostLoginHandler(PostLoginHandlerParams{
		UserStore:         userStore,
		SessionStore:      sessionStore,
		PasswordHash:      passwordHash,
		SessionCookieName: "session",
	})

	form := url.Values{}
	form.Set("email", "test@test.com")
	form.Set("password", "secret")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler.PostLogin(w, req)

	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "/home", resp.Header.Get("HX-Redirect"))
	require.Len(t, resp.Cookies(), 1)

	userStore.AssertExpectations(t)
	sessionStore.AssertExpectations(t)
	passwordHash.AssertExpectations(t)
}

func TestPostRegister_Success(t *testing.T) {
	userStore := &storemock.UserStoreMock{}

	userStore.On("UsernameExists", "testuser").Return(false, nil)
	userStore.On("EmailExists", "test@test.com").Return(false, nil)
	userStore.On("CreateUser", "testuser", "test@test.com", "secret").Return(nil)

	handler := NewPostRegisterHandler(PostRegisterHandlerParams{
		UserStore: userStore,
	})

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("email", "test@test.com")
	form.Set("password", "secret")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler.PostRegister(w, req)

	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "/login?from=register-success", resp.Header.Get("HX-Redirect"))

	userStore.AssertExpectations(t)
}
