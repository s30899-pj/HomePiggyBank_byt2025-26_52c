package basic

import (
	"net/http"

	templBasic "github.com/a-h/templ"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type BasicHandler struct{}

func NewBasicHandler() *BasicHandler {
	return &BasicHandler{}
}

func (h *BasicHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	isLoggedIn := user != nil

	var c, l templBasic.Component
	if isLoggedIn {
		c = templ.Index(user)
		l = templ.Layout(c, "Home | Home Piggy Bank", isLoggedIn, user)
	} else {
		c = templ.GuestIndex()
		l = templ.Layout(c, "Welcome to Home Piggy Bank", isLoggedIn, nil)
	}

	err := l.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
