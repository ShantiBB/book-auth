package service

import (
	"context"

	"auth/internal/entity"
	"auth/package/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *entity.User) error
	GetUserByID(ctx context.Context, u *entity.User) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUserByID(ctx context.Context, u *entity.User) error
	DeleteUserByID(ctx context.Context, id int64) error
}

func (s *Service) CreateUser(ctx context.Context, u *entity.User) error {
	const op = "user.service.Create"

	var err error
	u.PasswordHash, err = utils.HashPassword(u.PasswordHash)
	if err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	if err = s.repo.CreateUser(ctx, u); err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", u.ID)

	return nil
}

func (s *Service) GetUserByID(ctx context.Context, u *entity.User) error {
	const op = "user.service.GetByID"

	err := s.repo.GetUserByID(ctx, u)
	if err != nil {
		s.log.Error("failed to get book by id", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", u.ID)

	return nil
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	const op = "user.service.GetAll"

	books, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		s.log.Error("failed to get all books", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "count", len(books))
	return books, nil
}

func (s *Service) UpdateUserByID(ctx context.Context, u *entity.User) error {
	const op = "user.service.Update"

	err := s.repo.UpdateUserByID(ctx, u)
	if err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", u.ID)

	return nil
}

func (s *Service) DeleteUserByID(ctx context.Context, id int64) error {
	const op = "book.service.DeleteByID"

	if err := s.repo.DeleteUserByID(ctx, id); err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", id)

	return nil
}
