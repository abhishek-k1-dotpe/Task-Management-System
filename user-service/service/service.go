package service

import (
	"context"
	"user-service/domain"
	"user-service/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c context.Context, user domain.User) (int, error) {
	return models.Create(c, user)
}
func DeleteUser(c context.Context, userId int) error {
	return models.DeleteUser(c, userId)

}
func GetUserInfo(c *gin.Context, userId int) (*domain.User, error) {
	return models.GetUserInfo(c, userId)
}

func IsAdmin(c context.Context, userId int) (bool, error) {
	user, err := models.GetUserInfo(c, userId)
	if err != nil {
		return false, err
	}
	return user.Type == domain.Admin, nil
}
