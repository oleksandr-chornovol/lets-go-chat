package models

import (
	"github.com/google/uuid"
	"log"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

type User struct {
	Id string
	Name string
	Password string
}

func (u User) IsEmpty() bool {
	return u == User{}
}

func CreateUser(user User) (User, error) {
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

func GetUserById(id string) (User, error) {
	var user User

	whereAttributes := map[string]string{"id": id}
	users, err := database.Driver.Select("users", whereAttributes)
	if err != nil {
		return user, err
	}

	for users.Next() {
		err := users.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}

	return user, nil
}

func GetUserByName(name string) (User, error) {
	var user User

	whereAttributes := map[string]string{"name": name}
	users, err := database.Driver.Select("users", whereAttributes)
	if err != nil {
		return user, err
	}

	for users.Next() {
		err := users.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}

	return user, nil
}
