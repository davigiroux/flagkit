package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/davigiroux/flagkit/api/internal/db"
)

type contextKey string

const ActorKey contextKey = "actor"

func Auth(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			hash := HashToken(token)

			id, err := queries.GetAPIKeyByHash(r.Context(), hash)
			if err != nil {
				http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ActorKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func GetActor(ctx context.Context) string {
	if v, ok := ctx.Value(ActorKey).(string); ok {
		return v
	}
	return ""
}
