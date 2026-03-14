package service

import (
	"context"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type FriendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *FriendService) InviteFriend(ctx context.Context, userID int64, email string) (*domain.FriendInvite, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Email == email {
		return nil, fmt.Errorf("cannot invite yourself")
	}

	inv := &domain.FriendInvite{
		FromUserID: userID,
		ToEmail:    email,
	}
	if err := s.friendRepo.CreateInvite(ctx, inv); err != nil {
		return nil, fmt.Errorf("create invite: %w", err)
	}
	return inv, nil
}

func (s *FriendService) GetPendingInvites(ctx context.Context, userID int64) ([]domain.FriendInvite, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return s.friendRepo.GetPendingInvites(ctx, user.Email)
}

func (s *FriendService) RespondToInvite(ctx context.Context, userID, inviteID int64, accept bool) error {
	if !accept {
		return s.friendRepo.DeclineInvite(ctx, inviteID, userID)
	}

	if err := s.friendRepo.AcceptInvite(ctx, inviteID, userID); err != nil {
		return fmt.Errorf("accept invite: %w", err)
	}

	inv, err := s.friendRepo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("get invite: %w", err)
	}
	if inv == nil {
		return fmt.Errorf("invite not found")
	}

	return s.friendRepo.AddFriendship(ctx, userID, inv.FromUserID)
}

func (s *FriendService) GetFriends(ctx context.Context, userID int64) ([]domain.FriendInfo, error) {
	return s.friendRepo.GetFriends(ctx, userID)
}

func (s *FriendService) RemoveFriend(ctx context.Context, userID, friendID int64) error {
	return s.friendRepo.RemoveFriend(ctx, userID, friendID)
}

func (s *FriendService) GetLeaderboard(ctx context.Context, userID int64) ([]domain.LeaderboardEntry, error) {
	return s.friendRepo.GetLeaderboard(ctx, userID)
}
