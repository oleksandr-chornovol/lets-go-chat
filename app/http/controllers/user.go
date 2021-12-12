package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
)

func Register(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userData := getUserData(request)

	if len(userData.Name) < 4 || len(userData.Password) < 8 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Request data is invalid.")
		return
	}

	user, err := models.GetUserByName(userData.Name)
	if err != nil {
		log.Println("GetUserByName failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.IsEmpty() {
		user, err := models.CreateUser(userData)
		if err == nil {
			response.WriteHeader(http.StatusCreated)
			json.NewEncoder(response).Encode(map[string]string{
				"id":   user.Id,
				"name": user.Name,
			})
		} else {
			log.Println(err)
			response.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		response.WriteHeader(http.StatusConflict)
		fmt.Fprint(response, "Name is already taken.")
	}
}

func Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userData := getUserData(request)
	user, err := models.GetUserByName(userData.Name)
	if err != nil {
		log.Println("GetUserByName failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ! user.IsEmpty() {
		if hasher.CheckPasswordHash(userData.Password, user.Password) {
			token, err := models.CreateToken(user.Id)
			if err != nil {
				log.Println("CreateToken failed, err:", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(response).Encode(map[string]string {
				"url": "ws://" + request.Host + "/v1/chat?token=" + token,
			})
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(response, "Password is incorrect.")
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(response, "User does not exist.")
	}
}

func getUserData(request *http.Request) models.User {
	decoder := json.NewDecoder(request.Body)
	var userData models.User
	err := decoder.Decode(&userData)
	if err != nil {
		log.Println(err)
	}
	return userData
}
