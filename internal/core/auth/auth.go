package auth

import (
	"net/http"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	c := templ.Login()
	err := templ.Layout(c, "Log in").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	c := templ.Register()
	err := templ.Layout(c, "Sign up").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
