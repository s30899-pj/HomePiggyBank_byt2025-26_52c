package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
)

type AuthMiddleware struct {
	sessionStore      store.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(sessionStore store.SessionStore, sessionCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

type userContextKeyType struct{}

var userContextKey = userContextKeyType{}

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(m.sessionCookieName)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := base64.StdEncoding.DecodeString(sessionCookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := splitValue[0]
		userID := splitValue[1]

		user, err := m.sessionStore.GetUserFromSession(sessionID, userID)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *store.User {
	user, ok := ctx.Value(userContextKey).(*store.User)

	if !ok {
		return nil
	}

	return user
}
