package models

type MainMenu struct {
	Dialects               []Dialect `json:"dialects"`
	TotalOnlineTranslators int       `json:"totalOnlineTranslators"`
}
