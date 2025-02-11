package config

import (
	"go_base/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func LoadRoute(app *App) {
	app.FiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Content-Length, Accept, X-Requested-With, Authorization, X-Forwarded-For",
	}))
	app.FiberApp.Use(recover.New())

	app.FiberApp.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        10,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(utils.ResponseApi(fiber.StatusTooManyRequests, "Too many requests", nil))
		},
	}))

	app.FiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(utils.ResponseApi(fiber.StatusOK, "Welcome to go-base", nil))
	})

	api := app.FiberApp.Group("/api")
	v1 := api.Group("/v1")
	authentication := v1.Group("/auth")
	authentication.Post("/login", app.UserHandler.Login)

	user := v1.Group("/user", app.Middleware.CheckAuth)
	user.Get("/:id", app.UserHandler.GetUser)

}
