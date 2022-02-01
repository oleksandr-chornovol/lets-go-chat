package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

var router = chi.NewRouter()

type Route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

func InitRoutes() {
	userController := NewUserController()
	chatController := NewChatController()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/v1/user", userController.Register)
	router.Post("/v1/user/login", userController.Login)
	router.Get("/v1/user/active", chatController.GetActiveUsersCount)
	router.Get("/v1/chat", chatController.StartChat)
}
