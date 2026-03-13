package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rainyroot/bitewise/backend/internal/nutrition"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type NutritionHandler struct {
	provider nutrition.NutritionProvider
}

func NewNutritionHandler(provider nutrition.NutritionProvider) *NutritionHandler {
	return &NutritionHandler{provider: provider}
}

func (h *NutritionHandler) LookupBarcode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		httputil.Error(w, http.StatusBadRequest, "barcode required")
		return
	}

	item, err := h.provider.LookupBarcode(r.Context(), code)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, item)
}

func (h *NutritionHandler) SearchFood(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httputil.Error(w, http.StatusBadRequest, "query parameter 'q' required")
		return
	}

	items, err := h.provider.SearchFood(r.Context(), query)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, items)
}

func (h *NutritionHandler) GetSeasonal(w http.ResponseWriter, r *http.Request) {
	result := nutrition.GetSeasonalProduce()
	httputil.JSON(w, http.StatusOK, result)
}
