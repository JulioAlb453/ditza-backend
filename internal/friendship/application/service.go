package application

import (
	"context"

	friendshipdomain "ditza/internal/friendship/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
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

func (s *Service) SendRequest(ctx context.Context, command SendRequestCommand) (*friendshipdomain.Friendship, error) {
	existing, err := s.friendshipRepository.FindBetweenUsers(ctx, command.RequesterID, command.AddresseeID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerror.New("FRIENDSHIP_EXISTS", "ya existe una solicitud o amistad con este usuario", domainerror.ErrFriendshipExists)
	}

	friendshipEntity, err := friendshipdomain.NewRequest(command.RequesterID, command.AddresseeID)
	if err != nil {
		return nil, err
	}
	if err := s.friendshipRepository.Create(ctx, friendshipEntity); err != nil {
		return nil, err
	}
	return friendshipEntity, nil
}

func (s *Service) Accept(ctx context.Context, friendshipID valueobject.FriendshipID, responderID valueobject.UserID) error {
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

func (s *Service) Reject(ctx context.Context, friendshipID valueobject.FriendshipID, responderID valueobject.UserID) error {
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

func (s *Service) ListFriends(ctx context.Context, userID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return s.friendshipRepository.FindAcceptedByUserID(ctx, userID)
}

func (s *Service) ListPendingRequests(ctx context.Context, userID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return s.friendshipRepository.FindPendingByAddresseeID(ctx, userID)
}
