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

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/auth"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/core/basic"
	m "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/config"
	database "github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/db"
	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store/dbstore"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	db := database.MustOpen(cfg.DatabaseName)

	//userStore := dbstore.NewUserStore(
	//	dbstore.NewUserStoreParams{
	//		DB: db,
	//	},
	//)

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
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

		r.Get("/", basic.NewBasicHandler().Index)

		r.Get("/login", auth.NewAuthHandler().GetLogin)

		r.Get("/register", auth.NewAuthHandler().GetRegister)
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
