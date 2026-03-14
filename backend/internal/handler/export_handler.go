package handler

import (
	"net/http"

	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type ExportHandler struct {
	svc *service.ExportService
}

func NewExportHandler(svc *service.ExportService) *ExportHandler {
	return &ExportHandler{svc: svc}
}

func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=bitewise_food_logs.csv")

	if err := h.svc.ExportFoodLogsCSV(r.Context(), userID, w); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *ExportHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=bitewise_export.json")

	if err := h.svc.ExportJSON(r.Context(), userID, w); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *ExportHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	if err := h.svc.DeleteAccount(r.Context(), userID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.NoContent(w)
}
