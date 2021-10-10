package routes

import (
	"btradoc/api/controllers"
	"btradoc/pkg/dialect"
	"btradoc/pkg/email"
	"btradoc/pkg/translation"
	"btradoc/pkg/translator"

	"github.com/gofiber/fiber/v2"
)

func PublicEndpoints(app fiber.Router, secretKey string, translatorService translator.Service, dialectService dialect.Service, translationService translation.Service, emailService email.Service) {
	app.Post("/login", controllers.Login(secretKey, translatorService, dialectService))
	app.Get("/refreshtoken", controllers.RefreshToken(secretKey, dialectService, translatorService))
	app.Post("/register", controllers.Register(translatorService))
	app.Post("/send_pwd_reset", controllers.SendPasswordReset(translatorService, emailService))
	app.Get("/secret_questions", controllers.SecretQuestions(translatorService))
	app.Get("/confirm_token/:token", controllers.ConfirmPasswordResetToken(translatorService))
	app.Post("/update_pwd", controllers.UpdatePassword(translatorService))
	app.Get("/translations_files", controllers.TranslationsFiles(translationService))
}
