package models

import (
	"github.com/google/uuid"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

type UserInterface interface {
	CreateUser(user User) (User, error)
	GetUserByField(field string, value string) (User, error)
	IsEmpty() bool
}

type User struct {
	Id string
	Name string
	Password string
}

func (u User) IsEmpty() bool {
	return u == User{}
}

func (u User) CreateUser(user User) (User, error) {
	user.Id = uuid.New().String()
	user.Password, _ = hasher.HashPassword(user.Password)

	attributes := map[string]string{
		"id": user.Id,
		"name": user.Name,
		"password": user.Password,
	}
	err := database.Driver.Insert("users", attributes)

	return user, err
}

func (u User) GetUserByField(field string, value string) (User, error) {
	whereAttributes := map[string]string{field: value}
	result := database.Driver.SelectRow("users", whereAttributes)

	var user User
	err := result.Scan(&user.Id, &user.Name, &user.Password)

	return user, err
}
