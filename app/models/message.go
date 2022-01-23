package models

import (
	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"log"
	"time"
)

type MessageInterface interface {
	CreateMessage(message Message) (Message, error)
	GetMessagesFromTime(time string) []Message
}

type Message struct {
	Id        int
	UserId    string
	Text      string
	CreatedAt string
}

func NewMessageModel() *Message {
	return &Message{}
}

func (m *Message) CreateMessage(message Message) (Message, error) {
	message.CreatedAt = time.Now().String()

	attributes := map[string]string{
		"user_id":    message.UserId,
		"text":       message.Text,
		"created_at": message.CreatedAt,
	}
	err := database.Driver.Insert("messages", attributes)

	return message, err
}

func (m *Message) GetMessagesFromTime(time string) []Message {
	var message Message
	var messages []Message

	var whereAttributes = [][3]string{
		{"created_at", ">", time},
	}
	result, err := database.Driver.Select("messages", whereAttributes)
	if err != nil {
		log.Println(err)
		return messages
	}

	for result.Next() {
		err := result.Scan(&message.Id, &message.UserId, &message.Text, &message.CreatedAt)
		if err != nil {
			log.Println(err)
		} else {
			messages = append(messages, message)
		}
	}

	return messages
}
