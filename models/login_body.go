package models

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
