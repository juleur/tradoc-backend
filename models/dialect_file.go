package models

import "time"

type AllFiles struct {
	DialectFiles        []DialectFile `json:"dialectFiles"`
	LastDatetimeGenFile *time.Time    `json:"lastDatetimeGenFile"`
}

type DialectFile struct {
	Nom             string           `json:"nom"`
	SubdialectFiles []SubdialectFile `json:"subdialectFiles"`
}

type SubdialectFile struct {
	Nom        string  `json:"nom"`
	FilepathFr *string `json:"filepathFr"`
	FilepathEn *string `json:"filepathEn"`
}
