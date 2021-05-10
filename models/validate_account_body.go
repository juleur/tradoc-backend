package models

type ValidateAccountBody struct {
	Username string `json:"username"`
	Code     int    `json:"code"`
	Password string `json:"password"`
}
