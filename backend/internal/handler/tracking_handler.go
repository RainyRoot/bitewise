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

type TrackingHandler struct {
	tracking *service.TrackingService
}

func NewTrackingHandler(tracking *service.TrackingService) *TrackingHandler {
	return &TrackingHandler{tracking: tracking}
}

func (h *TrackingHandler) LogFood(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.FoodLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log, err := h.tracking.LogFood(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, log)
}

func (h *TrackingHandler) GetFoodLogs(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	date := r.URL.Query().Get("date")
	logs, err := h.tracking.GetFoodLogs(r.Context(), userID, date)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, logs)
}

func (h *TrackingHandler) DeleteFoodLog(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.tracking.DeleteFoodLog(r.Context(), id, userID); err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.NoContent(w)
}

func (h *TrackingHandler) LogWater(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.WaterLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log, err := h.tracking.LogWater(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, log)
}

func (h *TrackingHandler) GetWaterLogs(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	date := r.URL.Query().Get("date")
	logs, err := h.tracking.GetWaterLogs(r.Context(), userID, date)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, logs)
}

func (h *TrackingHandler) GetNutritionSummary(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	date := r.URL.Query().Get("date")
	summary, err := h.tracking.GetNutritionSummary(r.Context(), userID, date)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, summary)
}
