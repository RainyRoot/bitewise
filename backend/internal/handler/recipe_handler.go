package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type RecipeHandler struct {
	recipes *service.RecipeService
}

func NewRecipeHandler(recipes *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipes: recipes}
}

func (h *RecipeHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := domain.RecipeFilter{
		Query:    r.URL.Query().Get("q"),
		Category: r.URL.Query().Get("category"),
	}

	if v := r.URL.Query().Get("allergen_free"); v != "" {
		filter.AllergenFree = strings.Split(v, ",")
	}

	if v := r.URL.Query().Get("max_prep_time"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			filter.MaxPrepTime = i
		}
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			filter.Limit = i
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			filter.Offset = i
		}
	}

	recipes, total, err := h.recipes.Search(r.Context(), filter)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"recipes": recipes,
		"total":   total,
	})
}

func (h *RecipeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	userID := UserIDFromContext(r.Context())
	recipe, err := h.recipes.GetByID(r.Context(), id, userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, recipe)
}

func (h *RecipeHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	recipeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	if err := h.recipes.AddFavorite(r.Context(), userID, recipeID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}

func (h *RecipeHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	recipeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	if err := h.recipes.RemoveFavorite(r.Context(), userID, recipeID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.NoContent(w)
}

func (h *RecipeHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			limit = i
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			offset = i
		}
	}

	recipes, err := h.recipes.GetFavorites(r.Context(), userID, limit, offset)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, recipes)
}
