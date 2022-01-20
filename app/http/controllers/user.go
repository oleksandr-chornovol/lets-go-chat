package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
)

type UserController struct {
	TokenModel models.TokenInterface
	UserModel  models.UserInterface
}

func (c *UserController) Register(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userData := getUserData(request)

	if len(userData.Name) < 4 || len(userData.Password) < 8 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Request data is invalid.")
		return
	}

	_, err := c.UserModel.GetUserByField("name", userData.Name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("GetUserByField failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		response.WriteHeader(http.StatusConflict)
		fmt.Fprint(response, "Name is already taken.")
		return
	}

	user, err := c.UserModel.CreateUser(userData)
	if err != nil {
		log.Println("CreateUser failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(map[string]string{
		"id":   user.Id,
		"name": user.Name,
	})
}

func (c *UserController) Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	userData := getUserData(request)
	user, err := c.UserModel.GetUserByField("name", userData.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(response, "User does not exist.")
			return
		} else {
			log.Println("GetUserByField failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if hasher.CheckPasswordHash(userData.Password, user.Password) {
		token, err := c.TokenModel.CreateToken(models.Token{UserId: user.Id})
		if err != nil {
			log.Println("CreateToken failed, err:", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(response).Encode(map[string]string{
			"url": "ws://" + request.Host + "/v1/chat?token=" + token.Id,
		})
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(response, "Password is incorrect.")
	}
}

func getUserData(request *http.Request) models.User {
	decoder := json.NewDecoder(request.Body)
	var userData models.User
	err := decoder.Decode(&userData)
	if err != nil {
		log.Println("Decode failed, err:", err)
	}
	return userData
}
