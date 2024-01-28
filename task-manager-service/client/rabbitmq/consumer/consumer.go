package consumer

import (
	"fmt"
	"log"
	"os"
	rabbitmq "task-manger-service/client/rabbitmq/configuration"
)

func startConsumer(consumer consumerObj) {
	conn, err := rabbitmq.GetConnection()
	if err != nil {
		fmt.Println("Failed to open a rabbitmq connection")
		fmt.Println(err)
		os.Exit(0)
	}
	defer conn.Close()
	exchangeName := consumer.ExchangeName
	exchangeType := consumer.ExchangeType
	amqpChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer amqpChannel.Close()

	err = amqpChannel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		fmt.Println("Failed to declare a rabbitmq exchange")
		fmt.Println(err)
		os.Exit(0)
	}

	queue, err := amqpChannel.QueueDeclare(
		consumer.QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = amqpChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = amqpChannel.QueueBind(
		consumer.QueueName,    // queue name
		consumer.RoutingKey,   // routing key
		consumer.ExchangeName, // exchange
		false,
		nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	msgs, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message for %s ) : %s ExchangeName: %s", consumer.QueueName, d.Body, d.Exchange)
			consumeMessage(queue.Name, d, consumer)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C : %s", consumer.QueueName)
	<-forever
}
