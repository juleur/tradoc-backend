package routes

import (
	"btradoc/api/controllers"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translator"

	"github.com/gofiber/fiber/v2"
)

func PublicEndpoints(app fiber.Router, secretKey string, translatorService translator.Service, dialectService dialect.Service) {
	app.Post("/login", controllers.Login(secretKey, translatorService, dialectService))
	app.Get("/refreshtoken", controllers.RefreshToken(secretKey, dialectService, translatorService))
	app.Post("/register", controllers.Register(translatorService))
	app.Post("/pwd_reset", controllers.PasswordReset(translatorService))
	// files translations
}
