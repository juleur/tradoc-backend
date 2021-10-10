package translator

import "btradoc/entities"

type Service interface {
	FindTranslatorByUsername(username string) (*entities.Translator, error)
	CreateTranslator(newTranslator entities.NewTranslator) error
	SetRefreshToken(translatorID string) (string, error)
	FindRefreshToken(refreshToken string) (*entities.Translator, error)
	DeleteRefreshToken(refreshToken string) error
	FetchSecretQuestionsByToken(token string) (*entities.TranslatorSecretQuestions, error)
	FetchSecretQuestions() ([]string, error)
	ProceedResetPassword(email string) (*entities.TranslatorResetPassword, error)
	ResetPassword(translatorID string, hashedPassword string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FindTranslatorByUsername(username string) (*entities.Translator, error) {
	return s.repository.GetTranslatorByUsername(username)
}

func (s *service) CreateTranslator(newTranslator entities.NewTranslator) error {
	return s.repository.InsertNewTranslator(newTranslator)
}

func (s *service) SetRefreshToken(translatorID string) (string, error) {
	return s.repository.InsertNewRefreshToken(translatorID)
}

func (s *service) FindRefreshToken(refreshToken string) (*entities.Translator, error) {
	return s.repository.GetRefreshToken(refreshToken)
}

func (s *service) DeleteRefreshToken(refreshToken string) error {
	return s.repository.RemoveRefreshToken(refreshToken)
}

func (s *service) FetchSecretQuestionsByToken(token string) (*entities.TranslatorSecretQuestions, error) {
	return s.repository.GetSecretQuestionsByToken(token)
}

func (s *service) FetchSecretQuestions() ([]string, error) {
	return s.repository.GetSecretQuestions()
}

func (s *service) ProceedResetPassword(email string) (*entities.TranslatorResetPassword, error) {
	return s.repository.CreateTokenResetPassword(email)
}

func (s *service) ResetPassword(translatorID string, hashedPassword string) error {
	return s.repository.UpdatePassword(translatorID, hashedPassword)
}
