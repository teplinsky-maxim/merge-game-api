package tasks

import (
	"encoding/json"
	"errors"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
)

var ErrorMessageNotConfirmed = errors.New("message was not ack'ed")

func SendTask(mq *rabbitmq.RabbitMQ, task task.Task) error {
	channel, err := mq.Connection.Channel()
	if err != nil {
		return err
	}
	err = sendTo(channel, rabbitmq.ExchangeName, rabbitmq.QueueKey, task)
	if err != nil {
		return err
	}
	return nil
}

func sendTo(ch *amqp091.Channel, exchange, key string, message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))

	err = ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	confirmed := <-confirms
	if confirmed.Ack {
		return nil
	} else {
		return ErrorMessageNotConfirmed
	}
}
