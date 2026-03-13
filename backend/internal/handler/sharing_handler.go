package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type SharingHandler struct {
	sharing *service.SharingService
}

func NewSharingHandler(sharing *service.SharingService) *SharingHandler {
	return &SharingHandler{sharing: sharing}
}

func (h *SharingHandler) ShareRecipe(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	recipeID, err := parseIDParam(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	shared, err := h.sharing.ShareRecipe(r.Context(), recipeID, userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, shared)
}

func (h *SharingHandler) GetSharedRecipe(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		httputil.Error(w, http.StatusBadRequest, "missing share code")
		return
	}

	recipe, err := h.sharing.GetSharedRecipe(r.Context(), code)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, recipe)
}

func (h *SharingHandler) SaveSharedRecipe(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := chi.URLParam(r, "code")
	if code == "" {
		httputil.Error(w, http.StatusBadRequest, "missing share code")
		return
	}

	recipe, err := h.sharing.SaveSharedRecipe(r.Context(), code, userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.Created(w, recipe)
}
