package controllers

import (
	"net/http"
	"tradoc/db"
	"tradoc/message"
	"tradoc/models"
	"tradoc/pkg/store"
	"tradoc/web/middlewares"

	"tradoc/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func fetchDialectes(db *db.DBPsql, onlineTranslatorsStore *store.OnlineTranslatorsStore, DIALECTS []models.Dialect, dialectsTableName []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		translatorID := c.Locals("translatorID").(int)

		for i, dialect := range DIALECTS {
			for y, sub := range dialect.Subdialects {
				dialectTableName, err := utils.FindDialectByAbbre(dialectsTableName, sub.Abbr)
				if err != nil {
					c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("translatorID", int64(translatorID)).Str("subdialect", sub.Abbr).Msg(err.Error.Error())
					continue
				}
				total, err := db.FetchTotalSentenceTranslatedByDialect(dialectTableName)
				if err != nil {
					c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialectTableName", dialectTableName).Msg(err.Error.Error())
					continue
				}
				DIALECTS[i].Subdialects[y].TotalSentencesTranslated = total
			}
		}

		mainMenu := models.MainMenu{
			Dialects:               DIALECTS,
			TotalOnlineTranslators: onlineTranslatorsStore.TotalOnlineTranslators(),
		}

		return c.JSON(mainMenu)
	}
}

func newTranslation(db *db.DBPsql, onGoingTranslationsStore *store.OnGoingTranslationsStore, dialectsTableName []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		translatorID := c.Locals("translatorID").(int)

		newTranslationBody := new(models.NewTranslationBody)
		if err := c.BodyParser(newTranslationBody); err != nil {
			c.Locals("logger").(*log.Logger).Error().Interface("body", newTranslationBody).Msg(message.ErrBadBodyContent.Error())
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseContentBody,
			})
		}

		if newTranslationBody.Translation.Abbr == "" {
			c.Locals("logger").(*log.Logger).Error().Msg(message.ErrNoDialectAbbrProvided.Error())
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseDialectNotProvided,
			})
		}

		dialectTableName, err := utils.FindDialectByAbbre(dialectsTableName, newTranslationBody.Translation.Abbr)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialect abbre", newTranslationBody.Translation.Abbr).Msg(message.ErrNoDialectAbbrFound.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		err = db.AddTranslation(*&newTranslationBody.Translation, dialectTableName, translatorID)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialectTableName", dialectTableName).StrInt("translatorID", int64(translatorID)).Interface("newTranslationBody", newTranslationBody.Translation).Msg(message.ErrAddTranslation.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		go func() {
			onGoingTranslationsStore.Delete(dialectTableName, newTranslationBody.Translation.ID)
		}()

		return c.SendStatus(http.StatusAccepted)
	}
}

func fetchSentencesByDialect(db *db.DBPsql, onGoingTranslationsStore *store.OnGoingTranslationsStore, dialectsTableName []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dialectParam := c.Params("dialect")

		translatorID := c.Locals("translatorID").(int)

		dialectTableName, err := utils.FindDialectByAbbre(dialectsTableName, dialectParam)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialect abbr", dialectParam).Msg(message.ErrNoDialectAbbrFound.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseErrServer,
			})
		}

		// gérer les message erreur
		if onGoingTranslationsStore.LimitRate(dialectTableName, translatorID) {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialectTable", dialectTableName).StrInt("translatorID", int64(translatorID)).Msg(message.ErrNoSentencesFound.Error())
			return c.Status(fiber.ErrForbidden.Code).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		ruleOutTextesIDs := onGoingTranslationsStore.GetTexteIDs(dialectTableName)

		textes, err := db.FetchAllSentencesByDialect(dialectTableName, translatorID, ruleOutTextesIDs)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialectTable", dialectTableName).StrInt("translatorID", int64(translatorID)).Msg(message.ErrNoSentencesFound.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		var texteIDs []int
		for _, texte := range textes {
			texteIDs = append(texteIDs, texte.ID)
		}
		onGoingTranslationsStore.Put(dialectTableName, translatorID, texteIDs)

		return c.JSON(textes)
	}
}

func dialectFiles(db *db.DBPsql, DIALECTS []models.Dialect, dialectsTableName []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		allFiles := models.AllFiles{}
		for _, dialect := range DIALECTS {
			dialectFile := models.DialectFile{
				Nom: dialect.Nom,
			}
			for _, sub := range dialect.Subdialects {
				subdialectFile := models.SubdialectFile{
					Nom: sub.Nom,
				}
				dialectTableName, err := utils.FindDialectByAbbre(dialectsTableName, sub.Abbr)
				if err != nil {
					c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("subdialect", sub.Abbr).Msg(err.Error.Error())
				}
				filepathFr, filepathEn, err := db.FetchTranslatedFiles(dialectTableName)
				if err != nil {
					c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("dialectTable", dialectTableName).Msg(err.Error.Error())
				}
				subdialectFile.FilepathFr = filepathFr
				subdialectFile.FilepathEn = filepathEn

				dialectFile.SubdialectFiles = append(dialectFile.SubdialectFiles, subdialectFile)
			}
			allFiles.DialectFiles = append(allFiles.DialectFiles, dialectFile)
		}
		allFiles.LastDatetimeGenFile = db.FetchLastGeneratedFile()

		return c.JSON(allFiles)
	}
}

func MakeDialectControllers(app *fiber.App, dbRepo *db.DBPsql, onGoingTranslationsStore *store.OnGoingTranslationsStore, onlineTranslatorsStore *store.OnlineTranslatorsStore, DIALECTS []models.Dialect, dialectsTableName []string) {
	app.Get("/dialect_files", dialectFiles(dbRepo, DIALECTS, dialectsTableName))
	api := app.Group("/auth", middlewares.Authorization())
	api.Get("/dialects", fetchDialectes(dbRepo, onlineTranslatorsStore, DIALECTS, dialectsTableName))
	api.Post("/new_translation", newTranslation(dbRepo, onGoingTranslationsStore, dialectsTableName))
	api.Get("/sentences/:dialect", fetchSentencesByDialect(dbRepo, onGoingTranslationsStore, dialectsTableName))
}
