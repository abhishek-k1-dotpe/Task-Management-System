package publisher

import (
	rabbitmq "task-manger-service/client/rabbitmq/configuration"

	"errors"

	"github.com/streadway/amqp"
)

type PublishTaskRequest struct {
	ReqBytes     []byte `json:"reqBytes"`
	ExchangeName string `json:"exchangeName"`
	RoutingKey   string `json:"routingKey"`
	Headers      amqp.Table
}

func (publishTaskRequest *PublishTaskRequest) PublishTask() error {

	reqBytes := publishTaskRequest.ReqBytes
	exchangeName := publishTaskRequest.ExchangeName
	routingKey := publishTaskRequest.RoutingKey
	headers := publishTaskRequest.Headers

	if exchangeName == "" {
		return errors.New("exchange name missing")
	}

	if routingKey == "" {
		return errors.New("routing key missing")

	}

	if reqBytes == nil {
		return errors.New("no data to publish")
	}

	conn, err := rabbitmq.GetConnection()
	if err != nil {
		return err
	}

	amqpChannel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer amqpChannel.Close()

	err = amqpChannel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = amqpChannel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         reqBytes,
			Headers:      headers,
		})
	if err != nil {
		return err
	}
	return nil
}
