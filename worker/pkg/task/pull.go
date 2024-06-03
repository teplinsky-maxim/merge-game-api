package task

import (
	"encoding/json"
	"log"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/rabbitmq"
)

func StartPullTasks(mq *rabbitmq.RabbitMQ, manager *TaskExecutorsManager) error {
	channel, err := mq.Connection.Channel()
	if err != nil {
		return err
	}

	consumeChannel, err := channel.Consume(
		rabbitmq.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for d := range consumeChannel {
		log.Printf("Received a message: %s", d.Body)
		_, unmarshalledTask, err := unmarshalTask(d.Body)
		if err != nil {
			log.Printf("Error unmarshalling task: %s", err)
			continue
		}

		_, err = manager.ExecuteTask(&unmarshalledTask)
		if err != nil {
			log.Printf("Error executing task: %s", err)
		}

		err = d.Ack(false)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalTask(message []byte) (task.Type, task.Task, error) {
	t := new(task.Task)
	err := json.Unmarshal(message, t)
	if err != nil {
		return "", task.Task{}, UnknownTask
	}
	return t.Type, *t, nil
}
