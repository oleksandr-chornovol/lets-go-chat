package http

import (
	"log"
	"net/http"
	"os"
)

func StartServer() {
	//err := http.ListenAndServe("localhost:8080", nil) // local
	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil) // heroku
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
