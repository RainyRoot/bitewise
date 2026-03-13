package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type ShoppingHandler struct {
	shopping *service.ShoppingService
}

func NewShoppingHandler(shopping *service.ShoppingService) *ShoppingHandler {
	return &ShoppingHandler{shopping: shopping}
}

func (h *ShoppingHandler) GenerateFromMealPlan(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		MealPlanID int64 `json:"meal_plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.MealPlanID == 0 {
		httputil.Error(w, http.StatusBadRequest, "meal_plan_id required")
		return
	}

	list, err := h.shopping.GenerateFromMealPlan(r.Context(), userID, req.MealPlanID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, list)
}

func (h *ShoppingHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.shopping.GetCurrent(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}

func (h *ShoppingHandler) ToggleItem(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.shopping.ToggleItem(r.Context(), id, userID); err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.NoContent(w)
}
