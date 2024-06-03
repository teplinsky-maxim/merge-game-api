package task

import (
	"merge-api/shared/pkg/rabbitmq"
)

type TaskService struct {
	//repo repo.Task
	rmq *rabbitmq.RabbitMQ
}
