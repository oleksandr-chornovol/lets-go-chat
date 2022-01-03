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

type ChatController struct {
	ActiveUsersCache cache.ActiveUsersCacheInterface
	TokenModel models.TokenInterface
	UserModel models.UserInterface
}

func (cc ChatController) StartEcho(response http.ResponseWriter, request *http.Request) {
	tokenId := request.URL.Query().Get("token")

	token, err := cc.TokenModel.GetTokenById(tokenId)
	if err != nil {
		log.Println("GetTokenById failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ! token.IsEmpty() {
		if time.Now().String() < token.ExpiresAt {
			wsUpgrader := websocket.Upgrader{}
			conn, _ := wsUpgrader.Upgrade(response, request, nil)
			defer conn.Close()

			user, _ := cc.UserModel.GetUserByField("id", token.UserId)
			cc.ActiveUsersCache.AddUser(user)
			defer cc.ActiveUsersCache.DeleteUser(user.Id)

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

func (cc ChatController) GetActiveUsersCount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]int {
		"count_of_users": len(cc.ActiveUsersCache.GetAllUsers()),
	})
}
