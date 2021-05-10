package models

type NewTranslationBody struct {
	Translation Translation `json:"translation"`
}

type Translation struct {
	ID       int       `json:"id"`
	Abbr     string    `json:"abbr"`
	Default  Default   `json:"default"`
	Feminize *Feminize `json:"feminize,omitempty"`
}

type Default struct {
	Occitan string  `json:"occitan"`
	French  string  `json:"french"`
	English *string `json:"english,omitempty"`
}

type Feminize struct {
	Occitan string  `json:"occitan"`
	French  string  `json:"french"`
	English *string `json:"english,omitempty"`
}
