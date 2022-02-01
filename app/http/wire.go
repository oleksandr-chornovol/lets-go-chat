//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/oleksandr-chornovol/lets-go-chat/app/http/controllers"
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"github.com/oleksandr-chornovol/lets-go-chat/cache"
)

func NewUserController() *controllers.UserController {
	wire.Build(
		models.NewTokenModel,
		models.NewUserModel,

		wire.Bind(new(models.TokenInterface), new(*models.Token)),
		wire.Bind(new(models.UserInterface), new(*models.User)),

		controllers.NewUserController,
	)

	return &controllers.UserController{}
}

func NewChatController() *controllers.ChatController {
	wire.Build(
		cache.NewActiveUsersCache,
		models.NewMessageModel,
		models.NewTokenModel,
		models.NewUserModel,

		wire.Bind(new(cache.ActiveUsersCacheInterface), new(*cache.ActiveUsersCache)),
		wire.Bind(new(models.MessageInterface), new(*models.Message)),
		wire.Bind(new(models.TokenInterface), new(*models.Token)),
		wire.Bind(new(models.UserInterface), new(*models.User)),

		controllers.NewChatController,
	)

	return &controllers.ChatController{}
}