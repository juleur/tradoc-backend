package translation

import "btradoc/entities"

type Service interface {
	FetchSentencesToTranslate(translatorID string, dialectName string, subdialectName string) (*[]entities.Dataset, error)
	AddTranslations(translations []entities.Translation) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchSentencesToTranslate(translatorID string, dialectName string, subdialectName string) (*[]entities.Dataset, error) {
	return s.repository.GetDatasets(translatorID, dialectName, subdialectName)
}

func (s *service) AddTranslations(translations []entities.Translation) error {
	return s.repository.InsertTranslations(translations)
}
