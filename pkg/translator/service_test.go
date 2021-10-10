package translator

import (
	"btradoc/entities"
	"btradoc/storage/mongodb"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
)

func TestCreateTranslator(t *testing.T) {
	db := mongodb.NewMongoClient()

	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	newTranslator := entities.NewTranslator{
		Email:    "lois@sfr.fr",
		Username: "loisssvfvfvfvf",
		Hpwd:     "$argon2id$v=19$m=16,t=2,p=1$RXhFUTNTWGFueUI5Umx6Qg$wtoGpTHQ43Vt5kFiR342Xg",
		SecretQuestions: []entities.SecretQuestion{
			{
				Question: "Qué siguèt ton faus-nom quand ères un enfant ?",
				Response: "$argon2id$v=19$m=16,t=2,p=1$NkFLZFFPeEllS2VLeGpjWg$HlmRZGzfHlpD8ItbTdWABA",
			},
			{
				Question: "Qué siguèt lo premier filme que veguères au cinemà ?",
				Response: "$argon2id$v=19$m=16,t=2,p=1$dW54UThKWjJQTU9WM05Tdw$7KTR8EuJOF9QZ+jqtMDAPw",
			},
		},
	}

	err := translatorService.CreateTranslator(newTranslator)
	assert.Nil(t, err)
}

func TestGetTranslatorByUsername(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	username := "Dàrius"
	translator, err := translatorService.FindTranslatorByUsername(username)
	assert.Nil(t, err)
	t.Logf("%+v\n", translator)
}

func TestCreateRefreshToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	translatorID := "6148c3f1ba78b40cdeb49289"
	refreshToken, err := translatorService.SetRefreshToken(translatorID)
	assert.Nil(t, err)
	t.Logf("%s\n", refreshToken)
}

func TestGetRefreshToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	refreshToken := "JfQt53reoJR"

	transl, err := translatorService.FindRefreshToken(refreshToken)
	assert.Nil(t, err)
	t.Logf("%+v\n", transl)
}

func TestRemoveRefreshToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	refreshToken := "JfQtQ53reoJR"

	err := translatorService.DeleteRefreshToken(refreshToken)

	assert.Nil(t, err)
}

func TestFetchSecretQuestionsByToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	token := "0yCCKCPlRgJQnP5jDG22"
	// token := "pf4r7azacvn9h7"
	sq, err := translatorService.FetchSecretQuestionsByToken(token)
	assert.Nil(t, err)
	t.Logf("%s\n", sq)
}

func TestFetchSecretQuestions(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	secretQuestions, err := translatorService.FetchSecretQuestions()
	assert.Nil(t, err)
	t.Logf("%s", secretQuestions)
}

func TestUpdatePassword(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	translatorID := "6148c3f1ba78b40cdeb49289"
	hashedPaswword, err := argon2id.CreateHash("test", argon2id.DefaultParams)
	assert.Nil(t, err)

	err = translatorService.ResetPassword(translatorID, hashedPaswword)
	assert.Nil(t, err)
}

func TestProceedResetPassword(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	email := "dubois-darius@sfr.fr"
	transl, err := translatorService.ProceedResetPassword(email)
	assert.Nil(t, err)
	t.Logf("%+v\n", transl)
}
