package auth

import (
	"github.com/gofiber/fiber/v2"
)


func SetupRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	authGroup.Get("/login", loginHandler)
}

func loginHandler(c *fiber.Ctx) error {
	return c.SendString("Login page")
}
