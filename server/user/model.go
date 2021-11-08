package user

import "pkg/hasher"

type User struct {
	Id int
	Name string
	Password string
}

var users []User

func CreateUser(user User) User {
	user.Id = generateIdForNewUser()
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

func generateIdForNewUser() int {
	if len(users) == 0 {
		return 1
	} else {
		lastUserId := users[len(users)-1].Id
		return lastUserId + 1
	}
}
