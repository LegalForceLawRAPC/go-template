package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"

	"github.com/LegalForceLawRAPC/go-template/api/cache"
	"github.com/LegalForceLawRAPC/go-template/api/db"
	"github.com/LegalForceLawRAPC/go-template/api/migrations"
	"github.com/LegalForceLawRAPC/go-template/api/router"
	gofibersentry "github.com/LegalForceLawRAPC/go-template/api/sentry"
	"github.com/LegalForceLawRAPC/go-template/api/utils"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func main() {
	// Set global configuration
	utils.ImportEnv()

	// Init redis
	cache.GetRedis()

	// Init Validators
	utils.InitValidators()

	// Create Fiber
	app := fiber.New(fiber.Config{})

	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)

	// initialize sentry
	gofibersentry.SentryInit()
	sentryHandler := gofibersentry.New(gofibersentry.Options{})
	app.Use(sentryHandler.Handle)

	app.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "api")
	}}))

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "",
		AllowHeaders: "*",
	}))

	//Connect and migrate the db
	if viper.GetBool("MIGRATE") {
		migrations.Migrate()
	}

	// Initialize DB
	db.InitServices()

	// Mount Routes
	router.MountRoutes(app)

	// Get Port
	port := utils.GetPort()

	// Start Fiber
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
