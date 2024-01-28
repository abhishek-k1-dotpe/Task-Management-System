package handler

import (
	"net/http"
	"strconv"
	"user-service/domain"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	request := domain.User{}

	//json binding to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}
	// jsno validation
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}
	// checking if user type is default or admin for creation of admin validation is not required
	if request.CreatedBy != 0 {
		isAdmin, err := service.IsAdmin(c, request.CreatedBy)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.UserResonse{
				Status:    domain.SomethingWentWrong,
				RespError: err.Error(),
			})
			return
		}
		if !isAdmin {
			c.JSON(http.StatusBadRequest, domain.UserResonse{
				Status:    domain.SomethingWentWrong,
				RespError: "user is not admin",
			})
			return
		}
	}
	// calling service to create user
	userId, err := service.CreateUser(c, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.UserResonse{
		UserId: userId,
		Status: "user created succesfully",
	})
}
func DeleteUser(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}
	// validation for admin can be done by requesting createdBy in delete payload but since we can easily do by authentication system
	err = service.DeleteUser(c, int(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.UserResonse{
		Status: "user deleted succesfully",
	})

}

func GetUserInfo(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}
	user, err := service.GetUserInfo(c, int(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.UserResonse{
			Status:    domain.SomethingWentWrong,
			RespError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)

}
