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

func (c *ChatController) StartEcho(response http.ResponseWriter, request *http.Request) {
	tokenId := request.URL.Query().Get("token")

	token, err := c.TokenModel.GetTokenById(tokenId)
	if err != nil {
		log.Println("GetTokenById failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ! token.IsEmpty() {
		if time.Now().String() < token.ExpiresAt {
			wsUpgrader := websocket.Upgrader{}
			conn, err := wsUpgrader.Upgrade(response, request, nil)
			if err != nil {
				log.Println("Websocket Upgrade failed, err:", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			user, err := c.UserModel.GetUserByField("id", token.UserId)
			if err != nil {
				log.Println("GetUserByField failed, err:", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			c.ActiveUsersCache.AddUser(user)
			defer c.ActiveUsersCache.DeleteUser(user.Id)

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

func (c *ChatController) GetActiveUsersCount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]int {
		"count_of_users": len(c.ActiveUsersCache.GetAllUsers()),
	})
}
