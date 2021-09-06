package translator

import (
	"btradoc/entities"
	"btradoc/tools"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTranslator(t *testing.T) {
	db := tools.ArangoDBConnection()

	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	newTranslator := entities.NewTranslator{
		Email:     "test2-test2@sfr.fr",
		Username:  "test2 test2",
		Hpwd:      "$argon2id$v=19$m=16,t=2,p=1$VjM2OW4xNnJGSzNQdzE5dw$NVI8+jmyE+hzn7udhAMfZA",
		CreatedAt: time.Now(),
	}

	err := translatorService.CreateTranslator(newTranslator)
	assert.Nil(t, err)
}

func TestCreateRefreshToken(t *testing.T) {
	db := tools.ArangoDBConnection()

	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	translatorID := "translators/123456"
	refreshToken, err := translatorService.SetRefreshToken(translatorID)
	assert.Nil(t, err)
	t.Log(refreshToken)
}

func TestGetRefreshToken(t *testing.T) {
	db := tools.ArangoDBConnection()

	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	refreshToken := "HJm1ugQCq7"
	translator, err := translatorService.FindRefreshToken(refreshToken)
	assert.Nil(t, err)
	t.Log(translator)
}
