package http

import (
	"log"
	"net/http"
	"os"

	"github.com/oleksandr-chornovol/lets-go-chat/config"
)

func StartServer() {
	err := http.ListenAndServe(":" + getPort(), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = config.LocalPort
	}

	return port
}
