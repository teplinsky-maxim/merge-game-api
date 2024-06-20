package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"merge-api/api/internal/service"
	"merge-api/api/internal/service/task"
	"net/http"
)

type taskRoutes struct {
	taskService service.Task
}

func NewTasksRouter(router *fiber.Router, service service.Task) {
	r := &taskRoutes{taskService: service}

	(*router).Add("POST", "/task/board", r.handleCreateBoardTask())
	(*router).Add("POST", "/task/move", r.handleMoveItemTask())
	(*router).Add("POST", "/task/merge", r.handleMergeItemsTask())
}

func sendTaskUUID(c *fiber.Ctx, uuid uuid.UUID) error {
	return c.Status(200).JSON(map[string]string{
		"task_uuid": uuid.String(),
	})
}

func (r *taskRoutes) handleCreateBoardTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(task.CreateNewBoardTaskInput)
		err := c.BodyParser(body)
		if err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		createdTask, err := r.taskService.CreateTaskNewBoard(ctx, body.Width, body.Height)
		if err != nil {
			return err
		}

		return c.JSON(createdTask)
	}
}

func (r *taskRoutes) handleMoveItemTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(task.MoveItemTaskInput)
		err := c.BodyParser(body)
		if err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		createdTask, err := r.taskService.CreateTaskMoveItem(ctx, body.BoardID, body.W1, body.H1, body.W2, body.H2)
		if err != nil {
			return err
		}

		return c.JSON(createdTask)
	}
}

func (r *taskRoutes) handleMergeItemsTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(task.MergeItemsTaskInput)
		err := c.BodyParser(body)
		if err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		createdTask, err := r.taskService.CreateTaskMergeItems(ctx, body.BoardID, body.W1, body.H1, body.W2, body.H2)
		if err != nil {
			return err
		}

		return c.JSON(createdTask)
	}
}
