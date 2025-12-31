package basic

import (
	"net/http"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
)

type BasicHandler struct{}

func NewBasicHandler() *BasicHandler {
	return &BasicHandler{}
}

func (h *BasicHandler) Index(w http.ResponseWriter, r *http.Request) {
	c := templ.Index()
	err := templ.Layout(c, "Welcome in Home Piggy Bank").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
