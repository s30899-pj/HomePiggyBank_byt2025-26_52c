package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/hash"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
	templAlerts "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ/alerts"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	c := templ.Register()
	err := templ.Layout(c, "Sign up").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// TODO: add usage of alert from register page
// TODO: add alerts for wrong email address or password
func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	c := templ.Login()
	err := templ.Layout(c, "Log in").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type PostRegisterHandler struct {
	userStore store.UserStore
}

type PostRegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewPostRegisterHandler(params PostRegisterHandlerParams) *PostRegisterHandler {
	return &PostRegisterHandler{
		userStore: params.UserStore,
	}
}

// TODO: add verification for existing email address or username
// TODO: add alerts for an existing email address or username
func (h *PostRegisterHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(username, email, password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templAlerts.Error("Registration failed", "There was a problem creating your account. Please check your details and try again.")
		c.Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

type PostLoginHandler struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	passwordHash      hash.PasswordHash
	sessionCookieName string
}

type PostLoginHandlerParams struct {
	UserStore         store.UserStore
	SessionStore      store.SessionStore
	PasswordHash      hash.PasswordHash
	SessionCookieName string
}

func NewPostLoginHandler(params PostLoginHandlerParams) *PostLoginHandler {
	return &PostLoginHandler{
		userStore:         params.UserStore,
		sessionStore:      params.SessionStore,
		passwordHash:      params.PasswordHash,
		sessionCookieName: params.SessionCookieName,
	}
}

// TODO: Add an expiration time setting to login and registration alerts and fix the close button
func (h *PostLoginHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	remember := r.FormValue("remember")

	user, err := h.userStore.GetUser(email)

	//TODO: change alert for not existing user
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c := templAlerts.Error("Login failed", "Unable to sign in. Please try again.")
		c.Render(r.Context(), w)
		return
	}

	passwordValid, err := h.passwordHash.ComparePasswordAndHash(password, user.Password)

	// TODO: add alert description
	if err != nil || !passwordValid {
		w.WriteHeader(http.StatusUnauthorized)
		c := templAlerts.Error("Login failed", "")
		c.Render(r.Context(), w)
		return
	}

	session, err := h.sessionStore.CreateSession(&store.Session{
		UserID: user.ID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := user.ID
	sessionID := session.SessionID

	cookieValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", sessionID, userID)))

	var expiration time.Time
	if remember == "yes" {
		expiration = time.Now().Add(30 * 24 * time.Hour)
	}

	cookie := http.Cookie{
		Name:     h.sessionCookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if !expiration.IsZero() {
		cookie.Expires = expiration
	}

	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

type PostLogoutHandler struct {
	sessionCookieName string
}

type PostLogoutHandlerParams struct {
	SessionCookieName string
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLogoutHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now(),
		Path:    "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}
