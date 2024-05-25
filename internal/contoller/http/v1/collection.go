package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"merge-api/internal/service"
	"merge-api/internal/service/collection"
	"net/http"
	"strconv"
)

type collectionRoutes struct {
	collectionService service.Collection
}

func NewCollectionRouters(router *fiber.Router, service service.Collection) {
	r := &collectionRoutes{collectionService: service}

	(*router).Add("GET", "/collections", r.getCollections())
	(*router).Add("GET", "/collection/:id", r.getCollection())
	(*router).Add("POST", "/collection", r.createCollection())
}

func sendError(c *fiber.Ctx, errorCode int, err error) error {
	return c.Status(errorCode).JSON(map[string]string{
		"error": err.Error(),
	})
}

func (r *collectionRoutes) getCollection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collectionId := c.Params("id", "")
		if collectionId == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		params := new(collection.GetCollectionInput)
		collectionIdNumber, err := strconv.Atoi(collectionId)
		if err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}
		params.Id = uint(collectionIdNumber)

		ctx := context.Background()
		result, err := r.collectionService.GetCollection(ctx, params)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func (r *collectionRoutes) getCollections() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(collection.GetCollectionsInput)
		if err := c.QueryParser(params); err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		result, err := r.collectionService.GetCollections(ctx, params)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func (r *collectionRoutes) createCollection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(collection.CreateCollectionInput)
		if err := c.BodyParser(body); err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		result, err := r.collectionService.CreateCollection(ctx, body)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
