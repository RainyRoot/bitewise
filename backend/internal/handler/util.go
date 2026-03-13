package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseIDParam(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	if idStr == "" {
		return 0, fmt.Errorf("missing param %s", param)
	}
	return strconv.ParseInt(idStr, 10, 64)
}
