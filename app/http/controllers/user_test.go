package controllers

import (
	"errors"
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
	mocksmodels "github.com/oleksandr-chornovol/lets-go-chat/mocks/app/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	cases := map[string]struct {
		requestBody *strings.Reader
		expectedResponseCode int
		expectedResponseBody string
		setupUserModelMock func(userModelMock *mocksmodels.UserInterface)
	}{
		"success registration": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusCreated,
			expectedResponseBody: `{"id":"user_id","name":"name"}` + "\n",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				userModelMock.On("GetUserByField", "name", "name").
					Return(models.User{}, nil)
				createUserArgument := models.User{Name: "name", Password: "password"}
				createUserResult := models.User{Id: "user_id", Name: "name", Password: "password"}
				userModelMock.On("CreateUser", createUserArgument).
					Return(createUserResult, nil)
			},
		},
		"invalid data": {
			requestBody: strings.NewReader(`{"name":"","password":""}`),
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: "Request data is invalid.",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {},
		},
		"name is already taken": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusConflict,
			expectedResponseBody: "Name is already taken.",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				user := models.User{Id: "user_id", Name: "name", Password: "password"}
				userModelMock.On("GetUserByField", "name", "name").
					Return(user, nil)
			},
		},
		"error in GetUserByField": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: "",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				userModelMock.On("GetUserByField", "name", "name").
					Return(models.User{}, errors.New("cannot get user"))
			},
		},
		"error in CreateUser": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: "",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				userModelMock.On("GetUserByField", "name", "name").
					Return(models.User{}, nil)
				createUserArgument := models.User{Name: "name", Password: "password"}
				createUserResult := models.User{Id: "user_id", Name: "name", Password: "password"}
				userModelMock.On("CreateUser", createUserArgument).
					Return(createUserResult, errors.New("cannot create user"))
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, "/v1/user", c.requestBody)
			assert.Nil(t, err)

			userModelMock := new(mocksmodels.UserInterface)
			c.setupUserModelMock(userModelMock)

			userController := UserController{
				TokenModel: models.Token{},
				UserModel:  userModelMock,
			}

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(userController.Register)
			handler.ServeHTTP(response, request)

			assert.Equal(t, c.expectedResponseCode, response.Code)
			assert.Equal(t, c.expectedResponseBody, response.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	cases := map[string]struct {
		requestBody *strings.Reader
		expectedResponseCode int
		expectedResponseBody string
		setupUserModelMock func(userModelMock *mocksmodels.UserInterface)
		setupTokenModelMock func(tokenModelMock *mocksmodels.TokenInterface)
	}{
		"success login": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"url":"ws:///v1/chat?token=token_id"}` + "\n",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				getUserByFieldResult := models.User{Id: "user_id", Name: "name", Password: "$2a$10$Kt0YB3SgXJuUpek5anTDguHyKXUEbE4EIyzQXrfzYzsNB9ExZflSe"}
				userModelMock.On("GetUserByField", "name", "name").
					Return(getUserByFieldResult, nil)
			},
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {
				tokenModelMock.On("CreateToken", models.Token{UserId: "user_id"}).
					Return(models.Token{Id: "token_id"}, nil)
			},
		},
		"incorrect password": {
			requestBody: strings.NewReader(`{"name":"name","password":"incorrect_password"}`),
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: "Password is incorrect.",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				getUserByFieldResult := models.User{Id: "user_id", Name: "name", Password: "$2a$10$Kt0YB3SgXJuUpek5anTDguHyKXUEbE4EIyzQXrfzYzsNB9ExZflSe"}
				userModelMock.On("GetUserByField", "name", "name").
					Return(getUserByFieldResult, nil)
			},
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {},
		},
		"user does not exist": {
			requestBody: strings.NewReader(`{"name":"name","password":"incorrect_password"}`),
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: "User does not exist.",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				userModelMock.On("GetUserByField", "name", "name").
					Return(models.User{}, nil)
			},
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {},
		},
		"error in GetUserByField": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: "",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				userModelMock.On("GetUserByField", "name", "name").
					Return(models.User{}, errors.New("cannot get user"))
			},
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {},
		},
		"error in CreateToken": {
			requestBody: strings.NewReader(`{"name":"name","password":"password"}`),
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: "",
			setupUserModelMock: func(userModelMock *mocksmodels.UserInterface) {
				getUserByFieldResult := models.User{Id: "user_id", Name: "name", Password: "$2a$10$Kt0YB3SgXJuUpek5anTDguHyKXUEbE4EIyzQXrfzYzsNB9ExZflSe"}
				userModelMock.On("GetUserByField", "name", "name").
					Return(getUserByFieldResult, nil)
			},
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {
				tokenModelMock.On("CreateToken", models.Token{UserId: "user_id"}).
					Return(models.Token{}, errors.New("cannot create token"))
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, "/v1/user/login", c.requestBody)
			assert.Nil(t, err)

			tokenModelMock := new(mocksmodels.TokenInterface)
			c.setupTokenModelMock(tokenModelMock)

			userModelMock := new(mocksmodels.UserInterface)
			c.setupUserModelMock(userModelMock)

			userController := UserController{
				TokenModel: tokenModelMock,
				UserModel:  userModelMock,
			}

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(userController.Login)
			handler.ServeHTTP(response, request)

			assert.Equal(t, c.expectedResponseCode, response.Code)
			assert.Equal(t, c.expectedResponseBody, response.Body.String())
		})
	}
}

func TestGetUserData(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"name","password":"password"}`))
	assert.Nil(t, err)

	userData := getUserData(request)

	assert.Equal(t, "name", userData.Name)
	assert.Equal(t, "password", userData.Password)
}
