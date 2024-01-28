package domain

import "errors"

type User struct {
	UserId    int    `json:"userId,omitempty"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Type      string `json:"type"` // "default" or "admin"
	CreatedBy int    `json:"createdBy"`
}

func (u *User) Validate() error {
	if len(u.Username) <= 0 {
		return errors.New("invalid username")
	}

	if len(u.Email) <= 0 {
		return errors.New("invalid email")
	}
	if len(u.Type) <= 0 {
		return errors.New("invalid type")
	}
	if u.CreatedBy < 0 {
		return errors.New("invalid createdBy")
	}

	if u.CreatedBy == 0 && u.Type != Admin {
		return errors.New("invalid createdBy")
	}

	return nil

}

type UserResonse struct {
	UserId    int    `json:"userId,omitempty"`
	Status    string `json:"status"`
	RespError string `json:"error,omitempty"`
}
