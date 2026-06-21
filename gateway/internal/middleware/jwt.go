package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cmiski/sawako/gateway/internal/jwt"
)

const userIDHeader = "X-User-ID"

func JWTAuth(
	validator *jwt.Validator,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			token, err := bearerToken(r)
			if err != nil {
				writeUnauthorized(w, err.Error())
				return
			}

			userID, err := validator.Validate(token)
			if err != nil {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			r.Header.Set(
				userIDHeader,
				userID.String(),
			)

			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return "", errMissingAuthorization
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", errInvalidAuthorization
	}

	token := strings.TrimSpace(
		header[len(prefix):],
	)

	if token == "" {
		return "", errInvalidAuthorization
	}

	return token, nil
}

func writeUnauthorized(
	w http.ResponseWriter,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
