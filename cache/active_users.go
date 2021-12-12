package cache

import "github.com/oleksandr-chornovol/lets-go-chat/app/models"

var users = make(map[string]models.User)

func AddUser(user models.User) {
	users[user.Id] = user
}

func DeleteUser(userId string) {
	delete(users, userId)
}

func GetAllUsers() map[string]models.User {
	return users
}