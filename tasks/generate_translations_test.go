package tasks

import (
	"os"
	"testing"
	"tradoc/db"
)

func TestGenerateTranslations(t *testing.T) {
	homePath := os.Getenv("HOME")

	TRANSLATIONS_PATH := homePath + "/" + "files/translations"

	database := db.OpenDB()
	dbRepo := db.NewDBPsql(database)

	dialectsTableName, err := dbRepo.FindAllDialect()
	if err != nil {
		t.Fatalf(err.Error.Error())
	}

	errCh := GenerateTranslatedSentences(database, dialectsTableName, TRANSLATIONS_PATH)
	for e := range errCh {
		t.Log(e.Error())
	}
}
