package models

type Tokens struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refreshToken"`
}
