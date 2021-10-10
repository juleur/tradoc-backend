package controllers

import (
	"btradoc/entities"
	"btradoc/helpers"
	"btradoc/pkg"
	"btradoc/pkg/dialect"
	"btradoc/pkg/email"
	"btradoc/pkg/translator"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// TODO: refactor this file

func Login(secretKey string, translatorService translator.Service, dialectService dialect.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*logrus.Logger)

		loginBody := new(entities.LoginBody)
		if err := c.BodyParser(&loginBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		// TODO: rework this part
		transl, err := translatorService.FindTranslatorByUsername(strings.TrimSpace(loginBody.Username))
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

		if transl == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cet utilisateur n'existe pas",
			})
		}

		if !transl.Confirmed {
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

		// for frontend UI
		if len(transl.Permissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Aucune dialect ne t'a été attribué, contacte l'administrateur",
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://occitanofon.fr"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.fr"
		claims["sub"] = transl.ID
		claims["dperms"] = transl.CompressPerms()

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
		logger := c.Locals("logger").(*logrus.Logger)

		registerBody := new(entities.RegisterBody)
		if err := c.BodyParser(&registerBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrDefault,
			})
		}

		registerBody.Email = strings.TrimSpace(registerBody.Email)
		registerBody.Username = strings.TrimSpace(registerBody.Username)
		registerBody.Password = strings.TrimSpace(registerBody.Password)

		// check email validity
		if !helpers.IsEmailValid(registerBody.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		// check username validity
		if err := helpers.UsernameValidity(registerBody.Username); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// check password valitidy
		if len(registerBody.Password) < 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ton mot de passe doit contenir 10 caracters",
			})
		}

		// check if there is 2 secret questions
		if len(registerBody.SecretQuestions) != 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "il n'y a pas de 2 secret questions",
			})
		}

		// hash password
		hpwd, err := argon2id.CreateHash(registerBody.Password, argon2id.DefaultParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ISSOU",
			})
		}

		// hash responses from secret question for more privacy
		var sQHashResponses []entities.SecretQuestion
		for _, secretQuestion := range registerBody.SecretQuestions {
			// cleaning
			secretQuestion.Question = strings.TrimSpace(secretQuestion.Question)
			secretQuestion.Response = strings.TrimSpace(secretQuestion.Response)
			secretQuestion.Response = strings.ToLower(secretQuestion.Response)

			var sq entities.SecretQuestion
			sq.Question = secretQuestion.Question
			sQhashResp, err := argon2id.CreateHash(secretQuestion.Response, argon2id.DefaultParams)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "ISSOU",
				})
			}
			sq.Response = sQhashResp

			sQHashResponses = append(sQHashResponses, sq)
		}

		newTranslator := entities.NewTranslator{
			Email:           registerBody.Email,
			Username:        registerBody.Username,
			Hpwd:            hpwd,
			SecretQuestions: sQHashResponses,
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

		// besoin pour le frontend
		if len(transl.Permissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Aucune permission ne vous a été attribué",
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = "https://occitanofon.fr"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		claims["iss"] = "https://api.occitanofon.fr"
		claims["sub"] = transl.ID
		claims["dperms"] = transl.CompressPerms()

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
		logger := c.Locals("logger").(*logrus.Logger)

		refreshToken := c.Cookies("refreshToken")

		if err := translatorService.DeleteRefreshToken(refreshToken); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		// c.Context().Response.Header.DelClientCookie("refreshToken")
		c.ClearCookie("refreshToken")

		return c.SendStatus(fiber.StatusOK)
	}
}

func SendPasswordReset(translatorService translator.Service, emailService email.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		emailBody := new(entities.ResetPasswordBody)
		if err := c.BodyParser(&emailBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		transl, err := translatorService.ProceedResetPassword(emailBody.Email)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		emailService.SendResetPasswordLink(transl)

		return c.JSON(fiber.Map{
			"msg": "Un email vous sera prochainement envoyé afin de procéder à la réinitialisation de votre mot de passe",
		})
	}
}

func ConfirmPasswordResetToken(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		token := c.Params("token")
		translatorSecretQuestions, err := translatorService.FetchSecretQuestionsByToken(token)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		} else if len(translatorSecretQuestions.SecretQuestions) == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// pull out responses
		var secretQuestionsOnly []string
		for question, _ := range translatorSecretQuestions.SecretQuestions {
			secretQuestionsOnly = append(secretQuestionsOnly, question)
		}

		return c.JSON(secretQuestionsOnly)
	}
}

func SecretQuestions(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		secretQuestions, err := translatorService.FetchSecretQuestions()
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		} else if len(secretQuestions) == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.JSON(secretQuestions)
	}
}

func UpdatePassword(translatorService translator.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logger, _ := c.Locals("logger").(*logrus.Logger)

		updatePasswordBody := new(entities.UpdatePasswordBody)
		if err := c.BodyParser(&updatePasswordBody); err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		}

		// cleaning body
		for _, secret := range updatePasswordBody.SecretQuestions {
			secret.Question = strings.TrimSpace(secret.Question)
			secret.Response = strings.TrimSpace(secret.Response)
		}

		translatorSecretQuestions, err := translatorService.FetchSecretQuestionsByToken(updatePasswordBody.Token)
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "A REMPLIR",
			})
		} else if len(translatorSecretQuestions.SecretQuestions) == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}

		for _, secret := range updatePasswordBody.SecretQuestions {
			for questionStored, responseStored := range translatorSecretQuestions.SecretQuestions {
				if secret.Question == questionStored {
					match, _, err := argon2id.CheckHash(secret.Response, responseStored)

					if !match {
						return c.Status(fiber.StatusConflict).JSON(fiber.Map{
							"error": "une des reponse ne match pas",
						})
					} else if err != nil {
						logger.Error(err)
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "A REMPLIR",
						})
					}
				}
			}
		}

		// hash password
		hpwd, err := argon2id.CreateHash(strings.TrimSpace(updatePasswordBody.Password), argon2id.DefaultParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ISSOU",
			})
		}

		if err = translatorService.ResetPassword(translatorSecretQuestions.ID.Hex(), hpwd); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ISSOU",
			})
		}

		return c.SendStatus(fiber.StatusAccepted)
	}
}
