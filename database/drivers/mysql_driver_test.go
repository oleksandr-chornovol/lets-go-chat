package drivers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectRow(t *testing.T) {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer db.Close()
	mySqlDriver := MySqlDriver{DB: db}

	columns := []string{"id", "name"}
	name := "name_value"
	dbMock.ExpectQuery("select * from entities where id = ?").
		WillReturnRows(sqlmock.NewRows(columns).AddRow("id_value", name))

	result := mySqlDriver.SelectRow("entities", [][3]string{{"id", "=", "id_value"}})

	var entity struct {
		Id   string
		Name string
	}
	err = result.Scan(&entity.Id, &entity.Name)
	assert.Nil(t, err)

	assert.Equal(t, name, entity.Name)
}

func TestInsert(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mySqlDriver := MySqlDriver{DB: db}

	dbMock.ExpectExec("insert into entities").
		WithArgs("id_value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = mySqlDriver.Insert("entities", map[string]string{"id": "id_value"})
	assert.Nil(t, err)
}
