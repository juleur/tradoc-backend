package entities

type TranslationsBody struct {
	Translations []*Translation `json:"translations"`
}

type Translation struct {
	Oc           string `json:"oc"`
	Fr           string `json:"fr"`
	En           string `json:"en"`
	DatasetID    string `json:"datasetID"`
	TranslatorID string `json:"translatorID"`
	FullDialect  string `json:"fullDialect"`
}

type TranslationFile struct {
	Dialect        string           `json:"dialect" bson:"dialect"`
	SubdialectFile []SubdialectFile `json:"subdialects" bson:"subdialects"`
}

type SubdialectFile struct {
	Name  string `json:"name" bson:"name"`
	Files Files  `json:"files" bson:"files"`
}

type Files struct {
	Fr string `json:"fr" bson:"fr"`
	En string `json:"en" bson:"en"`
}
