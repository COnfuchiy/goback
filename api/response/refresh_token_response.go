package response

type RefreshTokenResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
