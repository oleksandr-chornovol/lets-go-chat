package drivers

import (
	"database/sql"
	"log"
)

type MySqlDriver struct {
	DB *sql.DB
}

func (d MySqlDriver) Select(table string, attributes map[string]string) *sql.Rows {
	query := "select * from " + table
	var values []interface{}

	if len(attributes) > 0 {
		query += " where "
		for column, value := range attributes {
			query += column + " = ?,"
			values = append(values, value)
		}
		query = query[:len(query) - 1]
	}

	result, err := d.DB.Query(query, values...)
	if err != nil {
		log.Println(err)
	}

	return result
}

func (d MySqlDriver) Insert(table string, attributes map[string]string) {
	query := "insert into " + table + " ("
	columns := ""
	valuesPlaceholders := ""
	var values []interface{}

	for column, value := range attributes {
		columns += column + ", "
		valuesPlaceholders += "?, "
		values = append(values, value)
	}
	columns = columns[:len(columns) - 2]
	valuesPlaceholders = valuesPlaceholders[:len(valuesPlaceholders) - 2]
	query += columns + ") values (" + valuesPlaceholders + ")"

	_, err :=  d.DB.Exec(query, values...)
	if err != nil {
		log.Println(err)
	}
}
