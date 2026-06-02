package application

import (
	"context"

	friendshipdomain "ditza/internal/friendship/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type Service struct {
	friendshipRepository friendshipdomain.Repository
}

type SendRequestCommand struct {
	RequesterID valueobject.UserID
	AddresseeID valueobject.UserID
}

func NewService(friendshipRepository friendshipdomain.Repository) *Service {
	return &Service{friendshipRepository: friendshipRepository}
}

func (s *Service) SendRequest(ctx context.Context, command SendRequestCommand) (friendshipEntity *friendshipdomain.Friendship, err error) {
	tracker := monitoring.StartService(logger.ModelFriendship, "send_request", map[string]any{
		"requester_id": command.RequesterID,
		"addressee_id": command.AddresseeID,
	})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"friendship_id": friendshipEntity.ID})
	}()

	existing, err := s.friendshipRepository.FindBetweenUsers(ctx, command.RequesterID, command.AddresseeID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerror.New("FRIENDSHIP_EXISTS", "ya existe una solicitud o amistad con este usuario", domainerror.ErrFriendshipExists)
	}

	friendshipEntity, err = friendshipdomain.NewRequest(command.RequesterID, command.AddresseeID)
	if err != nil {
		return nil, err
	}
	if err := s.friendshipRepository.Create(ctx, friendshipEntity); err != nil {
		return nil, err
	}
	return friendshipEntity, nil
}

func (s *Service) Accept(ctx context.Context, friendshipID valueobject.FriendshipID, responderID valueobject.UserID) (err error) {
	tracker := monitoring.StartService(logger.ModelFriendship, "accept", map[string]any{
		"friendship_id": friendshipID,
		"responder_id":  responderID,
	})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	friendshipEntity, err := s.friendshipRepository.FindByID(ctx, friendshipID)
	if err != nil {
		return err
	}
	if friendshipEntity == nil {
		return domainerror.New("FRIENDSHIP_NOT_FOUND", "solicitud de amistad no encontrada", domainerror.ErrNotFound)
	}
	if err := friendshipEntity.Accept(responderID); err != nil {
		return err
	}
	return s.friendshipRepository.Update(ctx, friendshipEntity)
}

func (s *Service) Reject(ctx context.Context, friendshipID valueobject.FriendshipID, responderID valueobject.UserID) (err error) {
	tracker := monitoring.StartService(logger.ModelFriendship, "reject", map[string]any{
		"friendship_id": friendshipID,
		"responder_id":  responderID,
	})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	friendshipEntity, err := s.friendshipRepository.FindByID(ctx, friendshipID)
	if err != nil {
		return err
	}
	if friendshipEntity == nil {
		return domainerror.New("FRIENDSHIP_NOT_FOUND", "solicitud de amistad no encontrada", domainerror.ErrNotFound)
	}
	if err := friendshipEntity.Reject(responderID); err != nil {
		return err
	}
	return s.friendshipRepository.Update(ctx, friendshipEntity)
}

func (s *Service) ListFriends(ctx context.Context, userID valueobject.UserID) (friends []friendshipdomain.Friendship, err error) {
	tracker := monitoring.StartService(logger.ModelFriendship, "list_friends", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(friends)})
	}()

	friends, err = s.friendshipRepository.FindAcceptedByUserID(ctx, userID)
	return friends, err
}

func (s *Service) ListPendingRequests(ctx context.Context, userID valueobject.UserID) (requests []friendshipdomain.Friendship, err error) {
	tracker := monitoring.StartService(logger.ModelFriendship, "list_pending_requests", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(requests)})
	}()

	requests, err = s.friendshipRepository.FindPendingByAddresseeID(ctx, userID)
	return requests, err
}
