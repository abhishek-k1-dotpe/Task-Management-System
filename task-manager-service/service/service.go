package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"task-manger-service/client"
	rabbitmq "task-manger-service/client/rabbitmq/configuration"
	"task-manger-service/client/rabbitmq/consumer"
	"task-manger-service/client/rabbitmq/publisher"
	"task-manger-service/domain"
	"task-manger-service/models"

	"github.com/streadway/amqp"
)

func CreateTask(ctx context.Context, task domain.Task) (int, error) {
	// checking if provided id is admin id
	adminUser, err := client.GetUser(task.AssignedBy)
	if err != nil {
		return 0, err
	}

	//checking if provided id is user id
	defaultUser, err := client.GetUser(task.AssignedTo)

	if adminUser.Type != "admin" || defaultUser.Type != "default" {
		return 0, errors.New("userType is invalid")
	}

	if err != nil {
		return 0, err
	}
	return models.CreateTask(ctx, task)
}

func UpdateTask(ctx context.Context, task domain.Task) error {
	return models.UpdateTask(ctx, task)
}

func GetTasks(ctx context.Context, userId int, sortBy string, isAdmin bool) ([]*domain.Task, error) {
	// validation  for admin
	if isAdmin {
		user, err := client.GetUser(userId)
		if err != nil {
			return nil, err
		}
		if user.Type != "admin" {
			return nil, errors.New("user is not admin")
		}
	}
	return models.GetTasks(ctx, userId, sortBy, isAdmin)
}
func MarkComplete(ctx context.Context, taskId int) error {
	err := models.MarkComplete(ctx, taskId)
	if err != nil {
		return err
	}
	err = raiseEvent(taskId)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func raiseEvent(taskId int) error {

	publishData := consumer.TaskRequestData{
		Data: taskId,
	}

	reqBytes, err := json.Marshal(publishData)
	if err != nil {
		return err
	}

	headers := make(amqp.Table)
	headers["x-attempt-count"] = 0
	publishRequest := publisher.PublishTaskRequest{}

	publishRequest.ExchangeName = rabbitmq.TaskEventExchange
	publishRequest.RoutingKey = rabbitmq.TaskRoutingKey
	publishRequest.ReqBytes = reqBytes
	publishRequest.Headers = headers

	err = publishRequest.PublishTask()
	if err != nil {
		return err
	}

	return nil
}
