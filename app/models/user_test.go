package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/oleksandr-chornovol/lets-go-chat/database"
	"github.com/oleksandr-chornovol/lets-go-chat/database/drivers"
)

func TestCreateUser(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	database.SetDriver(drivers.MySqlDriver{DB: db})

	dbMock.ExpectExec("insert into users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	userModel := User{}
	result, err := userModel.CreateUser(User{Name: "user_name", Password: "user_password"})
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "user_name", result.Name)
}

func TestGetUserByField(t *testing.T) {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer db.Close()
	database.SetDriver(drivers.MySqlDriver{DB: db})

	columns := []string{"id", "name", "password", "last_session_end"}
	user := User{Id: "user_id", Name: "user_name", Password: "hashed_password", LastSessionEnd: time.Now().Add(-time.Minute).String()}
	dbMock.ExpectQuery("select * from users where name = ?").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(user.Id, user.Name, user.Password, user.LastSessionEnd))

	userModel := User{}
	result, err := userModel.GetUserByField("name", user.Name)
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, user.Id, result.Id)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Password, result.Password)
	assert.Equal(t, user.LastSessionEnd, result.LastSessionEnd)
}
