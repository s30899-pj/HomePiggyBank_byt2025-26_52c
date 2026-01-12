package basic

import (
	"net/http"

	templBasic "github.com/a-h/templ"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type GetBasicHandler struct{}

func NewGetBasicHandler() *GetBasicHandler {
	return &GetBasicHandler{}
}

func (h *GetBasicHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user != nil {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	c := templ.GuestIndex()
	err := templ.Layout(c, "Welcome to Home Piggy Bank", false, nil).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *GetBasicHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	isHX := r.Header.Get("HX-Request") == "true"

	c := templ.Home(user, isHX)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Home | Home Piggy Bank", true, user)
	}

	err := out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
