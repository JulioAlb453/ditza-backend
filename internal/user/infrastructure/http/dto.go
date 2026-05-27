package http

type RegisterRequestDTO struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Timezone     string `json:"timezone"`
}

type RegisterResponseDTO struct {
	UserID     uint64 `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Timezone   string `json:"timezone"`
	FriendCode string `json:"friend_code"`
}

type UserProfileResponseDTO struct {
	UserID     uint64 `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Timezone   string `json:"timezone"`
	Coins      int    `json:"coins"`
	FriendCode string `json:"friend_code"`
}
