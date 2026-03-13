package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type MealPlanHandler struct {
	plans *service.MealPlanService
}

func NewMealPlanHandler(plans *service.MealPlanService) *MealPlanHandler {
	return &MealPlanHandler{plans: plans}
}

func (h *MealPlanHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.MealPlanGenerateRequest
	if r.ContentLength > 0 {
		json.NewDecoder(r.Body).Decode(&req)
	}

	plan, err := h.plans.Generate(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, plan)
}

func (h *MealPlanHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plan, err := h.plans.GetCurrent(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid plan id")
		return
	}

	plan, err := h.plans.GetByID(r.Context(), userID, planID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid plan id")
		return
	}

	entryID, err := strconv.ParseInt(chi.URLParam(r, "entryId"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	var req domain.MealPlanEntryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.plans.UpdateEntry(r.Context(), userID, planID, entryID, req); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}

func (h *MealPlanHandler) Activate(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid plan id")
		return
	}

	if err := h.plans.Activate(r.Context(), userID, planID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}
