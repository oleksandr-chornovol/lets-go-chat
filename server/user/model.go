package user

import (
	"github.com/google/uuid"
	"pkg/hasher"
)

type User struct {
	Id string
	Name string
	Password string
}

var users []User

func CreateUser(user User) User {
	user.Id = uuid.New().String()
	user.Password, _ = hasher.HashPassword(user.Password)
	users = append(users, user)
	return user
}

func GetUserByName (name string) (User, bool) {
	for _, user := range users {
		if user.Name == name {
			return user, true
		}
	}
	return User{}, false
}
