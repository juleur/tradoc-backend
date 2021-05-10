package models

type Dialect struct {
	Nom         string       `json:"nom"`
	Subdialects []Subdialect `json:"subdialects"`
}

type Subdialect struct {
	Nom                      string `json:"nom"`
	TotalSentencesTranslated int    `json:"total_sentences_translated"`
	Abbr                     string `json:"abbr"`
}
