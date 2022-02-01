package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/hasher"

	"github.com/oleksandr-chornovol/lets-go-chat/app"
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
)

type UserController struct {
	TokenModel models.TokenInterface
	UserModel  models.UserInterface
}

func NewUserController(tokenModel models.TokenInterface, userModel models.UserInterface) *UserController {
	return &UserController{
		TokenModel: tokenModel,
		UserModel:  userModel,
	}
}

func (c *UserController) Register(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	userData := &app.CreateUserRequest{}
	err := json.NewDecoder(request.Body).Decode(userData)
	if err != nil {
		log.Println("json decode failed, err:", err)
	}

	if len(userData.Name) < 4 || len(userData.Password) < 8 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "Request data is invalid.")
		return
	}

	_, err = c.UserModel.GetUserByField("name", userData.Name)
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

	user, err := c.UserModel.CreateUser(models.User{
		Name: userData.Name,
		Password: userData.Password,
	})

	if err != nil {
		log.Println("CreateUser failed, err:", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(app.CreateUserResponse{
		Id: user.Id,
		Name: user.Name,
	})
}

func (c *UserController) Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	userData := &app.LoginUserRequest{}
	err := json.NewDecoder(request.Body).Decode(userData)
	if err != nil {
		log.Println("json decode failed, err:", err)
	}

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

		json.NewEncoder(response).Encode(app.LoginUserResonse{
			Url: "ws://" + request.Host + "/v1/chat?token=" + token.Id,
		})
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(response, "Password is incorrect.")
	}
}
