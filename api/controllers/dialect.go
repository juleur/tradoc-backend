package controllers

import (
	"btradoc/pkg/dialect"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func FetchDialectsWithSubdialects(dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		translatorID := "translators/12504"

		dialectsSubdialects, err := dialectService.FetchDialectsSubdialect(translatorID)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		return c.JSON(dialectsSubdialects)
	}
}
