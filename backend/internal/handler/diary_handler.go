package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type DiaryHandler struct {
	svc *service.DiaryService
}

func NewDiaryHandler(svc *service.DiaryService) *DiaryHandler {
	return &DiaryHandler{svc: svc}
}

func (h *DiaryHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	var req domain.DiaryEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.svc.CreateOrUpdate(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, entry)
}

func (h *DiaryHandler) GetByDate(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	date := r.URL.Query().Get("date")
	if date == "" {
		httputil.Error(w, http.StatusBadRequest, "date parameter required")
		return
	}

	entry, err := h.svc.GetByDate(r.Context(), userID, date)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entry == nil {
		httputil.JSON(w, http.StatusOK, nil)
		return
	}
	httputil.JSON(w, http.StatusOK, entry)
}

func (h *DiaryHandler) GetMonthly(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	month := r.URL.Query().Get("month")
	if month == "" {
		httputil.Error(w, http.StatusBadRequest, "month parameter required (YYYY-MM)")
		return
	}

	summary, err := h.svc.GetMonthly(r.Context(), userID, month)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, summary)
}

func (h *DiaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.NoContent(w)
}
