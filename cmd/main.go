package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/config"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/auth"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/basic"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/expenses"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/households"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/reports"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/hash/passwordhash"
	m "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"
	database "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/db"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/dbstore"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	db := database.MustOpen(cfg.DatabaseName)

	passwordhash := passwordhash.NewPasswordHash()

	userStore := dbstore.NewUserStore(
		dbstore.NewUserStoreParams{
			DB:           db,
			PasswordHash: passwordhash,
		},
	)

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		},
	)

	householdStore := dbstore.NewHouseholdStore(
		dbstore.NewHouseholdStoreParams{
			DB: db,
		},
	)

	membershipStore := dbstore.NewMembershipStore(
		dbstore.NewMembershipStoreParams{
			DB: db,
		},
	)

	expenseStore := dbstore.NewExpenseStore(
		dbstore.NewExpenseStoreParams{
			DB: db,
		},
	)

	expenseShareStore := dbstore.NewExpenseShareStore(
		dbstore.NewExpenseShareStoreParams{
			DB: db,
		},
	)

	reportStore := dbstore.NewReportStore(
		dbstore.NewReportStoreParams{
			DB: db,
		},
	)

	fileServer := http.FileServer(http.Dir("./web/static"))

	r.Get("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
	}))

	authMiddleware := m.NewAuthMiddleware(sessionStore, cfg.SessionCookieName)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			authMiddleware.AddUserToContext,
		)

		//BASIC
		r.Get("/", basic.NewGetBasicHandler().GetIndex)

		r.Get("/home", basic.NewGetBasicHandler().GetHome)

		//AUTH
		r.Get("/register", auth.NewGetAuthHandler().GetRegister)

		r.Post("/register", auth.NewPostRegisterHandler(auth.PostRegisterHandlerParams{
			UserStore: userStore,
		}).PostRegister)

		r.Get("/login", auth.NewGetAuthHandler().GetLogin)

		r.Post("/login", auth.NewPostLoginHandler(auth.PostLoginHandlerParams{
			UserStore:         userStore,
			SessionStore:      sessionStore,
			PasswordHash:      passwordhash,
			SessionCookieName: cfg.SessionCookieName,
		}).PostLogin)

		r.Post("/logout", auth.NewPostLogoutHandler(auth.PostLogoutHandlerParams{
			SessionStore:      sessionStore,
			SessionCookieName: cfg.SessionCookieName,
		}).PostLogout)

		//HOUSEHOLDS
		r.Get("/households", households.NewGetHouseholdsHandler(households.GetHouseholdsHandlerParams{
			HouseholdStore: householdStore,
			UserStore:      userStore,
		}).GetHouseholds)

		r.Get("/household/{id}/members", households.NewGetHouseholdMembersHandler(households.GetHouseholdMembersHandlerParams{
			MembershipStore: membershipStore,
		}).GetHouseholdMembers)

		r.Get("/household/{id}/expenses", households.NewGetHouseholdExpensesHandler(households.GetHouseholdExpensesHandlerParams{
			ExpenseStore: expenseStore,
		}).GetHouseholdExpenses)

		r.Post("/household", households.NewPostHouseholdHandler(households.PostHouseholdHandlerParams{
			HouseholdStore:  householdStore,
			MembershipStore: membershipStore,
			UserStore:       userStore,
		}).PostHousehold)

		//EXPENSES
		r.Get("/expenses", expenses.NewGetExpensesHandler(expenses.GetExpensesHandlerParams{
			HouseholdStore:    householdStore,
			ExpenseShareStore: expenseShareStore,
		}).GetExpenses)

		r.Get("/expenses/chart", expenses.NewGetExpensesChartHandler(expenses.GetExpensesChartHandlerParams{
			ExpenseShareStore: expenseShareStore,
		}).GetExpensesChart)

		r.Post("/expense", expenses.NewPostExpenseHandler(expenses.PostExpenseHandlerParams{
			ExpenseStore:      expenseStore,
			ExpenseShareStore: expenseShareStore,
			MembershipStore:   membershipStore,
			UserStore:         userStore,
		}).PostExpense)

		r.Post("/expense/{id}/pay", expenses.NewPostExpenseShareHandler(expenses.PostExpenseShareHandlerParams{
			ExpenseShareStore: expenseShareStore,
		}).PostPayExpenseShare)

		//REPORTS
		r.Get("/reports", reports.NewGetReportsHandler(reports.GetReportsHandlerParams{
			ReportStore: reportStore,
		}).GetReports)

		r.Get("/reports/files/{file}", reports.NewGetReportHandler(reports.GetReportHandlerParams{
			ReportStore: reportStore,
		}).DownloadPDF)

		r.Post("/report", reports.NewPostReportsHandler(reports.PostReportHandlerParams{
			ReportStore: reportStore,
		}).PostGenerateReport)
	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port))
	<-killSig

	logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shoutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
