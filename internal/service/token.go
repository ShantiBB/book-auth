package service

import (
	"context"

	"auth/internal/entity"
	"auth/internal/http/lib/jwt"
	"auth/package/utils"
)

type TokenRepository interface {
	GetUserCredentialsByUsername(ctx context.Context, u *entity.User) error
	CreateUser(ctx context.Context, u *entity.User) error
}

func (s *Service) Register(ctx context.Context, u *entity.User) (*entity.Token, error) {
	const op = "user.service.RegisterUser"

	var err error
	u.PasswordHash, err = utils.HashPassword(u.PasswordHash)
	if err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return nil, err
	}

	if err = s.repo.CreateUser(ctx, u); err != nil {
		s.log.Error("failed to create user", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("user create success", "op", op, "id", u.ID)

	var tokens *entity.Token
	tokens, err = jwt.GenerateAllTokens(u.ID, u.Role)
	if err != nil {
		s.log.Error("failed to generate tokens", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", u.ID)
	return tokens, nil
}

func (s *Service) Login(ctx context.Context, u *entity.User) (*entity.Token, error) {
	const op = "user.service.LoginToken"

	var password = u.PasswordHash
	err := s.repo.GetUserCredentialsByUsername(ctx, u)
	if err != nil {
		s.log.Error("failed to get user by id", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("user credentials success", "op", op, "id", u.ID)

	err = utils.CheckPasswordHash(u.PasswordHash, password)
	if err != nil {
		s.log.Error("failed to check password", "op", op, "error", err)
		return nil, err
	}

	var tokens *entity.Token
	tokens, err = jwt.GenerateAllTokens(u.ID, u.Role)
	if err != nil {
		s.log.Error("failed to generate tokens", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", u.ID)
	return tokens, nil
}

func (s *Service) Refresh(token string) (string, error) {
	const op = "user.service.RefreshToken"

	claims, err := jwt.GetClaimsRefreshToken(token)
	if err != nil {
		s.log.Error("failed to parse refresh token", "op", op, "error", err)
		return "", err
	}

	s.log.Debug("refresh token success", "op", op, "id", claims.Sub)

	var accessToken string
	accessToken, err = jwt.GenerateAccessToken(claims.Sub, claims.Role)
	if err != nil {
		s.log.Error("failed to generate access token", "op", op, "error", err)
		return "", err
	}

	s.log.Debug("success", "op", op, "id", claims.Sub)
	return accessToken, nil
}
