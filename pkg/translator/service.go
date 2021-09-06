package translator

import "btradoc/entities"

type Service interface {
	FindTranslatorByUsername(username string) (*entities.Translator, error)
	CreateTranslator(newTranslator entities.NewTranslator) error
	SetRefreshToken(translatorID string) (string, error)
	FindRefreshToken(refreshToken string) (*entities.Translator, error)
	DeleteRefreshToken(refreshToken string) error
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
