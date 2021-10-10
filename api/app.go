package main

import (
	"btradoc/api/middlewares"
	"btradoc/api/routes"
	"btradoc/helpers"
	"btradoc/pkg/dialect"
	"btradoc/pkg/email"
	"btradoc/pkg/translation"
	"btradoc/pkg/translator"
	"btradoc/storage/inmemory"
	"btradoc/storage/mongodb"
	"flag"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

//        /!\ NOT SAFE /!\        //
const SECRET_KEY string = "WE1Lb9XqN0P0REmFWLSccKNxjGikZzAECereA5bH17dxFX0rIx4DGIbHC4NUfUwu"

var (
	TRANSLATIONS_FILES_PATH string
	PRODUCTION_MOD          *bool
	LOG_FILE                string
	ALLOW_ORIGINS           string
)

func init() {
	LOG_FILE = "logrus.log"

	PRODUCTION_MOD = flag.Bool("prod", false, "a string")
	flag.Parse()

	if *PRODUCTION_MOD {
		ALLOW_ORIGINS = "https://occitanofon.xyz, https://occitanofon.xyz"
	} else {
		ALLOW_ORIGINS = "http://127.0.0.1:3333"
	}

	mongodb.InitMongoDatabase()
}

func main() {
	logger := logrus.New()
	logger.ReportCaller = true

	if *PRODUCTION_MOD {
		file := helpers.CreateLogFile(LOG_FILE)
		defer file.Close()

		logger.SetOutput(file)
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	db := mongodb.NewMongoClient()
	activeTranslatorsTracker := inmemory.NewActiveTranslatorsTracker()

	translatorRepo := translator.NewRepo(db)
	translatorService := translator.NewService(translatorRepo)

	dialectRepo := dialect.NewRepo(db)
	dialectService := dialect.NewService(dialectRepo)

	translationRepo := translation.NewRepo(db)
	translationService := translation.NewService(translationRepo)

	emailService := email.NewService(db)
	emailService.Mailer(logger)

	app := fiber.New()

	if !*PRODUCTION_MOD {
		app.Use(log.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     ALLOW_ORIGINS,
		AllowMethods:     "GET,POST,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Cookie, Authorization, Set-Cookie",
		AllowCredentials: true,
	}))

	app.Use(middlewares.Logrus(logger))

	routes.PrivateEndpoints(app, SECRET_KEY, translatorService, dialectService, translationService, activeTranslatorsTracker)
	routes.PublicEndpoints(app, SECRET_KEY, translatorService, dialectService, translationService, emailService)

	_ = app.Listen(":9321")
}
