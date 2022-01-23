package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

type UserInterface interface {
	CreateUser(user User) (User, error)
	GetUserByField(field string, value string) (User, error)
	UpdateUser(user User) (User, error)
}

type User struct {
	Id             string
	Name           string
	Password       string
	LastSessionEnd string
	Connection     *websocket.Conn
}

func NewUserModel() *User {
	return &User{}
}

func (u *User) CreateUser(user User) (User, error) {
	user.Id = uuid.New().String()
	user.Password, _ = hasher.HashPassword(user.Password)

	attributes := map[string]string{
		"id":       user.Id,
		"name":     user.Name,
		"password": user.Password,
	}
	err := database.Driver.Insert("users", attributes)

	return user, err
}

func (u *User) GetUserByField(field string, value string) (User, error) {
	var whereAttributes = [][3]string{
		{field, "=", value},
	}
	result := database.Driver.SelectRow("users", whereAttributes)

	var user User
	err := result.Scan(&user.Id, &user.Name, &user.Password, &user.LastSessionEnd)

	return user, err
}

func (u *User) UpdateUser(user User) (User, error) {
	whereAttributes := map[string]string{
		"id": user.Id,
	}

	updateAttributes := map[string]string{
		"name":             user.Name,
		"password":         user.Password,
		"last_session_end": user.LastSessionEnd,
	}

	err := database.Driver.Update("users", whereAttributes, updateAttributes)

	return user, err
}
