package households

import (
	"log"
	"net/http"
	"strconv"

	templBasic "github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ"
	templAlerts "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/templ/alerts"
)

type GetHouseholdsHandler struct {
	householdStore store.HouseholdStore
	userStore      store.UserStore
}

type GetHouseholdsHandlerParams struct {
	HouseholdStore store.HouseholdStore
	UserStore      store.UserStore
}

func NewGetHouseholdsHandler(params GetHouseholdsHandlerParams) *GetHouseholdsHandler {
	return &GetHouseholdsHandler{
		householdStore: params.HouseholdStore,
		userStore:      params.UserStore,
	}
}

func (h *GetHouseholdsHandler) GetHouseholds(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	isHX := r.Header.Get("HX-Request") == "true"

	households, err := h.householdStore.GetHouseholdsByUserID(user.ID)
	if err != nil {
		http.Error(w, "cannot fetch households", http.StatusInternalServerError)
		return
	}

	ownedHouseholds := make([]store.Household, 0, len(households))
	for _, h := range households {
		if h.CreatedByID == user.ID {
			ownedHouseholds = append(ownedHouseholds, h)
		}
	}

	allUsers, err := h.userStore.GetAllUsers()
	filteredUsers := make([]store.User, 0, len(allUsers))
	for _, u := range allUsers {
		if u.ID != user.ID {
			filteredUsers = append(filteredUsers, u)
		}
	}

	if err != nil {
		http.Error(w, "cannot fetch users", http.StatusInternalServerError)
		return
	}

	c := templ.Households(isHX, households, ownedHouseholds, filteredUsers)

	var out templBasic.Component
	if isHX {
		out = c
	} else {
		out = templ.Layout(c, "Households | Home Piggy Bank", true, user)
	}

	err = out.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type PostHouseholdHandler struct {
	householdStore  store.HouseholdStore
	membershipStore store.MembershipStore
	userStore       store.UserStore
}

type PostHouseholdHandlerParams struct {
	HouseholdStore  store.HouseholdStore
	MembershipStore store.MembershipStore
	UserStore       store.UserStore
}

func NewPostHouseholdHandler(params PostHouseholdHandlerParams) *PostHouseholdHandler {
	return &PostHouseholdHandler{
		householdStore:  params.HouseholdStore,
		membershipStore: params.MembershipStore,
		userStore:       params.UserStore,
	}
}

func (h *PostHouseholdHandler) PostHousehold(w http.ResponseWriter, r *http.Request) {
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
	description := r.FormValue("description")
	memberUsernames := r.Form["members[]"]

	nameBusy, err := h.householdStore.NameExists(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if nameBusy {
		w.WriteHeader(http.StatusConflict)
		c := templAlerts.Error("Create failed", "Household with this name already exists")
		c.Render(r.Context(), w)
		return
	}

	householdID, err := h.createHouseholdWithMembership(name, description, user.ID, "owner")
	if err != nil {
		http.Error(w, "could not create household", http.StatusInternalServerError)
		return
	}

	if err := h.addMembers(memberUsernames, householdID); err != nil {
		http.Error(w, "could not add members", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/households")
	w.WriteHeader(http.StatusOK)
}

func (h *PostHouseholdHandler) createHouseholdWithMembership(householdName string, description string, userID uint, role string) (uint, error) {

	householdID, err := h.householdStore.CreateHousehold(householdName, description, userID)
	if err != nil {
		return 0, err
	}

	if err := h.membershipStore.CreateMembership(userID, householdID, role); err != nil {
		return 0, err
	}

	return householdID, nil
}

func (h *PostHouseholdHandler) addMembers(usernames []string, householdID uint) error {
	for _, username := range usernames {
		user, err := h.userStore.GetUserByUsername(username)
		if err != nil {
			continue
		}

		if err := h.membershipStore.CreateMembership(user.ID, householdID, "member"); err != nil {
			log.Printf("failed to add member %s: %v", username, err)
			continue
		}
	}
	return nil
}

type GetHouseholdMembersHandler struct {
	membershipStore store.MembershipStore
}

type GetHouseholdMembersHandlerParams struct {
	MembershipStore store.MembershipStore
}

func NewGetHouseholdMembersHandler(params GetHouseholdMembersHandlerParams) *GetHouseholdMembersHandler {
	return &GetHouseholdMembersHandler{
		membershipStore: params.MembershipStore,
	}
}

func (h *GetHouseholdMembersHandler) GetHouseholdMembers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	householdIDStr := chi.URLParam(r, "id")
	householdID, err := strconv.ParseUint(householdIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid household id", http.StatusBadRequest)
		return
	}

	members, err := h.membershipStore.GetMembersByHouseholdID(uint(householdID))
	if err != nil {
		http.Error(w, "cannot fetch members", http.StatusInternalServerError)
		return
	}

	err = templ.HouseholdMembers(members).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type GetHouseholdExpensesHandler struct {
	expenseStore store.ExpenseStore
}

type GetHouseholdExpensesHandlerParams struct {
	ExpenseStore store.ExpenseStore
}

func NewGetHouseholdExpensesHandler(params GetHouseholdExpensesHandlerParams) *GetHouseholdExpensesHandler {
	return &GetHouseholdExpensesHandler{
		expenseStore: params.ExpenseStore,
	}
}

func (h *GetHouseholdExpensesHandler) GetHouseholdExpenses(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	householdIDStr := chi.URLParam(r, "id")
	householdID, err := strconv.ParseUint(householdIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid household id", http.StatusBadRequest)
		return
	}

	expenses, err := h.expenseStore.GetExpensesByHouseholdID(uint(householdID))
	if err != nil {
		http.Error(w, "cannot fetch members", http.StatusInternalServerError)
		return
	}

	if expenses == nil {
		expenses = []store.Expense{}
	}

	err = templ.HouseholdExpenses(expenses).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
