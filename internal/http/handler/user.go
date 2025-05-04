package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"auth/internal/entity"
	"auth/internal/http/lib/schema/request"
	"auth/internal/http/lib/schema/response"
	"auth/internal/http/lib/utils"
	"auth/internal/http/lib/validate"
	"auth/internal/repository/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	CreateUser(ctx context.Context, u *entity.User) error
	GetUserByID(ctx context.Context, u *entity.User) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUserByID(ctx context.Context, u *entity.User) error
	DeleteUserByID(ctx context.Context, id int64) error
}

// CreateUser @Summary      Create new user
// @Description  Register a new user and return basic info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      request.UserCreate  true  "New user"
// @Success      201   {object}  response.UserShort
// @Failure      400   {object}  response.Response
// @Failure      409   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /users [post]
func (h *Handler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.UserCreate

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		v := validator.New()
		_ = v.RegisterValidation("passwd", validate.Password)
		if err := v.Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, validate.Error(validateErr))

			return
		}

		user := &entity.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: req.Password,
		}

		err := h.svc.CreateUser(r.Context(), user)
		if errors.Is(err, postgres.DuplicateError) {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Error("username or email already exists"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, response.UserShort{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get all user info by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.User
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [get]
func (h *Handler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		user := &entity.User{ID: id}
		err = h.svc.GetUserByID(r.Context(), user)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

// GetUserAll godoc
// @Summary      Get all users
// @Description  Get short users info
// @Tags         users
// @Produce      json
// @Success      200  {array}   response.UserShort
// @Failure      500  {object}  response.Response
// @Router       /users [get]
func (h *Handler) GetUserAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.svc.GetAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get users"))
			return
		}

		var userList []response.UserShort
		for _, user := range users {
			userList = append(userList, response.UserShort{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				Role:     user.Role,
			})
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, userList)
	}
}

// UpdateUserByID godoc
// @Summary      Update         user by ID
// @Description  Updates        user data based on the provided ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int                 true  "User ID"
// @Param        user body           request.UserUpdate  true  "User update request"
// @Success      200  {object}  response.UserShort
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [put]
func (h *Handler) UpdateUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		var req request.UserUpdate
		if err = render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, validate.Error(validateErr))
			return
		}

		user := &entity.User{
			ID:       id,
			Username: req.Username,
			Email:    req.Email,
		}

		err = h.svc.UpdateUserByID(r.Context(), user)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.UserShort{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}
}

// DeleteUserByID godoc
// @Summary      Delete user by ID
// @Description  Deletes a user by their unique ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		if err = h.svc.DeleteUserByID(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("user not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete user"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
