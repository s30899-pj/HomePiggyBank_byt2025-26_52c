package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	templBasic "github.com/a-h/templ"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/hash"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
	templAlerts "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ/alerts"
)

type GetAuthHandler struct{}

func NewGetAuthHandler() *GetAuthHandler {
	return &GetAuthHandler{}
}

func (h *GetAuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	c := templ.Register()
	err := templ.Layout(c, "Sign up | Home Piggy Bank", false, nil).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *GetAuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")

	var alert templBasic.Component
	if from == "register-success" {
		alert = templAlerts.Success(
			"Registration successful",
			"Your account has been created. You can now log in.",
		)
	}

	c := templ.Login(alert)
	err := templ.Layout(c, "Log in | Home Piggy Bank", false, nil).Render(r.Context(), w)

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

func (h *PostRegisterHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	registerError := func(code int, title string, description string) {
		w.WriteHeader(code)
		c := templAlerts.Error(title, description)
		c.Render(r.Context(), w)
	}

	usernameBusy, err := h.userStore.UsernameExists(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if usernameBusy {
		registerError(http.StatusConflict, "Registration failed", "User with this username already exists")
		return
	}

	emailBusy, err := h.userStore.EmailExists(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if emailBusy {
		registerError(http.StatusConflict, "Registration failed", "User with this email already exists")
		return
	}

	err = h.userStore.CreateUser(username, email, password)

	if err != nil {
		registerError(http.StatusBadRequest, "Registration failed", "There was a problem creating your account. Please check your details and try again.")
		return
	}

	w.Header().Set("HX-Redirect", "/login?from=register-success")
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

func (h *PostLoginHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	remember := r.FormValue("remember")

	user, err := h.userStore.GetUser(email)

	passwordValid := false
	if err == nil {
		passwordValid, err = h.passwordHash.ComparePasswordAndHash(password, user.Password)
	}

	if err != nil || !passwordValid {
		w.WriteHeader(http.StatusUnauthorized)
		c := templAlerts.Error("Login failed", "Invalid email or password.")
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
	sessionStore      store.SessionStore
	sessionCookieName string
}

type PostLogoutHandlerParams struct {
	SessionStore      store.SessionStore
	SessionCookieName string
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLogoutHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := user.ID

	err := h.sessionStore.DeleteSession(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now(),
		Path:    "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}
