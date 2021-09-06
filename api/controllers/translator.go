package controllers

import (
	"btradoc/entities"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translator"
	"btradoc/tools"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func Login(secretKey string, translatorService translator.Service, dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		loginBody := new(entities.LoginBody)
		if err := c.BodyParser(&loginBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		transl, err := translatorService.FindTranslatorByUsername(strings.TrimSpace(loginBody.Username))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}
		if transl == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cet utilisateur n'existe pas",
			})
		}

		if !transl.Activated {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Votre compte a besoin d'être activer afin de pouvoir vous connecter",
			})
		}

		if transl.Suspended {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Votre compte est suspendu, contactez l'administrateur",
			})
		}

		match, err := argon2id.ComparePasswordAndHash(loginBody.Password, transl.Hpwd)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		} else if !match {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		dialectPermissions, err := dialectService.FindDialectPermissions(transl.ID)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		// for frontend UI
		if len(*dialectPermissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Aucune dialect ne t'a été attribué, contacte l'administrateur",
			})
		}

		dperms := tools.MakeDialectPermissionsIntoAbbreviations(dialectPermissions)

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://occitanofon.fr"
		claims["exp"] = time.Now().Add(time.Second * 30).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.fr"
		claims["sub"] = transl.Key
		claims["dperms"] = dperms

		accessToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error SET accessToken",
			})
		}

		refreshToken, err := translatorService.SetRefreshToken(transl.ID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "error SET refresh token",
			})
		}

		c.Set("Set-Cookie", fmt.Sprintf("refreshToken=%s;HttpOnly;Path=/;SameSite=Lax", refreshToken))

		return c.JSON(accessToken)
	}
}

func Register(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		registerBody := new(entities.RegisterBody)
		if err := c.BodyParser(&registerBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		registerBody.Email = strings.TrimSpace(registerBody.Email)
		registerBody.Username = strings.TrimSpace(registerBody.Username)
		registerBody.Password = strings.TrimSpace(registerBody.Password)

		if !tools.IsEmailValid(registerBody.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		if err := tools.UsernameValidity(registerBody.Username); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		hpwd, err := argon2id.CreateHash(strings.TrimSpace(registerBody.Password), argon2id.DefaultParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ISSOU",
			})
		}

		newTranslator := entities.NewTranslator{
			Email:     registerBody.Email,
			Username:  registerBody.Username,
			Hpwd:      hpwd,
			CreatedAt: time.Now(),
		}

		if err = translatorService.CreateTranslator(newTranslator); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func RefreshToken(secretKey string, dialectService dialect.Service, translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		refreshToken := c.Cookies("refreshToken")

		transl, err := translatorService.FindRefreshToken(refreshToken)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		dialectPermissions, err := dialectService.FindDialectPermissions(transl.ID)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		// besoin pour le frontend
		if len(*dialectPermissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Aucune permission ne vous a été attribué ",
			})
		}

		dperms := tools.MakeDialectPermissionsIntoAbbreviations(dialectPermissions)

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://occitanofon.fr"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.fr"
		claims["sub"] = transl.Key
		claims["dperms"] = dperms

		accessToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error SET accessToken",
			})
		}

		c.Set("Set-Cookie", fmt.Sprintf("refreshToken=%s;HttpOnly;Path=/;SameSite=Lax", refreshToken))

		return c.JSON(accessToken)
	}
}

func Logout(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		refreshToken := c.Cookies("refreshToken")

		if err := translatorService.DeleteRefreshToken(refreshToken); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		c.Context().Response.Header.DelClientCookie("refreshToken")

		return c.SendStatus(fiber.StatusOK)
	}
}

func PasswordReset(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// logger, _ := c.Locals("logger").(*logrus.Logger)

		//
		return c.SendStatus(fiber.StatusAccepted)
	}
}
