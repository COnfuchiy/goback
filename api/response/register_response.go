package response

type RegisterResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
