package usersService

import "mybot/repository/users"

type UserResponse struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Active []string `json:"active"`
	Equity float64  `json:"equity"`
}

type UserService interface {
	Create(users.User) (string, error)
	GetAll() ([]UserResponse, error)
	GetUserById(string) (*UserResponse, error)
	UpdateUserById(string, users.User) (*UserResponse, error)
	DeleteUserById(string) (int, error)
}
