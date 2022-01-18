package http

import (
	"log"
	"net/http"

	"github.com/oleksandr-chornovol/lets-go-chat/config"
)

func StartServer() {
	err := http.ListenAndServe(":" + config.Get("port"), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
