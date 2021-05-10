package tasks

import (
	"tradoc/pkg/store"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
)

func RunAllTasks(cron *cron.Cron, logger *log.Logger, database *pgxpool.Pool, onGoingTranslationsStore *store.OnGoingTranslationsStore, onlineTranslatorsStore *store.OnlineTranslatorsStore, dialectsTableName []string, TRANSLATIONS_PATH string) {
	// supprime les traducteurs qui ne sont plus en ligne
	cron.AddFunc("*/3 * * * *", func() {
		onlineTranslatorsStore.CleanOnlineTranslators()
	})

	// supprime les traductions en cours
	cron.AddFunc("*/10 * * * *", func() {
		translationsRemoved := onGoingTranslationsStore.CleanOnGoingTranslations()

		var dialectsIDs []string
		for translationRemoved := range translationsRemoved {
			dialectsIDs = append(dialectsIDs, translationRemoved)
		}
		if len(dialectsIDs) == 0 {
			logger.Info().Str("task", "CleanOnGoingTranslations").Msg("No translation has been removed from OnGoingTranslationsStore")
		} else {
			logger.Info().Str("task", "CleanOnGoingTranslations").Msgf("Translations n°%v have been removed from OnGoingTranslationsStore", dialectsIDs)
		}
	})

	// rotate log files
	cron.AddFunc("0 0 * * *", func() {
		if err := logger.Writer.(*log.FileWriter).Rotate(); err != nil {
			logger.Error().Str("task", "Rotating log").Msg(err.Error())
			return
		}
	})

	// supprime les anciens fichiers .txt de traductions
	cron.AddFunc("0 4 * * 1", func() {
		resultErrors := CleanOldTranslations(database, TRANSLATIONS_PATH)
		for resultErr := range resultErrors {
			logger.Error().Str("task", "CleanOldTranslations").Msg(resultErr.Error())
		}
	})

	// génére automatiquement les traductions
	cron.AddFunc("30 4 * * 1", func() {
		resultErrors := GenerateTranslatedSentences(database, dialectsTableName, TRANSLATIONS_PATH)
		for resultErr := range resultErrors {
			logger.Error().Str("task", "GenerateTranslatedSentences").Msg(resultErr.Error())
		}
	})
}
