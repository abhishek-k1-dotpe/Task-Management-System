package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task-manger-service/config"
	"task-manger-service/domain"
)

func GetUser(userId int) (*domain.User, error) {

	url := fmt.Sprintf("%s?userId=%d", config.UserServiceUrl, userId)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}
	var user domain.User
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
