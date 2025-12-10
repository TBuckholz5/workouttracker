package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/TBuckholz5/workouttracker/internal/util/jwt"
)

func AuthMiddleware(jwtService jwt.JwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				w.WriteHeader(http.StatusUnauthorized)
			}
			authHeader = strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
			userID, err := jwtService.ValidateJwt(authHeader)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
