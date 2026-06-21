package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cmiski/sawako/services/auth/internal/credential"
	"github.com/cmiski/sawako/services/auth/internal/refreshtoken"
	"github.com/cmiski/sawako/services/auth/internal/user"
)

type RegisterRequest struct {
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequest struct {
	RefreshToken string
}

type Service struct {
	users                 UserService
	credentials           CredentialRepository
	refreshTokens         RefreshTokenRepository
	refreshTokenService   RefreshTokenService
	hasher                PasswordHasher
	issuer                TokenIssuer
	refreshTokenGenerator RefreshTokenGenerator
	tokenHasher           TokenHasher
	txManager             TransactionManager
}

func NewService(
	users UserService,
	credentials CredentialRepository,
	refreshTokens RefreshTokenRepository,
	refreshTokenService RefreshTokenService,
	hasher PasswordHasher,
	issuer TokenIssuer,
	refreshTokenGenerator RefreshTokenGenerator,
	tokenHasher TokenHasher,
	txManager TransactionManager,
) *Service {
	return &Service{
		users:                 users,
		credentials:           credentials,
		refreshTokens:         refreshTokens,
		refreshTokenService:   refreshTokenService,
		hasher:                hasher,
		issuer:                issuer,
		refreshTokenGenerator: refreshTokenGenerator,
		tokenHasher:           tokenHasher,
		txManager:             txManager,
	}
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) error {
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	passwordHash, err := s.hasher.Hash(
		req.Password,
	)
	if err != nil {
		return fmt.Errorf(
			"register user: hash password: %w",
			err,
		)
	}

	err = s.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context) error {
			u := &user.User{
				Email: email,
			}

			if err := s.users.Create(
				ctx,
				u,
			); err != nil {
				return fmt.Errorf(
					"create user: %w",
					err,
				)
			}

			c := &credential.Credential{
				UserID:         u.ID,
				CredentialType: credential.CredentialTypePassword,
				PasswordHash:   passwordHash,
			}

			if err := s.credentials.Create(
				ctx,
				c,
			); err != nil {
				return fmt.Errorf(
					"create credential: %w",
					err,
				)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf(
			"register user: %w",
			err,
		)
	}

	return nil
}

const refreshTokenTTL = 30 * 24 * time.Hour

func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
) (*LoginResponse, error) {
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	u, err := s.users.GetByEmail(
		ctx,
		email,
	)

	switch {
	case err == nil:
		// continue

	case errors.Is(err, user.ErrUserNotFound):
		return nil, fmt.Errorf(
			"login user: %w",
			ErrInvalidCredentials,
		)

	default:
		return nil, fmt.Errorf(
			"login user: get user: %w",
			err,
		)
	}

	cred, err := s.credentials.GetByUserIDAndType(
		ctx,
		u.ID,
		credential.CredentialTypePassword,
	)

	switch {
	case err == nil:
		// continue

	case errors.Is(err, credential.ErrCredentialNotFound):
		return nil, fmt.Errorf(
			"login user: %w",
			ErrInvalidCredentials,
		)

	default:
		return nil, fmt.Errorf(
			"login user: get credential: %w",
			err,
		)
	}

	if err := s.hasher.Verify(
		req.Password,
		cred.PasswordHash,
	); err != nil {
		return nil, fmt.Errorf(
			"login user: %w",
			ErrInvalidCredentials,
		)
	}

	claims := Claims{
		UserID: u.ID,
	}

	accessToken, err := s.issuer.IssueAccessToken(
		claims,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"login user: issue access token: %w",
			err,
		)
	}

	refreshToken, err := s.refreshTokenGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf(
			"login user: generate refresh token: %w",
			err,
		)
	}

	tokenHash := s.tokenHasher.Hash(
		refreshToken,
	)

	expiresAt := time.Now().
		Add(refreshTokenTTL)

	rt := s.refreshTokenService.New(
		u.ID,
		tokenHash,
		expiresAt,
	)

	if err := s.refreshTokens.Create(
		ctx,
		rt,
	); err != nil {
		return nil, fmt.Errorf(
			"login user: create refresh token: %w",
			err,
		)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	req RefreshRequest,
) (*LoginResponse, error) {
	refreshToken := strings.TrimSpace(
		req.RefreshToken,
	)

	if refreshToken == "" {
		return nil, fmt.Errorf(
			"refresh token: %w",
			ErrInvalidRefreshToken,
		)
	}

	tokenHash := s.tokenHasher.Hash(
		refreshToken,
	)

	storedToken, err := s.refreshTokens.GetByHash(
		ctx,
		tokenHash,
	)

	switch {
	case err == nil:
		// continue

	case errors.Is(err, refreshtoken.ErrRefreshTokenNotFound):
		return nil, fmt.Errorf(
			"refresh token: %w",
			ErrInvalidRefreshToken,
		)

	default:
		return nil, fmt.Errorf(
			"refresh token: get: %w",
			err,
		)
	}

	if storedToken.RevokedAt != nil {
		return nil, fmt.Errorf(
			"refresh token: %w",
			ErrInvalidRefreshToken,
		)
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, fmt.Errorf(
			"refresh token: %w",
			ErrInvalidRefreshToken,
		)
	}

	var (
		accessToken     string
		newRefreshToken string
	)

	err = s.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context) error {
			if err := s.refreshTokens.Revoke(
				ctx,
				storedToken.ID,
			); err != nil {
				return fmt.Errorf(
					"revoke refresh token: %w",
					err,
				)
			}

			claims := Claims{
				UserID: storedToken.UserID,
			}

			issuedAccessToken, err := s.issuer.IssueAccessToken(
				claims,
			)
			if err != nil {
				return fmt.Errorf(
					"issue access token: %w",
					err,
				)
			}

			issuedRefreshToken, err := s.refreshTokenGenerator.Generate()
			if err != nil {
				return fmt.Errorf(
					"generate refresh token: %w",
					err,
				)
			}

			newHash := s.tokenHasher.Hash(
				issuedRefreshToken,
			)

			expiresAt := time.Now().
				Add(refreshTokenTTL)

			newToken := s.refreshTokenService.New(
				storedToken.UserID,
				newHash,
				expiresAt,
			)

			if err := s.refreshTokens.Create(
				ctx,
				newToken,
			); err != nil {
				return fmt.Errorf(
					"create refresh token: %w",
					err,
				)
			}

			accessToken = issuedAccessToken
			newRefreshToken = issuedRefreshToken

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"refresh token: %w",
			err,
		)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(
	ctx context.Context,
	req RefreshRequest,
) error {
	refreshToken := strings.TrimSpace(
		req.RefreshToken,
	)

	if refreshToken == "" {
		return fmt.Errorf(
			"logout: %w",
			ErrInvalidRefreshToken,
		)
	}

	tokenHash := s.tokenHasher.Hash(
		refreshToken,
	)

	storedToken, err := s.refreshTokens.GetByHash(
		ctx,
		tokenHash,
	)

	switch {
	case errors.Is(err, refreshtoken.ErrRefreshTokenNotFound):
		return nil

	case err == nil:
		// continue

	default:
		return fmt.Errorf(
			"logout: get refresh token: %w",
			err,
		)
	}

	if storedToken.RevokedAt != nil {
		return nil
	}

	if err := s.refreshTokens.Revoke(
		ctx,
		storedToken.ID,
	); err != nil {
		if errors.Is(err, refreshtoken.ErrRefreshTokenNotFound) {
			return nil
		}

		return fmt.Errorf(
			"logout: revoke refresh token: %w",
			err,
		)
	}

	return nil
}
