package entities

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterBody struct {
	Username        string           `json:"username"`
	Email           string           `json:"email"`
	Password        string           `json:"password"`
	SecretQuestions []SecretQuestion `json:"secretQuestions"`
}

type ResetPasswordBody struct {
	Email string `json:"email"`
}

type SecretQuestion struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

type UpdatePasswordBody struct {
	Token           string           `json:"token"`
	SecretQuestions []SecretQuestion `json:"secretQuestions"`
	Password        string           `json:"password"`
}
