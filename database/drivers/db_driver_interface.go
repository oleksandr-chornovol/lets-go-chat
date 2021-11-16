package drivers

import "database/sql"

type DBDriverInterface interface {
	Select(table string, attributes map[string]string) *sql.Rows
	Insert(table string, attributes map[string]string)
	//Update(table string, attributes map[string]string)
	//Delete(table string, attributes map[string]string)
}