package server

import (
	"log"
	"net/http"
)

func Start() {
	RouterInit()
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
