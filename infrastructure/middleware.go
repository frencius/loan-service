package infrastructure

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: auth checking

		// return hardcoded user id
		ctx := context.WithValue(r.Context(), "userID", "f2c86f5c-6578-4d63-aa01-5bd4246c3bd8")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
