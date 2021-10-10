package entities

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Translator struct {
	ID              string           `json:"id"`
	Email           string           `json:"email"`
	Username        string           `json:"username"`
	Hpwd            string           `json:"hpwd"`
	Confirmed       bool             `json:"confirmed"`
	Suspended       bool             `json:"suspended"`
	CreatedAt       time.Time        `json:"createdAt"`
	Permissions     []string         `json:"permissions"`
	SecretQuestions []SecretQuestion `json:"secretQuestions,omitempty"`
}

func (t *Translator) CompressPerms() []string {
	var dperms []string
	for _, dp := range t.Permissions {
		dialectSubdialect := strings.Split(dp, "_")
		// convert to runes for extended ASCII characters
		dialectRunes := []rune(dialectSubdialect[0])
		subdialectRunes := []rune(dialectSubdialect[1])
		dperm := string(dialectRunes[:3]) + "_" + string(subdialectRunes[:3])
		dperms = append(dperms, dperm)
	}
	return dperms
}

type NewTranslator struct {
	Email           string
	Username        string
	Hpwd            string
	Confirmed       bool
	Suspended       bool
	SecretQuestions []SecretQuestion
}

type TranslatorResetPassword struct {
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Token    string `bson:"token"`
}

type TranslatorSecretQuestions struct {
	ID              primitive.ObjectID
	SecretQuestions map[string]string
}
