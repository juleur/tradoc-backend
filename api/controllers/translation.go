package controllers

import (
	"btradoc/pkg/translation"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDatasets(translationService translation.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		fullDialect := c.Params("full_dialect")
		if len(fullDialect) == 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Aucun sous dialect fourni",
			})
		}

		decodedFullDialect, err := url.QueryUnescape(fullDialect)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Aucun sous dialect fourni",
			})
		}

		fullDialectSplit := strings.Split(decodedFullDialect, "-")
		datasets, err := translationService.FetchSentencesToTranslate(translatorID, fullDialectSplit[0], fullDialectSplit[1])
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Aucun sous dialect fourni",
			})
		}

		fmt.Println(*datasets)
		return c.JSON(*datasets)
	}
}

func AddNewTranslations() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}
