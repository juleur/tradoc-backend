package password

import (
	"tradoc/message"
	"tradoc/models"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func ComparePasswords(password string, hashedPassword string) *models.Log {
	if password == "" || hashedPassword == "" {
		return &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrNoPasswords,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if !match {
		return &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrPasswordNoMatch,
			HttpStatusCode: fiber.ErrForbidden.Code,
			Message:        message.ResponseBadPassword,
		}
	} else if err != nil {
		return &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrNoPasswords,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}

	return nil
}

func HashPassword(password string) (string, *models.Log) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", &models.Log{
			Level:          log.ErrorLevel,
			Error:          message.ErrHashedPasswordFailed,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	return hashedPassword, nil
}
