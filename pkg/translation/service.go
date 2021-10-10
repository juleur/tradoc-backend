package translation

import "btradoc/entities"

type Service interface {
	FetchDatasets(fullDialect string) (*[]entities.Dataset, error)
	AddDatasetNewFullDialect(translations []*entities.Translation) error
	AddTranslations(translations []*entities.Translation) error
	FetchTotalOnGoingTranslations(fullDialect, translatorID string) (int, error)
	AddOnGoingTranslations(fullDialect, translatorID string, datasets *[]entities.Dataset) error
	RemoveOnGoingTranslations(translations []*entities.Translation) error
	FetchTranslationsFiles() (*[]entities.TranslationFile, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchDatasets(fullDialect string) (*[]entities.Dataset, error) {
	return s.repository.GetDatasets(fullDialect)
}

func (s *service) AddDatasetNewFullDialect(translations []*entities.Translation) error {
	return s.repository.AddNewFullDialectToDataset(translations)
}

func (s *service) AddTranslations(translations []*entities.Translation) error {
	return s.repository.InsertTranslations(translations)
}

func (s *service) FetchTotalOnGoingTranslations(fullDialect, translatorID string) (int, error) {
	return s.repository.GetTotalOnGoingTranslation(fullDialect, translatorID)
}

func (s *service) AddOnGoingTranslations(fullDialect, translatorID string, datasets *[]entities.Dataset) error {
	return s.repository.InsertDatasetsOnGoingTranslations(fullDialect, translatorID, datasets)
}

func (s *service) RemoveOnGoingTranslations(translations []*entities.Translation) error {
	return s.repository.RemoveDatasetsOnGoingTranslations(translations)
}

func (s *service) FetchTranslationsFiles() (*[]entities.TranslationFile, error) {
	return s.repository.GetTranslationsFiles()
}
