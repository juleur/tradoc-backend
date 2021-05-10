package middlewares

import (
	"bytes"
	"net/http"
	"tradoc/message"
	"tradoc/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func Authorization() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		SECRET_KEY := c.Locals("SECRET_KEY").([]byte)

		bearerJwt := c.Context().Request.Header.Peek("Authorization")

		if token.IsItAJwtToken(bearerJwt) {
			jwToken := bytes.Split(bearerJwt, []byte(" "))[1]

			translatorID, expired := token.GetJwtPayload(SECRET_KEY, jwToken)
			if expired {
				c.Locals("logger").(*log.Logger).Warn().StrInt("translatorID", int64(translatorID)).Msg(message.ErrJWTExpired.Error())
				return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
					"errorCode": 41,
					"message":   message.ErrJWTExpired.Error(),
				})
			} else if translatorID == 0 {
				c.Locals("logger").(*log.Logger).Error().Str("jwt", string(jwToken[1])).Msg(message.ErrJWTPayload.Error())
				return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
					"errorCode": 40,
					"message":   message.ResponseAuthFailed,
				})
			}
			////// gérer le cas lorsque le token est expiré ou invalid
			c.Locals("translatorID", translatorID)
			return c.Next()
		}
		c.Locals("logger").(*log.Logger).Error().Str("BearerJwt", string(bearerJwt)).Msg(message.ErrWrongBearerJWT.Error())
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"errorCode": 40,
			"message":   message.ResponseAuthFailed,
		})
	}
}
