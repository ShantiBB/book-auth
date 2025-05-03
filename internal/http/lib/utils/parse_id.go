package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"auth/internal/http/lib/schema/response"
)

func ParseID(w http.ResponseWriter, r *http.Request, strID string) (int64, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid id"))
		return 0, err
	}

	return int64(id), nil
}
