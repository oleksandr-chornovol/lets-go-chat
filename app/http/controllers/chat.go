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
	"github.com/oleksandr-chornovol/lets-go-chat/config"
)

type ChatController struct {
	ActiveUsersCache cache.ActiveUsersCacheInterface
	MessageModel models.MessageInterface
	TokenModel models.TokenInterface
	UserModel models.UserInterface
}

func (c *ChatController) StartChat(response http.ResponseWriter, request *http.Request) {
	tokenId := request.URL.Query().Get("token")

	token, err := c.TokenModel.GetTokenById(tokenId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(response, "Token does not exist.")
			return
		} else {
			log.Println("GetTokenById failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if time.Now().String() < token.ExpiresAt {
		wsUpgrader := websocket.Upgrader{}
		conn, err := wsUpgrader.Upgrade(response, request, nil)
		if err != nil {
			log.Println("Websocket Upgrade failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := c.UserModel.GetUserByField("id", token.UserId)
		if err != nil {
			log.Println("GetUserByField failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Connection = conn
		c.ActiveUsersCache.AddUser(user)

		c.WriteMissedMessages(user)

		c.ReadMessages(user)

		c.EndUserSession(user)
	} else {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Token is expired.")
	}
}

func (c *ChatController) ReadMessages(user models.User) {
	for {
		_, message, err := user.Connection.ReadMessage()
		if err != nil {
			break
		}

		_, err = c.MessageModel.CreateMessage(models.Message{UserId: user.Id, Text: string(message)})
		if err != nil {
			log.Println("CreateMessage failed, err:", err)
		}

		go c.BroadcastMessage(message)
	}
}

func (c *ChatController) BroadcastMessage(message []byte) {
	for _, user := range c.ActiveUsersCache.GetAllUsers() {
		err := user.Connection.WriteMessage(1, message)
		if err != nil {
			log.Println("WriteMessage failed:", err)
		}
	}
}

func (c *ChatController) WriteMissedMessages(user models.User) {
	for _, message := range c.MessageModel.GetMessagesFromTime(user.LastSessionEnd) {
		err := user.Connection.WriteMessage(1, []byte(message.Text))
		if err != nil {
			log.Println("WriteMessage failed:", err)
		}
	}
}

func (c *ChatController) EndUserSession(user models.User) {
	c.ActiveUsersCache.DeleteUser(user.Id)

	err := user.Connection.Close()
	if err != nil {
		log.Println("Connection close failed:", err)
	}

	user.LastSessionEnd = time.Now().String()
	_, err = c.UserModel.UpdateUser(user)
	if err != nil {
		log.Println("UpdateUser failed:", err)
	}
}

func (c *ChatController) GetActiveUsersCount(response http.ResponseWriter, request *http.Request) {
	log.Println("!!!DB_URL!!!", config.Get("db_url"))
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]int {
		"count_of_users": len(c.ActiveUsersCache.GetAllUsers()),
	})
}
