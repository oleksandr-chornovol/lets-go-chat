package models

import (
	"github.com/google/uuid"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"log"
	"time"
)

type Token struct {
	Id string
	UserId string
	ExpiresAt string
}

func CreateToken(userId string) string {
	token := map[string]string{
		"id": uuid.New().String(),
		"user_id": userId,
		"expires_at": time.Now().Add(time.Minute).String(),
	}
	database.Driver.Insert("tokens", token)

	return token["id"]
}

func GetTokenById(id string) (Token, bool) {
	whereAttributes := map[string]string{"id": id}
	tokens := database.Driver.Select("tokens", whereAttributes)

	var token Token
	for tokens.Next() {
		err := tokens.Scan(&token.Id, &token.UserId, &token.ExpiresAt)
		if err != nil {
			log.Println(err)
		}
	}

	return token, token != Token{}
}
