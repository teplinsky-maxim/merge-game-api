package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"merge-api/internal/service"
	"merge-api/internal/service/task"
	"net/http"
)

type taskRoutes struct {
	taskService service.Task
}

func NewTasksRouter(router *fiber.Router, service service.Task) {
	r := &taskRoutes{taskService: service}

	(*router).Add("POST", "/task", r.createTask())
}

func sendTaskUUID(c *fiber.Ctx, uuid uuid.UUID) error {
	return c.Status(200).JSON(map[string]string{
		"task_uuid": uuid.String(),
	})
}

func (r *taskRoutes) createTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(task.CreateNewBoardTaskInput)
		err := c.BodyParser(body)
		if err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		ctx := context.Background()
		taskId, err := r.taskService.CreateTaskNewBoard(ctx, body.Width, body.Height)
		if err != nil {
			return err
		}

		return sendTaskUUID(c, uuid.UUID(taskId))
	}
}
