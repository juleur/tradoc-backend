package tasks

import (
	"log"
	"testing"
	"tradoc/db"
)

func TestCleanOldTranslations(t *testing.T) {
	db := db.OpenDB()
	dirPath := "/Users/rd/files/translations"
	errCh := CleanOldTranslations(db, dirPath)
	for e := range errCh {
		log.Println(e)
	}
}
