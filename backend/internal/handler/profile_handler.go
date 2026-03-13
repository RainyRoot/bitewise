package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

// ProfileHandler handles profile HTTP requests.
type ProfileHandler struct {
	profile *service.ProfileService
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(profile *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profile: profile}
}

// GetProfile handles GET /profile.
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.profile.GetProfile(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, user)
}

// UpdateProfile handles PUT /profile.
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.profile.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, user)
}

// GetAllergies handles GET /profile/allergies.
func (h *ProfileHandler) GetAllergies(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	allergies, err := h.profile.GetAllergies(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, allergies)
}

// SetAllergies handles PUT /profile/allergies.
func (h *ProfileHandler) SetAllergies(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var allergens []string
	if err := json.NewDecoder(r.Body).Decode(&allergens); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.profile.SetAllergies(r.Context(), userID, allergens); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}

// GetPreferences handles GET /profile/preferences.
func (h *ProfileHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	prefs, err := h.profile.GetPreferences(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, prefs)
}

// SetPreferences handles PUT /profile/preferences.
func (h *ProfileHandler) SetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var prefs []domain.UserPreference
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.profile.SetPreferences(r.Context(), userID, prefs); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}
