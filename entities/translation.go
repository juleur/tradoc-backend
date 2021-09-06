package entities

type Translation struct {
	Oc           string `json:"oc"`
	Fr           string `json:"fr"`
	En           string `json:"en"`
	DatasetID    string `json:"dataset"`
	TranslatorID string `json:"translator"`
	SubdialectID string `json:"subdialect"`
}
