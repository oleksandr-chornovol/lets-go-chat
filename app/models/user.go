package models

import (
	"github.com/google/uuid"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"log"
	"pkg/hasher"
)

type User struct {
	Id string
	Name string
	Password string
}

func CreateUser(user User) User {
	user.Id = uuid.New().String()
	user.Password, _ = hasher.HashPassword(user.Password)

	attributes := map[string]string{
		"id": user.Id,
		"name": user.Name,
		"password": user.Password,
	}
	database.Driver.Insert("users", attributes)

	return user
}

func GetUserById(id string) (User, bool) {
	whereAttributes := map[string]string{"id": id}
	users := database.Driver.Select("users", whereAttributes)

	var user User
	for users.Next() {
		err := users.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}

	return user, user != User{}
}

func GetUserByName(name string) (User, bool) {
	whereAttributes := map[string]string{"name": name}
	users := database.Driver.Select("users", whereAttributes)

	var user User
	for users.Next() {
		err := users.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}

	return user, user != User{}
}
