package http

type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
	UserID      string `json:"user_id"`
	Alias       string `json:"alias"`
	Email       string `json:"email"`
}

type RegisterRequestDTO struct {
	Alias    string `json:"alias"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileResponseDTO struct {
	UserID string `json:"user_id"`
	Alias  string `json:"alias"`
	Email  string `json:"email"`
	Coins  int    `json:"coins"`
}
