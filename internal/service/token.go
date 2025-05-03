package service

import (
	"context"

	"auth/internal/entity"
	"auth/internal/http/lib/jwt"
	"auth/package/utils"
)

type TokenRepository interface {
	GetUserCredentialsByUsername(ctx context.Context, u *entity.User) error
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

	token := &entity.Token{}
	token.AccessToken, err = jwt.GenerateAccessToken(u.ID, u.Role)
	if err != nil {
		s.log.Error("failed to generate access token", "op", op, "error", err)
		return nil, err
	}

	token.RefreshToken, err = jwt.GenerateRefreshToken(u.ID, u.Role)
	if err != nil {
		s.log.Error("failed to generate refresh token", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", u.ID)

	return token, nil
}

func (s *Service) Refresh(token string) (*entity.Token, error) {
	const op = "user.service.RefreshToken"

	claims, err := jwt.ParseRefreshToken(token)
	if err != nil {
		s.log.Error("failed to parse refresh token", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("refresh token success", "op", op, "id", claims.Sub)

	accessToken := &entity.Token{}
	accessToken.AccessToken, err = jwt.GenerateAccessToken(claims.Sub, claims.Role)
	if err != nil {
		s.log.Error("failed to generate access token", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", claims.Sub)

	return accessToken, nil
}
