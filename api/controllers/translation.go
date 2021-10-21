package controllers

import (
	"btradoc/entities"
	"btradoc/pkg"
	"btradoc/pkg/activity"
	"btradoc/pkg/dataset"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translation"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewTranslations(datasetService dataset.Service, dialectService dialect.Service, translationService translation.Service, activityService activity.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		// translationsBody has original translation and might have feminize translation
		translationsBody := new([]entities.Translation)
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

		for _, translation := range *translationsBody {
			if !strings.Contains(translation.Occitan, "_") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": ErrBadFullDialectFormat,
				})
			}

			occitan := strings.Split(translation.Occitan, "_")
			dialect := occitan[0]
			subdialect := occitan[1]
			// verify dialect and subdialect name
			match, err := dialectService.Exists(dialect, subdialect)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrDefault,
				})
			} else if !match {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": pkg.ErrDialectNotFound,
				})
			}
		}

		translations := *translationsBody

		if err := translationService.AddTranslations(translatorID, translations); err != nil {
			logger.Error(err)
			if e, ok := err.(*pkg.DBError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		datasetID := translations[0].DatasetID
		fullDialect := translations[0].Occitan
		if err := datasetService.AddTranslatedIn(datasetID, fullDialect); err != nil {
			logger.Error(err)
			if e, ok := err.(*pkg.DBError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		if err := translationService.RemoveOnGoingTranslations(translations); err != nil {
			logger.Error(err)
			if e, ok := err.(*pkg.DBError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		// track translator activities
		activityService.AddOrKeepActive(translatorID)

		return c.SendStatus(fiber.StatusOK)
	}
}

func TranslationsFiles(translationService translation.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translationsFiles, err := translationService.FetchPathnameFiles()
		if err != nil {
			logger.Error(err)
			if e, ok := err.(*pkg.DBError); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		return c.JSON(translationsFiles)
	}
}

func TotalOnlineTranslators(activityService activity.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		totalOnlineTransl := activityService.Total()

		return c.JSON(totalOnlineTransl)
	}
}
