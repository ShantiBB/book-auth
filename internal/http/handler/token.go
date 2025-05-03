package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"auth/internal/entity"
	"auth/internal/http/lib/schema/request"
	"auth/internal/http/lib/schema/response"
	"auth/internal/http/lib/validate"
)

type TokenService interface {
	Login(ctx context.Context, u *entity.User) (*entity.Token, error)
	Refresh(token string) (*entity.Token, error)
}

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
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("failed to login"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, token)
	}
}

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

		token, err := h.svc.Refresh(req.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("failed to refresh token"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.AccessToken{
			AccessToken: token.AccessToken,
		})
	}
}
