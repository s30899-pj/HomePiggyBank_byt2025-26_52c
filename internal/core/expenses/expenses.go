package expenses

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	templBasic "github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
	templAlerts "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ/alerts"
)

type GetExpensesHandler struct {
	householdStore    store.HouseholdStore
	expenseShareStore store.ExpenseShareStore
}

type GetExpensesHandlerParams struct {
	HouseholdStore    store.HouseholdStore
	ExpenseShareStore store.ExpenseShareStore
}

func NewGetExpensesHandler(params GetExpensesHandlerParams) *GetExpensesHandler {
	return &GetExpensesHandler{
		householdStore:    params.HouseholdStore,
		expenseShareStore: params.ExpenseShareStore,
	}
}

func (h *GetExpensesHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	shares, err := h.expenseShareStore.GetExpensesByUserID(user.ID)
	if err != nil {
		http.Error(w, "Cannot load expenses", 500)
		return
	}

	shares = sortShares(shares)

	isHX := r.Header.Get("HX-Request") == "true"

	c := templ.Expenses(isHX, shares)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Expenses | Home Piggy Bank", true, user)
	}

	err = out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func sortShares(shares []store.ExpenseShare) []store.ExpenseShare {
	sort.SliceStable(shares, func(i, j int) bool {
		return !shares[i].Paid && shares[j].Paid
	})
	return shares
}

type GetExpensesChartHandler struct {
	expenseShareStore store.ExpenseShareStore
}

type GetExpensesChartHandlerParams struct {
	ExpenseShareStore store.ExpenseShareStore
}

func NewGetExpensesChartHandler(params GetExpensesChartHandlerParams) *GetExpensesChartHandler {
	return &GetExpensesChartHandler{
		expenseShareStore: params.ExpenseShareStore,
	}
}

func (h *GetExpensesChartHandler) GetExpensesChart(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mode := r.URL.Query().Get("mode")

	shares, err := h.expenseShareStore.GetExpensesByUserID(user.ID)
	if err != nil {
		http.Error(w, "Failed to load expenses", 500)
		return
	}

	var labels []string
	var values []float64

	var unpaidShares []store.ExpenseShare
	for _, s := range shares {
		if !s.Paid {
			unpaidShares = append(unpaidShares, s)
		}
	}

	switch mode {
	case "household":
		labels, values = sumByHousehold(unpaidShares)
	case "category":
		labels, values = sumByCategory(unpaidShares)
	case "status":
		labels, values = sumByStatus(shares)
	default:
		http.Error(w, "Invalid mode", 400)
		return
	}

	templ.ExpensesChart(labels, values).Render(r.Context(), w)
}

func sumByCategory(shares []store.ExpenseShare) ([]string, []float64) {
	m := map[string]float64{}

	for _, s := range shares {
		m[string(s.Expense.Category)] += s.Amount
	}

	var labels []string
	var values []float64
	for k, v := range m {
		labels = append(labels, k)
		values = append(values, v)
	}

	return labels, values
}

func sumByHousehold(shares []store.ExpenseShare) ([]string, []float64) {
	m := map[string]float64{}

	for _, s := range shares {
		m[s.Expense.Household.Name] += s.Amount
	}

	var labels []string
	var values []float64
	for k, v := range m {
		labels = append(labels, k)
		values = append(values, v)
	}

	return labels, values
}

func sumByStatus(shares []store.ExpenseShare) ([]string, []float64) {
	var paid, unpaid float64
	for _, s := range shares {
		if s.Paid {
			paid += s.Amount
		} else {
			unpaid += s.Amount
		}
	}
	labels := []string{"Unpaid", "Paid"}
	values := []float64{unpaid, paid}
	return labels, values
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

	const (
		maxExpenseNameLength = 40
		minExpenseAmount     = 10.0
	)

	name := r.FormValue("name")
	amountStr := r.FormValue("amount")
	categoryStr := r.FormValue("category")
	householdIDStr := r.FormValue("household_id")

	if len(name) > maxExpenseNameLength {
		w.WriteHeader(http.StatusBadRequest)
		c := templAlerts.Error(
			"Create failed",
			fmt.Sprintf("Expense name cannot be longer than %d characters", maxExpenseNameLength),
		)
		c.Render(r.Context(), w)
		return
	}

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

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templAlerts.Error("Create failed", "Invalid amount format")
		c.Render(r.Context(), w)
		return
	}

	if amount < minExpenseAmount {
		w.WriteHeader(http.StatusBadRequest)
		c := templAlerts.Error(
			"Create failed",
			fmt.Sprintf("Amount must be at least %.2f", minExpenseAmount),
		)
		c.Render(r.Context(), w)
		return
	}

	if math.Round(amount*100)/100 != amount {
		w.WriteHeader(http.StatusBadRequest)
		c := templAlerts.Error(
			"Create failed",
			"Amount can have at most 2 decimal places",
		)
		c.Render(r.Context(), w)
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

type PostExpenseShareHandler struct {
	expenseShareStore store.ExpenseShareStore
}

type PostExpenseShareHandlerParams struct {
	ExpenseShareStore store.ExpenseShareStore
}

func NewPostExpenseShareHandler(params PostExpenseShareHandlerParams) *PostExpenseShareHandler {
	return &PostExpenseShareHandler{
		expenseShareStore: params.ExpenseShareStore,
	}
}

func (h *PostExpenseShareHandler) PostPayExpenseShare(w http.ResponseWriter, r *http.Request) {
	expenseIDStr := chi.URLParam(r, "id")
	userIDStr := r.FormValue("user_id")

	expenseID, _ := strconv.Atoi(expenseIDStr)
	userID, _ := strconv.Atoi(userIDStr)

	share, err := h.expenseShareStore.GetExpenseShare(uint(expenseID), uint(userID))
	if err != nil {
		http.Error(w, "Share not found", http.StatusNotFound)
		return
	}

	share.Paid = true
	_ = h.expenseShareStore.UpdateExpenseShare(share)

	w.Header().Set("HX-Redirect", "/expenses")
	w.WriteHeader(http.StatusOK)
}
