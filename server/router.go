package server

import (
	"github.com/oleksandr-chornovol/lets-go-chat/server/user"
	"net/http"
)

type Route struct {
	Path string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route {
	{"/v1/user", user.Register},
	{"/v1/user/login", user.Login},
}

func RouterInit() {
	for _, route := range routes {
		http.HandleFunc(route.Path, route.Handler)
	}
}
