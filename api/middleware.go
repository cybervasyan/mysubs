package api

import (
	"context"
	"mysub/models"
	"net/http"
)

type contextKey string

const subsKey contextKey = "user"

func SetUserCtx(subs []models.Subscription) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), subsKey, subs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
