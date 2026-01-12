package households

import (
	"net/http"

	templBasic "github.com/a-h/templ"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type GetHouseholdsHandler struct{}

func NewGetHouseholdsHandler() *GetHouseholdsHandler { return &GetHouseholdsHandler{} }

func (h *GetHouseholdsHandler) GetHouseholds(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	isHX := r.Header.Get("HX-Request") == "true"

	c := templ.Households(isHX)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Households | Home Piggy Bank", true, user)
	}

	err := out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
