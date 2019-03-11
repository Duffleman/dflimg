package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UsernameKey contextKey = "username"

func AuthMiddleware(users map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authKey := r.Header.Get("Authorization")

			if authKey == "" {
				http.Error(w, "access_denied", 500)
				return
			}

			for username, key := range users {
				if key == authKey {
					ctx = context.WithValue(ctx, UsernameKey, username)
					next.ServeHTTP(w, r.WithContext(ctx))

					return
				}
			}

			http.Error(w, "access_denied", 500)
			return

		}

		return http.HandlerFunc(fn)
	}
}
