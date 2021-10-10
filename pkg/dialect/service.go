package dialect

import "btradoc/entities"

type Service interface {
	FetchDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error) {
	return s.repository.GetDialectsSubdialect(translatorID)
}
