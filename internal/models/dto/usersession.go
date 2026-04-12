package dto

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	UserID int `json:"user_id"`
}
