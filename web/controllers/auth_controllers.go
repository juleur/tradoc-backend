package controllers

import (
	"net/http"
	"strconv"
	"tradoc/db"
	"tradoc/message"
	"tradoc/models"
	"tradoc/pkg/password"
	"tradoc/pkg/store"
	"tradoc/pkg/token"

	"github.com/phuslu/log"

	"github.com/gofiber/fiber/v2"
)

func login(db *db.DBPsql, onlineTranslatorsStore *store.OnlineTranslatorsStore) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loginBody := new(models.LoginBody)
		if err := c.BodyParser(&loginBody); err != nil {
			c.Locals("logger").(*log.Logger).Error().Interface("loginBody", loginBody).Msg(message.ErrBadBodyContent.Error())
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseContentBody,
			})
		}

		translator, err := db.FindTranslator(*loginBody)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("username", loginBody.Username).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		err = password.ComparePasswords(loginBody.Password, translator.HashedPassword)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Str("password", loginBody.Password).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		translator.Dialects, err = db.FetchDialectsByTranslator(translator.ID)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("translatorID", int64(translator.ID)).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		SECRET_KEY := c.Locals("SECRET_KEY").([]byte)

		tokens, err := token.GenerateTokens(SECRET_KEY, translator)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Interface("tokens", tokens).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		go func() {
			if err := db.UpdateRefreshToken(translator.ID, tokens.RefreshToken); err != nil {
				c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("translatorID", int64(translator.ID)).Str("refreshToken", tokens.RefreshToken).Msg(err.Error.Error())
				return
			}
		}()

		onlineTranslatorsStore.Put(translator.Username)

		c.Locals("logger").(*log.Logger).Info().Msgf("%s has successfully logged in", loginBody.Username)
		return c.JSON(tokens)
	}
}

func refreshToken(db *db.DBPsql, onlineTranslatorsStore *store.OnlineTranslatorsStore) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		SECRET_KEY := c.Locals("SECRET_KEY").([]byte)

		refToken := c.Context().Request.Header.Peek("X-Refresh-Token")
		xUserID := c.Context().Request.Header.Peek("X-User-ID")

		if len(refToken) != 32 || !token.IsAlphanumeric(refToken) {
			c.Locals("logger").(*log.Logger).Error().Str("X-Refresh-Token", string(refToken)).Msg(message.ErrRefreshTokenBadFormat.Error())
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"errorCode": 40,
				"message":   message.ResponseAuthFailed,
			})
		}

		userID, er := strconv.Atoi(string(xUserID))
		if er != nil {
			c.Locals("logger").(*log.Logger).Error().Str("X-User-ID", string(xUserID)).Msg(message.ErrRefreshTokenBadFormat.Error())
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"errorCode": 40,
				"message":   message.ResponseAuthFailed,
			})
		}

		translator, err := db.IsRefreshTokenExists(userID, string(refToken))
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("userID", int64(userID)).Str("refreshToken", string(refToken)).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 40,
				"message":   err.Message,
			})
		}

		translator.Dialects, err = db.FetchDialectsByTranslator(translator.ID)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("translatorID", int64(translator.ID)).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		newTokens, err := token.GenerateTokens(SECRET_KEY, translator)
		if err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Interface("newTokens", newTokens).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		go func() {
			if err := db.UpdateRefreshToken(translator.ID, newTokens.RefreshToken); err != nil {
				c.Locals("logger").(*log.Logger).WithLevel(err.Level).StrInt("translatorID", int64(translator.ID)).Str("refreshToken", newTokens.RefreshToken).Msg(err.Error.Error())
				return
			}
		}()

		onlineTranslatorsStore.Put(translator.Username)

		c.Locals("logger").(*log.Logger).Info().Msgf("%s has successfully refreshed his tokens", translator.Username)
		return c.JSON(newTokens)
	}
}

func register(db *db.DBPsql) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		validateAccountBody := new(models.ValidateAccountBody)
		if err := c.BodyParser(&validateAccountBody); err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(log.ErrorLevel).Interface("validateAccountBody", validateAccountBody).Msg(message.ErrBadBodyContent.Error())
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseContentBody,
			})
		}

		if validateAccountBody.Code < 999 || validateAccountBody.Code > 9999 {
			c.Locals("logger").(*log.Logger).WithLevel(log.ErrorLevel).StrInt("validateAccountBody.Code", int64(validateAccountBody.Code)).Msg("")
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   message.ResponseContentBody,
			})
		}

		translatorID, er := db.FirstRegister(validateAccountBody.Username, validateAccountBody.Code)
		if er != nil {
			c.Locals("logger").(*log.Logger).WithLevel(er.Level).Str("username", validateAccountBody.Username).StrInt("code", int64(validateAccountBody.Code)).Msg(er.Error.Error())
			return c.Status(er.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   er.Message,
			})
		}

		hashedPWd, er := password.HashPassword(validateAccountBody.Password)
		if er != nil {
			c.Locals("logger").(*log.Logger).WithLevel(er.Level).Str("username", validateAccountBody.Username).StrInt("code", int64(validateAccountBody.Code)).Msg(er.Error.Error())
			return c.Status(er.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   er.Message,
			})
		}

		if err := db.CreateNewTranslator(translatorID, validateAccountBody.Username, hashedPWd); err != nil {
			c.Locals("logger").(*log.Logger).WithLevel(err.Level).Interface("validateAccountBody", validateAccountBody).Msg(err.Error.Error())
			return c.Status(err.HttpStatusCode).JSON(&fiber.Map{
				"errorCode": 1,
				"message":   err.Message,
			})
		}

		c.Locals("logger").(*log.Logger).Info().Msgf("%s has successfully validated his account", validateAccountBody.Username)
		return c.SendStatus(http.StatusAccepted)
	}
}

func MakeAuthControllers(app *fiber.App, dbRepo *db.DBPsql, onlineTranslatorsStore *store.OnlineTranslatorsStore) {
	app.Post("/login", login(dbRepo, onlineTranslatorsStore))
	app.Get("/refresh_token", refreshToken(dbRepo, onlineTranslatorsStore))
	app.Post("/register", register(dbRepo))
}
