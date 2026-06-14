package http

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, compressHandler *CompressHandler) {
	api := app.Group("/api/v1")

	api.Post("/compress", compressHandler.Compress)
}