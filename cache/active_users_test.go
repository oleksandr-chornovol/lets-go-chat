package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
)

func TestAddUser(t *testing.T) {
	activeUsersCache := NewActiveUsersCache()
	user := models.User{Id: "user_id", Name: "name", Password: "password"}

	activeUsersCache.AddUser(user)

	assert.Equal(t, user, activeUsersCache.users[user.Id])
}

func TestDeleteUser(t *testing.T) {
	activeUsersCache := NewActiveUsersCache()
	user := models.User{Id: "user_id", Name: "name", Password: "password"}
	activeUsersCache.users[user.Id] = user

	activeUsersCache.DeleteUser(user.Id)

	_, ok := activeUsersCache.users[user.Id]
	assert.Equal(t, false, ok)
}

func TestGetAllUsers(t *testing.T) {
	activeUsersCache := NewActiveUsersCache()
	user1 := models.User{Id: "user_id1", Name: "name", Password: "password"}
	user2 := models.User{Id: "user_id2", Name: "name", Password: "password"}
	activeUsersCache.users[user1.Id] = user1
	activeUsersCache.users[user2.Id] = user2

	allUsers := activeUsersCache.GetAllUsers()

	assert.Equal(t, 2, len(allUsers))
	assert.Equal(t, user1, allUsers[user1.Id])
	assert.Equal(t, user2, allUsers[user2.Id])
}
