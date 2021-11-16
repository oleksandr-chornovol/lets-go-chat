package user

import (
	"encoding/json"
	"fmt"
	userModel "github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"log"
	"net/http"
	"pkg/hasher"
)

func Register(response http.ResponseWriter, request *http.Request) {
	userData := getUserData(request)

	if len(userData.Name) < 4 || len(userData.Password) < 8 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Request data is invalid.")
		return
	}

	_, userExists := userModel.GetUserByName(userData.Name)
	if userExists {
		response.WriteHeader(http.StatusConflict)
		fmt.Fprint(response, "Name is already taken.")
	} else {
		user := userModel.CreateUser(userData)
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(map[string]string{
			"id": user.Id,
			"name": user.Name,
		})
	}
}

func Login(response http.ResponseWriter, request *http.Request) {
	userData := getUserData(request)
	user, userExists := userModel.GetUserByName(userData.Name)

	if userExists {
		if hasher.CheckPasswordHash(userData.Password, user.Password) {
			json.NewEncoder(response).Encode(map[string]string{
				"url": request.Host + "/ws%token=one-time-token",
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

func getUserData(request *http.Request) userModel.User {
	decoder := json.NewDecoder(request.Body)
	var userData userModel.User
	err := decoder.Decode(&userData)
	if err != nil {
		log.Println(err)
	}
	return userData
}
