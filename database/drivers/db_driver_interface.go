package drivers

import "database/sql"

type DBDriverInterface interface {
	Select(table string, attributes map[string]string) (*sql.Rows, error)
	Insert(table string, attributes map[string]string) error
	//Update(table string, attributes map[string]string)
	//Delete(table string, attributes map[string]string)
}