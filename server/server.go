package server

import (
	"log"
	"net/http"
	"os"
)

func Start() {
	RouterInit()
	// err := http.ListenAndServe("localhost:8080", nil) // local
	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
