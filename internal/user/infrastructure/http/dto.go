package http

type RegisterRequestDTO struct {
	Alias    string `json:"alias"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	UserID string `json:"user_id"`
	Alias  string `json:"alias"`
	Email  string `json:"email"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	UserID string `json:"user_id"`
	Alias  string `json:"alias"`
	Email  string `json:"email"`
}

type UserProfileResponseDTO struct {
	UserID string `json:"user_id"`
	Alias  string `json:"alias"`
	Email  string `json:"email"`
	Coins  int    `json:"coins"`
}
