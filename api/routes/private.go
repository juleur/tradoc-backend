package routes

import (
	"btradoc/api/controllers"
	"btradoc/api/middlewares"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translation"
	"btradoc/pkg/translator"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func PrivateEndpoints(app fiber.Router, secretKey string, translatorService translator.Service, dialectService dialect.Service, translationService translation.Service) {
	api := app.Group("/p", jwtware.New(jwtware.Config{
		SigningMethod:  jwt.SigningMethodHS256.Name,
		SigningKey:     []byte(secretKey),
		SuccessHandler: middlewares.JWTSuccessHandler(),
		ErrorHandler:   middlewares.JWTErrorHandler(secretKey),
	}))

	api.Add("GET", "/dialects_subdialects", controllers.FetchDialectsWithSubdialects(dialectService))
	api.Add("GET", "/get_datasets/:full_dialect", controllers.GetDatasets(translationService))
	// api.Add("POST", "/add_new_translations", controllers.AddNewTranslations(s))
	api.Add("GET", "/logout", controllers.Logout(translatorService))
}
