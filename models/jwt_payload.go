package models

import "github.com/gbrlsnchs/jwt/v3"

type JWTPayload struct {
	jwt.Payload
	TranslatorID int      `json:"translatorId"`
	Username     string   `json:"username"`
	Dialects     []string `json:"dialects"` // si utilisateur est vérifié
}
