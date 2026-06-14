package credential

import "errors"

var (
	ErrCredentialNotFound = errors.New(
		"credential not found",
	)
)
