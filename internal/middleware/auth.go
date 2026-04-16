package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/samaita/quick-go/internal/auth"
)

type SessionStore interface {
	Get(ctx context.Context, jti string) (string, error)
}

// JWT validates the token cookie on every request.
// If Redis is unreachable, it falls back to JWT-only validation (logged as warning).
// Returns 401 JSON on missing or invalid token.
func JWT(jwtSecret string, sessions SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeUnauthorized(w, "missing token")
				return
			}

			claims, err := auth.Parse(jwtSecret, cookie.Value)
			if err != nil {
				writeUnauthorized(w, "invalid token")
				return
			}

			// Redis session check (revocation support)
			stored, err := sessions.Get(r.Context(), claims.ID)
			if err != nil {
				// Redis unavailable: fall back to JWT-only validation
				log.Printf("middleware: WARNING: redis unavailable, skipping session check: %v", err)
			} else if stored == "" {
				// Key not found in Redis → session revoked or expired
				writeUnauthorized(w, "session expired")
				return
			}

			ctx := context.WithValue(r.Context(), auth.ClaimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
