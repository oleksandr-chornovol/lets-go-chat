package models

import (
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/oleksandr-chornovol/lets-go-chat/database"
)

type Token struct {
	Id string
	UserId string
	ExpiresAt string
}

func (t Token) IsEmpty() bool {
	return t == Token{}
}

func CreateToken(userId string) (string, error) {
	token := map[string]string{
		"id": uuid.New().String(),
		"user_id": userId,
		"expires_at": time.Now().Add(time.Minute).String(),
	}
	err := database.Driver.Insert("tokens", token)

	return token["id"], err
}

func GetTokenById(id string) (Token, error) {
	var token Token

	whereAttributes := map[string]string{"id": id}
	tokens, err := database.Driver.Select("tokens", whereAttributes)
	if err != nil {
		return token, err
	}

	for tokens.Next() {
		err := tokens.Scan(&token.Id, &token.UserId, &token.ExpiresAt)
		if err != nil {
			log.Println(err)
		}
	}

	return token, nil
}
