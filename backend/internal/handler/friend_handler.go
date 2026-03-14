package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

type FriendHandler struct {
	svc *service.FriendService
}

func NewFriendHandler(svc *service.FriendService) *FriendHandler {
	return &FriendHandler{svc: svc}
}

func (h *FriendHandler) InviteFriend(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	var req domain.FriendInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inv, err := h.svc.InviteFriend(r.Context(), userID, req.Email)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.Created(w, inv)
}

func (h *FriendHandler) GetPendingInvites(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	invites, err := h.svc.GetPendingInvites(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, invites)
}

func (h *FriendHandler) RespondToInvite(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	inviteID, err := parseIDParam(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid invite id")
		return
	}

	var body struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.RespondToInvite(r.Context(), userID, inviteID, body.Accept); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.NoContent(w)
}

func (h *FriendHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	friends, err := h.svc.GetFriends(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, friends)
}

func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	friendID, err := parseIDParam(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid friend id")
		return
	}

	if err := h.svc.RemoveFriend(r.Context(), userID, friendID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.NoContent(w)
}

func (h *FriendHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	entries, err := h.svc.GetLeaderboard(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, entries)
}
