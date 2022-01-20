package drivers

import (
	"database/sql"
)

type MySqlDriver struct {
	DB *sql.DB
}

func (d MySqlDriver) Select(table string, attributes [][3]string) (*sql.Rows, error) {
	query, values := getSelectQuery(table, attributes)

	return d.DB.Query(query, values...)
}

func (d MySqlDriver) SelectRow(table string, attributes [][3]string) *sql.Row {
	query, values := getSelectQuery(table, attributes)

	return d.DB.QueryRow(query, values...)
}

func (d MySqlDriver) Insert(table string, attributes map[string]string) error {
	query := "insert into " + table + " ("
	columns := ""
	valuesPlaceholders := ""
	var values []interface{}

	for column, value := range attributes {
		columns += column + ", "
		valuesPlaceholders += "?, "
		values = append(values, value)
	}
	columns = columns[:len(columns)-2]
	valuesPlaceholders = valuesPlaceholders[:len(valuesPlaceholders)-2]
	query += columns + ") values (" + valuesPlaceholders + ")"

	_, err := d.DB.Exec(query, values...)

	return err
}

func (d MySqlDriver) Update(table string, whereAttributes map[string]string, updateAttributes map[string]string) error {
	query := "update " + table + " set "
	var values []interface{}

	for column, value := range updateAttributes {
		query += column + " = ?, "
		values = append(values, value)
	}
	query = query[:len(query)-2] + " where "

	for column, value := range whereAttributes {
		query += column + " = ?, "
		values = append(values, value)
	}
	query = query[:len(query)-2]

	_, err := d.DB.Exec(query, values...)

	return err
}

func getSelectQuery(table string, attributes [][3]string) (string, []interface{}) {
	query := "select * from " + table
	var values []interface{}

	if len(attributes) > 0 {
		query += " where "
		for _, attribute := range attributes {
			query += attribute[0] + " " + attribute[1] + " ?,"
			values = append(values, attribute[2])
		}
		query = query[:len(query)-1]
	}

	return query, values
}
