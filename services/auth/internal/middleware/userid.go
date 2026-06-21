package middleware

import (
	"net/http"

	"github.com/cmiski/sawako/services/auth/internal/handlers"
	"github.com/cmiski/sawako/shared/contextx"
	"github.com/google/uuid"
)

func RequireUserID(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		userIDValue := r.Header.Get(
			contextx.UserIDHeader,
		)

		if userIDValue == "" {
			handlers.WriteUnauthorized(
				w,
				"missing user identity",
			)
			return
		}

		userID, err := uuid.Parse(userIDValue)
		if err != nil {
			handlers.WriteUnauthorized(
				w,
				"invalid user identity",
			)
			return
		}

		ctx := contextx.SetUserID(
			r.Context(),
			userID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
