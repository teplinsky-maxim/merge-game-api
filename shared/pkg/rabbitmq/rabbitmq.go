package rabbitmq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"merge-api/shared/config"
)

const (
	ExchangeName = "TaskExchange"
	QueueName    = "TaskQueue"
	QueueKey     = "TaskKey"
)

type RabbitMQ struct {
	Connection *amqp091.Connection
}

func NewRabbitMQ(config *config.Config) (RabbitMQ, error) {
	dsn := makeDsn(config)
	conn, err := amqp091.Dial(dsn)
	return RabbitMQ{
		Connection: conn,
	}, err
}

func Initialize(mq *RabbitMQ) error {
	ch, err := mq.Connection.Channel()
	if err != nil {
		return err
	}
	err = makeMainExchange(ch)
	if err != nil {
		return err
	}
	q, err := makeMainQueue(ch)
	err = ch.QueueBind(
		q.Name,       // queue name
		QueueKey,     // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = ch.Close()
	return err
}

func makeMainExchange(channel *amqp091.Channel) error {
	err := channel.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func makeMainQueue(channel *amqp091.Channel) (amqp091.Queue, error) {
	q, err := channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return q, err
}

func makeDsn(config *config.Config) string {
	template := "amqp://%v:%v@%v:%v/"
	dsn := fmt.Sprintf(template, config.RabbitMQ.User, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port)
	return dsn
}
