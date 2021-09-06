package translator

import (
	"btradoc/entities"
	"btradoc/tools"
	"context"
	"errors"

	arangodb "github.com/arangodb/go-driver"
)

type Repository interface {
	GetTranslatorByUsername(username string) (*entities.Translator, error)
	InsertNewTranslator(newTranslator entities.NewTranslator) error
	InsertNewRefreshToken(translatorID string) (string, error)
	GetRefreshToken(refreshToken string) (*entities.Translator, error)
	RemoveRefreshToken(refreshToken string) error
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangoDB arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangoDB,
	}
}

func (r *repository) GetTranslatorByUsername(username string) (*entities.Translator, error) {
	query := `FOR t IN translators FILTER t.username == @username return t`
	bindVars := map[string]interface{}{
		"username": username,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New("A REMPLIR")
	}
	defer cursor.Close()

	translator := entities.Translator{}
	_, err = cursor.ReadDocument(context.Background(), &translator)
	if arangodb.IsNoMoreDocuments(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.New("A REMPLIR")
	}

	return &translator, nil
}

func (r *repository) InsertNewTranslator(newTranslator entities.NewTranslator) error {
	query := "FOR t IN translators FILTER t.username == @username || t.email == @email RETURN t"
	bindVars := map[string]interface{}{
		"username": newTranslator.Username,
		"email":    newTranslator.Email,
	}

	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return errors.New("A REMPLIR")
	}
	defer cursor.Close()

	translator := entities.Translator{}
	_, err = cursor.ReadDocument(context.Background(), &translator)
	if arangodb.IsNoMoreDocuments(err) {
		col, err := r.ArangoDB.Collection(context.Background(), "translators")
		if err != nil {
			return err
		}
		_, err = col.CreateDocument(context.Background(), newTranslator)
		if err != nil {
			return err
		}

		return err
	} else if err != nil {
		return errors.New("A REMPLIR")
	}

	if newTranslator.Username == translator.Username {
		return errors.New("ce pseudo est déjà utilisé")
	}

	return errors.New("cet email est déjà utilisé")
}

func (r *repository) InsertNewRefreshToken(translatorID string) (string, error) {
	for i := 0; i < 50; i++ {
		refreshToken := tools.GenerateID(10)

		query := `INSERT {
			_key: @refreshToken, "translator": @translatorID, "createdAt": DATE_NOW()
		} IN refreshTokens OPTIONS { conflict: true }`
		bindVars := map[string]interface{}{
			"translatorID": translatorID,
			"refreshToken": refreshToken,
		}

		if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil {
			if arangodb.IsConflict(err) {
				continue
			}

			return "", err
		}

		return refreshToken, nil
	}
	return "", errors.New("INSERTNEWREFRESHTOKEN ERROR")
}

func (r *repository) GetRefreshToken(refreshToken string) (*entities.Translator, error) {
	query := `
		FOR rt IN refreshTokens
			FILTER rt._key == @key
			FOR t IN translators
				FILTER t._id == rt.translator
				RETURN t
	`
	bindVars := map[string]interface{}{
		"key": refreshToken,
	}

	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New("A REMPLIR")
	}
	defer cursor.Close()

	translator := entities.Translator{}
	_, err = cursor.ReadDocument(context.Background(), &translator)
	if arangodb.IsNoMoreDocuments(err) {
		return nil, errors.New("ce refresh token n'existe pas")
	} else if err != nil {
		return nil, errors.New("A REMPLIR")
	}

	return &translator, nil
}

func (r *repository) RemoveRefreshToken(refreshToken string) error {
	query := `FOR rt IN refreshTokens FILTER rt._key == @key REMOVE rt IN refreshTokens `
	bindVars := map[string]interface{}{
		"key": refreshToken,
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil {
		return err
	}
	return nil
}
