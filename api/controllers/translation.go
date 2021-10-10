package controllers

import (
	"btradoc/entities"
	"btradoc/pkg"
	"btradoc/pkg/translation"
	"btradoc/storage/inmemory"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDatasets(translationService translation.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		fullDialectParam := c.Params("full_dialect")
		if len(fullDialectParam) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrDialectNotProvided,
			})
		}

		fullDialect, err := url.QueryUnescape(fullDialectParam)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		totalOnGoingTranslations, err := translationService.FetchTotalOnGoingTranslations(fullDialect, translatorID)
		if err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		// prevent translator to reload then fetch again and again without limit
		if totalOnGoingTranslations > 300 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": ErrTooMuchTranslationsFetched,
			})
		}

		datasets, err := translationService.FetchDatasets(fullDialect)
		if err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		// assure that there is still dataset to translate
		if len(*datasets) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrNoMoreDataset,
			})
		}

		if err = translationService.AddOnGoingTranslations(fullDialect, translatorID, datasets); err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		return c.JSON(*datasets)
	}
}

func AddNewTranslations(translationService translation.Service, activeTranslatorsTracker *inmemory.ActiveTranslators) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		translationsBody := new(entities.TranslationsBody)
		if err := c.BodyParser(&translationsBody); err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		// add translatorID to each translations
		for _, translation := range translationsBody.Translations {
			translation.TranslatorID = translatorID
		}

		if err := translationService.AddTranslations(translationsBody.Translations); err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		if err := translationService.AddDatasetNewFullDialect(translationsBody.Translations); err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		if err := translationService.RemoveOnGoingTranslations(translationsBody.Translations); err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		// track translator activities
		activeTranslatorsTracker.AddOrKeepActive(translatorID)

		return c.SendStatus(fiber.StatusOK)
	}
}

func TranslationsFiles(translationService translation.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translationsFiles, err := translationService.FetchTranslationsFiles()
		if err != nil {
			logger.Error(err)
			switch err.(type) {
			case pkg.DBError:
				return c.Status(err.(*pkg.DBError).Code).JSON(fiber.Map{
					"error": err.(*pkg.DBError).Message,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			}
		}

		return c.JSON(translationsFiles)
	}
}
