package auth

import (
	"net/http"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
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

	err := h.userStore.CreateUser(username, email, password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	}

}
