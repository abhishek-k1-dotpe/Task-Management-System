package consumer

import (
	"encoding/json"
	"fmt"
	rabbitmq "task-manger-service/client/rabbitmq/configuration"

	"github.com/streadway/amqp"
)

type TaskRequestData struct {
	Data interface{} `json:"data"`
}

type consumerObj struct {
	QueueName    string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

func consumeMessage(queueName string, d amqp.Delivery, consumer consumerObj) {

	reqBytes := d.Body
	consumedData := TaskRequestData{}
	err := json.Unmarshal(reqBytes, &consumedData)
	if err != nil {
		d.Ack(false)
		return
	}
	switch queueName {
	case rabbitmq.TaskEventQueue:
		go consumeData(consumedData.Data, d)
	default:
		return
	}
}

// StartConsumers ..
func StartConsumers() {
	go startConsumer(consumerObj{
		QueueName:    rabbitmq.TaskEventQueue,
		ExchangeName: rabbitmq.TaskEventExchange,
		ExchangeType: "direct",
		RoutingKey:   "task.event.key",
	})
}

func consumeData(request interface{}, del amqp.Delivery) {
	del.Ack(false)

	data := request.(float64)

	msg := fmt.Sprintf("task %d is completed", int(data))
	fmt.Println(msg)
}
