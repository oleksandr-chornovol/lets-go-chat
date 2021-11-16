package user

import (
	"github.com/google/uuid"
	"github.com/oleksandr-chornovol/lets-go-chat/database/drivers"
	"log"
	"pkg/hasher"
)

type User struct {
	Id string
	Name string
	Password string
}

var DBDriver drivers.DBDriverInterface

func CreateUser(user User) User {
	user.Id = uuid.New().String()
	user.Password, _ = hasher.HashPassword(user.Password)

	attributes := map[string]string{
		"id": user.Id,
		"name": user.Name,
		"password": user.Password,
	}
	DBDriver.Insert("users", attributes)

	return user
}

func GetUserByName(name string) (User, bool) {
	whereAttributes := map[string]string{"name": name}
	users := DBDriver.Select("users", whereAttributes)

	var user User
	for users.Next() {
		err := users.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err)
		}
	}

	return user, user != User{}
}
