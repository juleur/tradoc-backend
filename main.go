package main

import (
	"tradoc/data"
	"tradoc/db"
	"tradoc/pkg/store"
	"tradoc/pkg/token"
	"tradoc/tasks"
	"tradoc/web/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
)

const TRANSLATIONS_PATH = "/var/www/files/translations"

var (
	SECRET_KEY string
)

func init() {
	SECRET_KEY = token.HexKeyGenerator(36)
}

func main() {
	logger := &log.Logger{
		Level:  log.ParseLevel("info"),
		Caller: 1,
		Writer: &log.FileWriter{
			Filename:     "logs/main.log",
			FileMode:     0600,
			MaxSize:      100 * 1024 * 1024,
			MaxBackups:   7,
			EnsureFolder: true,
			LocalTime:    true,
		},
	}

	database := db.OpenDB()
	defer database.Close()

	dbRepo := db.NewDBPsql(database)

	dialectsTableName, err := dbRepo.FindAllDialect()
	if err != nil {
		logger.Fatal().Msg(err.Error.Error())
	}

	onGoingTranslationsStore := store.NewOnGoingTranslationsStore()
	onlineTranslatorsStore := store.NewOnlineTranslatorsStore()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logger)
		c.Locals("SECRET_KEY", []byte(SECRET_KEY))
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://trad-oc.xyz",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Accept, Accept-Encoding, Accept-Language, Authorization, Cache-Control, Content-Type, Expires, Host, Origin, Pragma, Referer, User-Agent, X-Refresh-Token, X-User-ID",
	}))

	controllers.MakeAuthControllers(app, dbRepo, onlineTranslatorsStore)
	controllers.MakeDialectControllers(app, dbRepo, onGoingTranslationsStore, onlineTranslatorsStore, data.DIALECTS, dialectsTableName)

	c := cron.New()

	tasks.RunAllTasks(c, logger, database, onGoingTranslationsStore, onlineTranslatorsStore, dialectsTableName, TRANSLATIONS_PATH)

	c.Start()
	defer c.Stop()

	logger.Info().Msg("Server started")
	if err := app.Listen(":2345"); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}
