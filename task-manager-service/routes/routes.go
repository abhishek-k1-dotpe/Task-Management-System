package routes

import (
	"task-manger-service/handler"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	taskRoute := router.Group("/task")
	taskRoute.POST("/create", handler.CreateTask)
	taskRoute.PUT("admin/update", handler.UpdateTask)
	taskRoute.GET("/get", handler.SearchTasks)
	taskRoute.PUT("/complete", handler.MarkComplete)
}
