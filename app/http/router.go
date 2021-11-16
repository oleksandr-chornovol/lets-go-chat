package http

import (
	userController "github.com/oleksandr-chornovol/lets-go-chat/app/http/controllers"
	"net/http"
)

type Route struct {
	Path string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route {
	{"/v1/user", userController.Register},
	{"/v1/user/login", userController.Login},
}

func InitRoutes() {
	for _, route := range routes {
		http.HandleFunc(route.Path, route.Handler)
	}
}
