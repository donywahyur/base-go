package config

import (
	"go_base/internal/handlers"
	"go_base/internal/repositories"
	"go_base/internal/services"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	FiberApp    *fiber.App
	UserHandler *handlers.UserHandler
	Middleware  *Middleware
}

func InitializeApp() *App {
	fiberApp := fiber.New()
	postgreDb := GetDB()

	userRepo := repositories.NewUserRepository(postgreDb)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	middleware := NewMiddleware(userRepo)

	return &App{
		FiberApp:    fiberApp,
		UserHandler: userHandler,
		Middleware:  middleware,
	}
}
