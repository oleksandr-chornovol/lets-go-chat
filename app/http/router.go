package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oleksandr-chornovol/lets-go-chat/app/http/controllers"
	"net/http"
)

var router = chi.NewRouter()

type Route struct {
	Path string
	Handler func(http.ResponseWriter, *http.Request)
}

func InitRoutes() {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/v1/user", controllers.Register)
	router.Post("/v1/user/login", controllers.Login)
	router.Get("/v1/user/active", controllers.GetActiveUsersCount)
	router.Get("/v1/chat", controllers.StartEcho)
}
