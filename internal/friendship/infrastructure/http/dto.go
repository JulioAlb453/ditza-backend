package http

type SendFriendRequestDTO struct {
	AddresseeID string `json:"addressee_id"`
}

type FriendshipResponseDTO struct {
	FriendshipID uint64 `json:"friendship_id"`
	RequesterID  string `json:"requester_id"`
	AddresseeID  string `json:"addressee_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	RespondedAt  string `json:"responded_at,omitempty"`
}
