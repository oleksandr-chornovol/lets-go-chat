package cache

import (
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
)

type ActiveUsersCacheInterface interface {
	AddUser(user models.User)
	DeleteUser(userId string)
	GetAllUsers() map[string]models.User
}

type ActiveUsersCache struct {
	users map[string]models.User
}

func NewActiveUsersCache() *ActiveUsersCache {
	return &ActiveUsersCache{
		users: make(map[string]models.User),
	}
}

func (c *ActiveUsersCache) AddUser(user models.User) {
	c.users[user.Id] = user
}

func (c *ActiveUsersCache) DeleteUser(userId string) {
	delete(c.users, userId)
}

func (c *ActiveUsersCache) GetAllUsers() map[string]models.User {
	return c.users
}
