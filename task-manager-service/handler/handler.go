package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"task-manger-service/domain"
	"task-manger-service/service"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	request := domain.Task{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	fmt.Println(request)
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}

	taskId, err := service.CreateTask(c, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &domain.TaskCreateResponse{
		Status: "task created",
		TaskId: taskId,
	})
}

func UpdateTask(c *gin.Context) {
	request := domain.Task{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	if request.TaskId == 0 {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: "empty taskId",
		})
		return
	}
	err := service.UpdateTask(c, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &domain.TaskResponse{
		Status: "task successfully updated",
	})
}

func SearchTasks(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	sortBy := c.Query("sortBy")
	isAdminStr := c.Query("isAdmin")
	isAdmin := false
	if isAdminStr == "true" {
		isAdmin = true
	}
	tasks, err := service.GetTasks(c, int(userId), sortBy, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &domain.TaskSearchResponse{
		Status: "successfully fetched",
		Tasks:  tasks,
	})
}

func MarkComplete(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	err = service.MarkComplete(c, int(taskId))
	if err != nil {
		c.JSON(http.StatusBadRequest, &domain.TaskResponse{
			Status:    "something went wrong",
			RespError: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &domain.TaskResponse{
		Status: "task successfully updated",
	})
}
