package tasks

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
	"tradoc/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TranslatedSentence struct {
	ID      int
	Occitan string
	French  string
	English string
}

func GenerateTranslatedSentences(db *pgxpool.Pool, dialectsTableName []string, dirPath string) <-chan error {
	chErr := make(chan error, 10)
	go func() {
		for _, dialectTableName := range dialectsTableName {
			fmt.Println("GenerateTranslatedSentences: " + dialectTableName)
			translatedSentences, err := getAllTranslationsByDialect(db, dialectTableName)
			if err != nil {
				chErr <- err
				continue
			}

			var hasAtLeastOneEnglishSentence bool

			err = createDirectory(dirPath)
			if err != nil {
				chErr <- err
				continue
			}

			filepathFr, filepathEn := generateFilepaths(dirPath, dialectTableName)

			fileFr, err := createFile(filepathFr)
			if err != nil {
				chErr <- err
				continue
			}
			defer fileFr.Close()

			writerFr := bufio.NewWriter(fileFr)
			defer writerFr.Flush()

			for _, tr := range translatedSentences {
				line := fmt.Sprintf("%s~%s\n", tr.Occitan, tr.French)
				if _, err := writerFr.WriteString(line); err != nil {
					chErr <- err
					continue
				}

				if len(tr.English) > 0 {
					fileEn, err := createFile(filepathEn)
					if err != nil {
						chErr <- err
						continue
					}
					defer fileEn.Close()

					writerEn := bufio.NewWriter(fileEn)
					defer writerEn.Flush()

					line := fmt.Sprintf("%s~%s\n", tr.Occitan, tr.English)
					if _, err = writerEn.WriteString(line); err != nil {
						chErr <- err
						continue
					}
					hasAtLeastOneEnglishSentence = true
				}

				if hasAtLeastOneEnglishSentence {
					/// ajouter les fichier en base de donnée
					if _, err := db.Exec(context.Background(), `
					INSERT INTO translation_files(dialect_name, filename_fr, filename_en) VALUES ($1,$2,$3)
				`, dialectTableName, utils.ExtractFilenameOnly(filepathFr), utils.ExtractFilenameOnly(filepathEn)); err != nil {
						chErr <- err
					}
				} else {
					if _, err := db.Exec(context.Background(), `
					INSERT INTO translation_files(dialect_name, filename_fr) VALUES ($1,$2)
				`, dialectTableName, utils.ExtractFilenameOnly(filepathFr)); err != nil {
						chErr <- err
					}
				}
			}
		}
		close(chErr)
	}()
	return chErr
}

func getAllTranslationsByDialect(db *pgxpool.Pool, dialectTableName string) ([]TranslatedSentence, error) {
	splited := strings.Split(dialectTableName, "_")
	query := fmt.Sprintf(`
		SELECT id, frasa_%s_%s, frasa_fr, frasa_an FROM %s
		ORDER BY tradusit_lo DESC
	`, splited[0][0:3], splited[1][0:3], dialectTableName)
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return []TranslatedSentence{}, err
	}

	translatedSentences := []TranslatedSentence{}
	for rows.Next() {
		tr := TranslatedSentence{}
		english := &sql.NullString{}
		if err := rows.Scan(&tr.ID, &tr.Occitan, &tr.French, english); err != nil {
			return []TranslatedSentence{}, err
		}
		if english.Valid {
			tr.English = english.String
		}

		translatedSentences = append(translatedSentences, tr)
	}
	rows.Close()

	if len(translatedSentences) == 0 {
		return []TranslatedSentence{}, fmt.Errorf("no sentences found in %s", dialectTableName)
	}
	return translatedSentences, nil
}

func createDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, os.FileMode(0755))
}

func generateFilepaths(dirPath string, dialectTableName string) (string, string) {
	const layout = "01-02-2006"
	t := time.Now()
	filenameFr := dirPath + "/" + fmt.Sprintf("%s_fr_%s.txt", dialectTableName, t.Format(layout))
	filenameEn := dirPath + "/" + fmt.Sprintf("%s_en_%s.txt", dialectTableName, t.Format(layout))
	return filenameFr, filenameEn
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
