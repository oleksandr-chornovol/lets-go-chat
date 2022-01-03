package drivers

import "database/sql"

type DBDriverInterface interface {
	Insert(table string, attributes map[string]string) error
	SelectRow(table string, attributes map[string]string) *sql.Row
	//Update(table string, attributes map[string]string)
	//Delete(table string, attributes map[string]string)
}