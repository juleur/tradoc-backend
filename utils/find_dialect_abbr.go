package utils

import (
	"strings"
	"tradoc/message"
	"tradoc/models"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func FindDialectByAbbre(dialectsTableName []string, dialectAbbr string) (string, *models.Log) {
	dialectAbbrSplit := strings.Split(dialectAbbr, "_")
	for _, dialect := range dialectsTableName {
		dialectSplit := strings.Split(dialect, "_")
		if dialectAbbrSplit[0][0:3] == dialectSplit[0][0:3] && dialectAbbrSplit[1][0:3] == dialectSplit[1][0:3] {
			return dialect, nil
		}
	}
	return "", &models.Log{
		Level:          log.WarnLevel,
		Error:          message.ErrNoDialectAbbrFound,
		HttpStatusCode: fiber.ErrNotFound.Code,
		Message:        message.ResponseErrServer,
	}
}
