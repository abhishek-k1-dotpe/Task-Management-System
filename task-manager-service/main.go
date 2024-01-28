package main

import (
	"fmt"
	"os"
	rabbitmq "task-manger-service/client/rabbitmq/configuration"
	"task-manger-service/client/rabbitmq/consumer"
	"task-manger-service/db"
	"task-manger-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.InitRoutes(router)
	err := db.InitMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	// initializing rabbitMq
	err = rabbitmq.Init()
	if err != nil {
		fmt.Println("Hii")
		fmt.Println(err)
		os.Exit(0)
	}
	consumer.StartConsumers()

	err = router.Run(":8081")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
