package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"merge-api/internal/service"
	"merge-api/internal/service/collection"
	"net/http"
)

type collectionRoutes struct {
	collectionService service.Collection
}

func NewCollectionRouters(router *fiber.Router, service service.Collection) {
	r := &collectionRoutes{collectionService: service}

	(*router).Add("GET", "/collection", r.getCollection())
}

func sendError(c *fiber.Ctx, errorCode int, err error) error {
	return c.Status(errorCode).JSON(map[string]string{
		"error": err.Error(),
	})
}

func (r *collectionRoutes) getCollection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(collection.GetCollectionInput)
		if err := c.QueryParser(params); err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		result, err := r.collectionService.GetCollection(ctx, params)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
