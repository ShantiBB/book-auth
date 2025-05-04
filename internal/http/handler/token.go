package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"auth/internal/entity"
	"auth/internal/http/lib/schema/request"
	"auth/internal/http/lib/schema/response"
	"auth/internal/http/lib/validate"
	"auth/internal/repository/postgres"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type TokenService interface {
	Register(ctx context.Context, u *entity.User) (*entity.Token, error)
	Login(ctx context.Context, u *entity.User) (*entity.Token, error)
	Refresh(token string) (string, error)
}

// Register godoc
// @Summary      Register new user
// @Description  Creates a new user account and returns auth token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      request.Register  true  "Registration data"
// @Success      200   {object}  response.Tokens
// @Failure      400   {object}  response.Response
// @Failure      409   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /auth/register [post]
func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.Register

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

		var user = &entity.User{
			Username:     req.Username,
			PasswordHash: req.Password,
			Email:        req.Email,
		}

		token, err := h.svc.Register(r.Context(), user)
		if errors.Is(err, postgres.DuplicateError) {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Error("username or email already exists"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to register"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Tokens{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		})
	}
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns auth token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      request.Login  true  "Login credentials"
// @Success      200   {object}  response.Tokens
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /auth/login [post]
func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.Login

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, validate.Error(validateErr))
		}

		var user = &entity.User{
			Username:     req.Username,
			PasswordHash: req.Password,
		}

		token, err := h.svc.Login(r.Context(), user)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to login"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Tokens{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		})
	}
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Returns a new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  body      request.Refresh  true  "Refresh token request"
// @Success      200    {object}  response.AccessToken
// @Failure      400    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Router       /auth/refresh [post]
func (h *Handler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.Refresh

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, validate.Error(validateErr))
			return
		}

		accessToken, err := h.svc.Refresh(req.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to refresh token"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.AccessToken{
			AccessToken: accessToken,
		})
	}
}
