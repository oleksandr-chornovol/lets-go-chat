package models

import (
	"github.com/google/uuid"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"time"
)

type TokenInterface interface {
	CreateToken(token Token) (Token, error)
	GetTokenById(id string) (Token, error)
}

type Token struct {
	Id string
	UserId string
	ExpiresAt string
}

func (t *Token) CreateToken(token Token) (Token, error) {
	token.Id = uuid.New().String()
	token.ExpiresAt = time.Now().Add(time.Minute).String()

	attributes := map[string]string{
		"id": token.Id,
		"user_id": token.UserId,
		"expires_at": token.ExpiresAt,
	}
	err := database.Driver.Insert("tokens", attributes)

	return token, err
}

func (t *Token) GetTokenById(id string) (Token, error) {
	var whereAttributes = [][3]string{
		{"id", "=", id},
	}
	result := database.Driver.SelectRow("tokens", whereAttributes)

	var token Token
	err := result.Scan(&token.Id, &token.UserId, &token.ExpiresAt)

	return token, err
}
