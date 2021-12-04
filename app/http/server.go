package http

import (
	"log"
	"net/http"
)

func StartServer() {
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
