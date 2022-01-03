package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"github.com/oleksandr-chornovol/lets-go-chat/database/drivers"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	database.SetDriver(drivers.MySqlDriver{DB: db})

	dbMock.ExpectExec("insert into tokens").
		WillReturnResult(sqlmock.NewResult(1, 1))

	userId := "user_id"
	result, err := Token{}.CreateToken(Token{UserId: userId})
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, userId, result.UserId)
}

func TestGetTokenById(t *testing.T) {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer db.Close()
	database.SetDriver(drivers.MySqlDriver{DB: db})

	columns := []string{"id", "user_id", "expires_at"}
	token := Token{Id: "token_id", UserId: "user_id", ExpiresAt: time.Now().Add(time.Minute).String()}
	dbMock.ExpectQuery("select * from tokens where id = ?").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(token.Id, token.UserId, token.ExpiresAt))

	result, _ := Token{}.GetTokenById("token_id")
	assert.Nil(t, err)

	assert.Equal(t, token.Id, result.Id)
	assert.Equal(t, token.UserId, result.UserId)
	assert.Equal(t, token.ExpiresAt, result.ExpiresAt)
}

func TestTokenIsEmpty(t *testing.T) {
	result := Token{}.IsEmpty()

	assert.Equal(t, true, result)
}
