package controllers

import (
	"btradoc/pkg"
	"btradoc/pkg/dialect"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func FetchDialectsWithSubdialects(dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		translatorID := c.Locals("translatorID").(string)

		dialectsSubdialects, err := dialectService.FetchDialectsSubdialect(translatorID)
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

		return c.JSON(dialectsSubdialects)
	}
}
