package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"

	"github.com/oleksandr-chornovol/lets-go-chat/app"
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"github.com/oleksandr-chornovol/lets-go-chat/cache"
)

type ChatController struct {
	ActiveUsersCache cache.ActiveUsersCacheInterface
	MessageModel     models.MessageInterface
	TokenModel       models.TokenInterface
	UserModel        models.UserInterface
	chMessage        chan []byte
}

func NewChatController(
	activeUsersCache cache.ActiveUsersCacheInterface,
	messageModel models.MessageInterface,
	tokenModel models.TokenInterface,
	userModel models.UserInterface,
) *ChatController {
	return &ChatController{
		ActiveUsersCache: activeUsersCache,
		MessageModel:     messageModel,
		TokenModel:       tokenModel,
		UserModel:        userModel,
	}
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

		c.writeMissedMessages(user)

		c.chMessage = make(chan []byte, 1000)

		go c.readMessages(user)
		go c.processMessages(user)
	} else {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Token is expired.")
	}
}

func (c *ChatController) readMessages(user models.User) {
	for {
		_, message, err := user.Connection.ReadMessage()
		if err != nil {
			c.endUserSession(user)
			break
		}
		c.chMessage <- message
	}
}

func (c *ChatController) processMessages(user models.User) {
	for {
		msg, ok := <-c.chMessage
		if !ok {
			break
		}
		_, err := c.MessageModel.CreateMessage(models.Message{UserId: user.Id, Text: string(msg)})
		if err != nil {
			log.Println("CreateMessage failed, err:", err)
		}

		c.broadcastMessage(msg)
	}
}

func (c *ChatController) broadcastMessage(message []byte) {
	for _, user := range c.ActiveUsersCache.GetAllUsers() {
		err := user.Connection.WriteMessage(1, message)
		if err != nil {
			log.Println("WriteMessage failed:", err)
		}
	}
}

func (c *ChatController) writeMissedMessages(user models.User) {
	for _, message := range c.MessageModel.GetMessagesFromTime(user.LastSessionEnd) {
		err := user.Connection.WriteMessage(1, []byte(message.Text))
		if err != nil {
			log.Println("WriteMessage failed:", err)
		}
	}
}

func (c *ChatController) endUserSession(user models.User) {
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
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(app.ActiveUsersResponse{
		Count: len(c.ActiveUsersCache.GetAllUsers()),
	})
}
