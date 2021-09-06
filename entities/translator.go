package entities

import "time"

type Translator struct {
	ID        string    `json:"_id"`
	Key       string    `json:"_key"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Hpwd      string    `json:"hpwd"`
	Activated bool      `json:"activated"`
	Suspended bool      `json:"suspended"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewTranslator struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Hpwd      string    `json:"hpwd"`
	Activated bool      `json:"activated"`
	Suspended bool      `json:"suspended"`
	CreatedAt time.Time `json:"createdAt"`
}
