package http

type SendFriendRequestDTO struct {
	AddresseeID uint64 `json:"addressee_id"`
}

type FriendshipResponseDTO struct {
	FriendshipID uint64 `json:"friendship_id"`
	RequesterID  uint64 `json:"requester_id"`
	AddresseeID  uint64 `json:"addressee_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	RespondedAt  string `json:"responded_at,omitempty"`
}
