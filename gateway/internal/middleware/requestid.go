package middleware

import (
	"net/http"

	"github.com/cmiski/sawako/shared/contextx"
	"github.com/cmiski/sawako/shared/uuidx"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" || !uuidx.IsValid(requestID) {
			requestID = uuidx.NewV7().String()
		}

		ctx := contextx.SetRequestID(
			r.Context(),
			requestID,
		)

		r = r.WithContext(ctx)

		w.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(w, r)
	})

}
