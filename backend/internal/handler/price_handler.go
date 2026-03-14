package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type PriceHandler struct {
	svc *service.PriceService
}

func NewPriceHandler(svc *service.PriceService) *PriceHandler {
	return &PriceHandler{svc: svc}
}

func (h *PriceHandler) LogPrice(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	var req domain.PriceLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log, err := h.svc.LogPrice(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.Created(w, log)
}

func (h *PriceHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	logs, err := h.svc.GetLogs(r.Context(), userID, limit)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, logs)
}

func (h *PriceHandler) GetTrend(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	item := r.URL.Query().Get("item")
	if item == "" {
		httputil.Error(w, http.StatusBadRequest, "item parameter required")
		return
	}

	trend, err := h.svc.GetTrend(r.Context(), userID, item)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, trend)
}

func (h *PriceHandler) CompareStores(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	item := r.URL.Query().Get("item")
	if item == "" {
		httputil.Error(w, http.StatusBadRequest, "item parameter required")
		return
	}

	comparison, err := h.svc.CompareStores(r.Context(), userID, item)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, comparison)
}

func (h *PriceHandler) GetSpending(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	month := r.URL.Query().Get("month")

	summary, err := h.svc.GetSpendingSummary(r.Context(), userID, month)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, summary)
}
