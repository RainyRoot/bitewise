package handler

import (
	"net/http"

	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type AchievementHandler struct {
	achievements *service.AchievementService
}

func NewAchievementHandler(achievements *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievements: achievements}
}

func (h *AchievementHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	achievements, err := h.achievements.GetAll(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, achievements)
}

func (h *AchievementHandler) GetMine(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	achievements, err := h.achievements.GetUserAchievements(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, achievements)
}
