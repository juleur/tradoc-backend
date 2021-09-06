package entities

type Dialect struct {
	Name        string   `json:"name"`
	Subdialects []string `json:"subdialects"`
}

type DialectSubdialectDocument struct {
	Dialect    string `json:"dialect"`
	Subdialect string `json:"subdialect"`
}

type DialectSubdialects struct {
	DialectID   string       `json:"dialect_id"`
	Dialect     string       `json:"dialect"`
	Subdialects []Subdialect `json:"subdialects"`
}

type Subdialect struct {
	SubdialectID                 string `json:"subdialect_id"`
	Subdialect                   string `json:"subdialect"`
	TotalTranslations            int    `json:"total_translations"`
	TotalTranslationByTranslator int    `json:"total_translation_by_translator"`
}
