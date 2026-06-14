package refreshtoken

import "errors"

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)
