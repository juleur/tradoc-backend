package tasks

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CleanOldTranslations(db *pgxpool.Pool, dirPath string) <-chan error {
	errCh := make(chan error, 4)
	go func() {
		_ = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, _ error) error {
			fileInfo, err := d.Info()
			if err != nil {
				errCh <- err
				return nil
			}
			if !fileInfo.IsDir() {
				after2Month := time.Now().AddDate(0, -1, 0)
				if !after2Month.Before(fileInfo.ModTime()) {
					if _, err := db.Exec(context.Background(), `
						DELETE FROM translation_files WHERE filename_fr = $1 OR filename_en = $1
					`, fileInfo.Name()); err != nil {
						errCh <- err
						return nil
					}
					os.Remove(dirPath + "/" + fileInfo.Name())
				}
			}
			return nil
		})
		close(errCh)
	}()

	return errCh
}
