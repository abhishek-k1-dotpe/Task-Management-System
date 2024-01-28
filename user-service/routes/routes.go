package routes

import (
	"user-service/handler"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	userRoute := router.Group("/user")
	userRoute.POST("/admin/create", handler.CreateUser)
	userRoute.DELETE("/admin/delete", handler.DeleteUser)
	userRoute.GET("/get", handler.GetUserInfo)
}
