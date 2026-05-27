package domain

import (
	"time"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Friendship struct {
	ID          valueobject.FriendshipID
	RequesterID valueobject.UserID
	AddresseeID valueobject.UserID
	Status      Status
	CreatedAt   time.Time
	RespondedAt *time.Time
}

func NewRequest(requesterID, addresseeID valueobject.UserID) (*Friendship, error) {
	if requesterID == addresseeID {
		return nil, domainerror.New("CANNOT_FRIEND_SELF", "no puedes agregarte a ti mismo como amigo", domainerror.ErrCannotFriendSelf)
	}

	return &Friendship{
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      StatusPending,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (f *Friendship) Involves(userID valueobject.UserID) bool {
	return f.RequesterID == userID || f.AddresseeID == userID
}

func (f *Friendship) IsAccepted() bool {
	return f.Status == StatusAccepted
}

func (f *Friendship) Accept(responderID valueobject.UserID) error {
	if f.AddresseeID != responderID {
		return domainerror.New("UNAUTHORIZED", "solo el destinatario puede aceptar la solicitud", domainerror.ErrUnauthorized)
	}
	if !f.Status.CanAccept() {
		return domainerror.New("FRIENDSHIP_NOT_PENDING", "la solicitud de amistad no está pendiente", domainerror.ErrFriendshipNotPending)
	}

	now := time.Now().UTC()
	f.Status = StatusAccepted
	f.RespondedAt = &now
	return nil
}

func (f *Friendship) Reject(responderID valueobject.UserID) error {
	if f.AddresseeID != responderID {
		return domainerror.New("UNAUTHORIZED", "solo el destinatario puede rechazar la solicitud", domainerror.ErrUnauthorized)
	}
	if !f.Status.CanReject() {
		return domainerror.New("FRIENDSHIP_NOT_PENDING", "la solicitud de amistad no está pendiente", domainerror.ErrFriendshipNotPending)
	}

	now := time.Now().UTC()
	f.Status = StatusRejected
	f.RespondedAt = &now
	return nil
}

func (f *Friendship) FriendIDFor(userID valueobject.UserID) (valueobject.UserID, bool) {
	if !f.IsAccepted() {
		return "", false
	}
	if f.RequesterID == userID {
		return f.AddresseeID, true
	}
	if f.AddresseeID == userID {
		return f.RequesterID, true
	}
	return "", false
}
