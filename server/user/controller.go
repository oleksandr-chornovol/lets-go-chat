package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pkg/hasher"
	"strconv"
)

func Register(response http.ResponseWriter, request *http.Request) {
	userData := getUserData(request)

	if len(userData.Name) < 4 || len(userData.Password) < 8 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Request data is invalid.")
		return
	}

	_, userExists := GetUserByName(userData.Name)
	if userExists {
		response.WriteHeader(http.StatusConflict)
		fmt.Fprint(response, "Name is already taken.")
	} else {
		user := CreateUser(userData)
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(map[string]string{
			"id": strconv.Itoa(user.Id),
			"name": user.Name,
		})
	}
}

func Login(response http.ResponseWriter, request *http.Request) {
	userData := getUserData(request)
	user, userExists := GetUserByName(userData.Name)

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

func getUserData(request *http.Request) User {
	decoder := json.NewDecoder(request.Body)
	var userData User
	err := decoder.Decode(&userData)
	if err != nil {
		panic(err)
	}
	return userData
}
