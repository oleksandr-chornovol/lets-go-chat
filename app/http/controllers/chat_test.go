package controllers

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/oleksandr-chornovol/lets-go-chat/app/models"
	"github.com/oleksandr-chornovol/lets-go-chat/cache"
	mocksmodels "github.com/oleksandr-chornovol/lets-go-chat/mocks/app/models"
	mockscache "github.com/oleksandr-chornovol/lets-go-chat/mocks/cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStartEcho(t *testing.T) {
	tokenId := "token_id"

	tokenModelMock := new(mocksmodels.TokenInterface)
	getTokenByIdResult := models.Token{Id: tokenId, UserId: "user_id", ExpiresAt: time.Now().Add(time.Minute).String()}
	tokenModelMock.On("GetTokenById", tokenId).
		Return(getTokenByIdResult, nil)

	userModelMock := new(mocksmodels.UserInterface)
	getUserByFieldResult := models.User{Id: getTokenByIdResult.UserId, Name: "name", Password: "password"}
	userModelMock.On("GetUserByField", "id", getTokenByIdResult.UserId).
		Return(getUserByFieldResult, nil)

	activeUsersCacheMock := new(mockscache.ActiveUsersCacheInterface)
	activeUsersCacheMock.On("AddUser", getUserByFieldResult)
	activeUsersCacheMock.On("DeleteUser", getUserByFieldResult.Id)

	chatController := ChatController{
		ActiveUsersCache: activeUsersCacheMock,
		TokenModel:       tokenModelMock,
		UserModel:        userModelMock,
	}

	server := httptest.NewServer(getHandlerFunc(chatController.StartEcho))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + tokenId

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err)
	defer ws.Close()

	for i := 0; i < 10; i++ {
		message := "message"
		err := ws.WriteMessage(websocket.TextMessage, []byte(message))
		assert.Nil(t, err)

		_, receivedMessage, err := ws.ReadMessage()
		assert.Nil(t, err)

		assert.Equal(t, message, string(receivedMessage))
	}
}

func TestStartEchoNegativeCases(t *testing.T) {
	cases := map[string]struct {
		expectedResponseCode int
		expectedResponseBody string
		setupTokenModelMock func(tokenModelMock *mocksmodels.TokenInterface)
	}{
		"error in GetTokenById": {
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: "",
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {
				tokenModelMock.On("GetTokenById", "token_id").
					Return(models.Token{}, errors.New("cannot get token"))
			},
		},
		"token is expired": {
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: "Token is expired.",
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {
				expiredTime := time.Now().Add(-1 * time.Hour).String()
				token := models.Token{Id: "token_id", UserId: "user_id", ExpiresAt: expiredTime}
				tokenModelMock.On("GetTokenById", "token_id").
					Return(token, nil)
			},
		},
		"token does not exist": {
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: "Token does not exist.",
			setupTokenModelMock: func(tokenModelMock *mocksmodels.TokenInterface) {
				tokenModelMock.On("GetTokenById", "token_id").
					Return(models.Token{}, nil)
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/v1/chat?token=token_id", nil)
			assert.Nil(t, err)

			tokenModelMock := new(mocksmodels.TokenInterface)
			c.setupTokenModelMock(tokenModelMock)

			chatController := ChatController{
				ActiveUsersCache: cache.NewActiveUsersCache(),
				TokenModel:       tokenModelMock,
				UserModel:        models.User{},
			}

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(chatController.StartEcho)
			handler.ServeHTTP(response, request)

			assert.Equal(t, c.expectedResponseCode, response.Code)
			assert.Equal(t, c.expectedResponseBody, response.Body.String())
		})
	}
}

func TestGetActiveUsersCount(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/v1/user/active", nil)
	assert.Nil(t, err)

	activeUsersCacheMock := new(mockscache.ActiveUsersCacheInterface)
	getAllUsersResult := make(map[string]models.User)
	getAllUsersResult["user_id"] = models.User{}
	activeUsersCacheMock.On("GetAllUsers").
		Return(getAllUsersResult)

	chatController := ChatController {
		ActiveUsersCache: activeUsersCacheMock,
		TokenModel: models.Token{},
		UserModel: models.User{},
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(chatController.GetActiveUsersCount)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"count_of_users\":1}\n", response.Body.String())
}

func getHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r)
	}
}
