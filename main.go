package main

import (
	"github.com/oleksandr-chornovol/lets-go-chat/app/http"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

func main() {
	database.Init()
	database.Migrate()
	http.InitRoutes()
	http.StartServer()
}
