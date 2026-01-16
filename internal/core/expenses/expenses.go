package expenses

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	templBasic "github.com/a-h/templ"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
	templAlerts "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ/alerts"
)

type GetExpensesHandler struct {
	householdStore store.HouseholdStore
}

type GetExpensesHandlerParams struct {
	HouseholdStore store.HouseholdStore
}

func NewGetExpensesHandler(params GetExpensesHandlerParams) *GetExpensesHandler {
	return &GetExpensesHandler{
		householdStore: params.HouseholdStore,
	}
}

func (h *GetExpensesHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	isHX := r.Header.Get("HX-Request") == "true"

	c := templ.Expenses(isHX)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Expenses | Home Piggy Bank", true, user)
	}

	err := out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type PostExpenseHandler struct {
	expenseStore      store.ExpenseStore
	expenseShareStore store.ExpenseShareStore
	membershipStore   store.MembershipStore
	userStore         store.UserStore
}

type PostExpenseHandlerParams struct {
	ExpenseStore      store.ExpenseStore
	ExpenseShareStore store.ExpenseShareStore
	MembershipStore   store.MembershipStore
	UserStore         store.UserStore
}

func NewPostExpenseHandler(params PostExpenseHandlerParams) *PostExpenseHandler {
	return &PostExpenseHandler{
		expenseStore:      params.ExpenseStore,
		expenseShareStore: params.ExpenseShareStore,
		membershipStore:   params.MembershipStore,
		userStore:         params.UserStore,
	}
}

func (h *PostExpenseHandler) PostExpense(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	amountStr := r.FormValue("amount")
	categoryStr := r.FormValue("category")
	householdIDStr := r.FormValue("household_id")

	nameBusy, err := h.expenseStore.NameExists(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if nameBusy {
		w.WriteHeader(http.StatusConflict)
		c := templAlerts.Error("Create failed", "Expense with this name already exists")
		c.Render(r.Context(), w)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 32)
	if err != nil || amount <= 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	householdID64, err := strconv.ParseUint(householdIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid household", http.StatusBadRequest)
		return
	}
	householdID := uint(householdID64)

	category := store.ExpenseCategory(categoryStr)
	if !category.IsValid() {
		http.Error(w, "invalid category", http.StatusBadRequest)
		return
	}

	expenseID, err := h.expenseStore.CreateExpense(
		name,
		amount,
		category,
		time.Now(),
		householdID,
		user.ID,
	)
	if err != nil {
		http.Error(w, "cannot create expense", http.StatusInternalServerError)
		return
	}

	members, err := h.membershipStore.GetMembersByHouseholdID(householdID)
	if err != nil || len(members) == 0 {
		http.Error(w, "cannot fetch household members", http.StatusInternalServerError)
		return
	}

	shares := splitAmount(amount, len(members))

	for i, member := range members {
		if err := h.expenseShareStore.CreateExpenseShare(expenseID, member.UserID, shares[i]); err != nil {
			log.Printf("cannot create expense share for user %d: %v", member.UserID, err)
		}
	}

	w.Header().Set("HX-Redirect", "/households")
	w.WriteHeader(http.StatusOK)
}

func splitAmount(amount float64, membersCount int) []float64 {
	if membersCount == 0 {
		return nil
	}

	perPerson := math.Ceil(amount/float64(membersCount)*100) / 100
	shares := make([]float64, membersCount)
	for i := range shares {
		shares[i] = perPerson
	}

	return shares
}
