package router

import (
	"github.com/gofiber/fiber/v2"
	"go-shortfile/controllers"
)

func AddRouters(app *fiber.App) {
	app.Get("/api/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to shortFile")
	})

	app.Post("/api/u", controllers.UploadFile)
	app.Get("/api/i/:id", controllers.GetFileInfo)
	app.Get("/api/d/:id", controllers.DownloadFile)
}
