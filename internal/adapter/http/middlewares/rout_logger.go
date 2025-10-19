package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	DateTimeLayout = "15:04:05 02-01-2006"
)

func RoutLoggerMiddlewareInitializer(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${method} ${path} ${status} pid=${pid}\n",
		TimeFormat: DateTimeLayout,
	}))
}
