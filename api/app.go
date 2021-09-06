package main

import (
	"btradoc/api/middlewares"
	"btradoc/api/routes"
	"btradoc/entities"
	"btradoc/pkg/dialect"
	"btradoc/pkg/translation"
	"btradoc/pkg/translator"
	"btradoc/tools"
	"crypto/rsa"
	"flag"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

const SECRET_KEY = "WE1Lb9XqN0P0REmFWLSccKNxjGikZzAECereA5bH17dxFX0rIx4DGIbHC4NUfUwu"

var (
	TRANSLATIONS_FILES_PATH string
	PRODUCTION_MOD          *bool
	LOG_FILE                string
	ALLOW_ORIGINS           string
	PRIVATE_KEY             *rsa.PrivateKey
)

func init() {
	db := tools.ArangoDBConnection()
	dialects := tools.OpenDialectsJSONFile()

	if err := tools.CreateCollection(db, entities.COLLECTIONS); err == nil {
		tools.AddDialectsDocuments(db, dialects)
	}
	////////////

	PRODUCTION_MOD = flag.Bool("prod", false, "a string")
	flag.Parse()

	if *PRODUCTION_MOD {
		ALLOW_ORIGINS = "https://trad-oc.fr, https://trad-oc.fr, https://trad-oc.fr"
	} else {
		ALLOW_ORIGINS = "http://127.0.0.1:3333"
	}
	LOG_FILE = "logrus.log"
}

func main() {
	logger := logrus.New()
	logger.ReportCaller = true

	db := tools.ArangoDBConnection()

	translatorRepo := translator.NewRepo(db)
	translatorService := translator.NewService(translatorRepo)

	dialectRepo := dialect.NewRepo(db)
	dialectService := dialect.NewService(dialectRepo)

	translationRepo := translation.NewRepo(db)
	translationService := translation.NewService(translationRepo)

	app := fiber.New()

	if *PRODUCTION_MOD {
		file := tools.CreateLogFile(LOG_FILE)
		defer file.Close()

		logger.SetOutput(file)
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		app.Use(log.New())

		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     ALLOW_ORIGINS,
		AllowMethods:     "GET,POST,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Cookie, Authorization, Set-Cookie",
		AllowCredentials: true,
	}))

	app.Use(middlewares.Logrus(logger))

	routes.PrivateEndpoints(app, SECRET_KEY, translatorService, dialectService, translationService)
	routes.PublicEndpoints(app, SECRET_KEY, translatorService, dialectService)

	_ = app.Listen(":9321")
}
