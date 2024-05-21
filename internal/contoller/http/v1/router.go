package v1

import (
	"github.com/gofiber/fiber/v2"
	"merge-api/internal/service"
)

func NewRouter(app *fiber.App, services *service.Services) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	NewCollectionRouters(&v1, services.Collection)
}
