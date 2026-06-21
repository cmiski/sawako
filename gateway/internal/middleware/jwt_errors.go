package middleware

import "errors"

var (
	errMissingAuthorization = errors.New(
		"missing authorization header",
	)

	errInvalidAuthorization = errors.New(
		"invalid authorization header",
	)
)
