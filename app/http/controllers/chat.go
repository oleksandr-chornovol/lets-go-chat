package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"

	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"github.com/oleksandr-chornovol/lets-go-chat/cache"
)

func StartEcho(response http.ResponseWriter, request *http.Request) {
	tokenId := request.URL.Query().Get("token")

	token, err := models.GetTokenById(tokenId)
	if err != nil {
		log.Println("GetTokenById failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ! token.IsEmpty() {
		if time.Now().String() < token.ExpiresAt {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(response, request, nil)
			if err != nil {
				log.Print("upgrade failed: ", err)
				return
			}
			defer conn.Close()

			user, _ := models.GetUserById(token.UserId)
			cache.AddUser(user)
			defer cache.DeleteUser(user.Id)

			for {
				mt, message, err := conn.ReadMessage()

				err = conn.WriteMessage(mt, message)
				if err != nil {
					log.Println("write failed:", err)
					break
				}
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(response, "Token is expired.")
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(response, "Token does not exist.")
	}
}

func GetActiveUsersCount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]int {
		"count_of_users": len(cache.GetAllUsers()),
	})
}
