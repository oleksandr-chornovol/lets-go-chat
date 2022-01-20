package drivers

import "database/sql"

type DBDriverInterface interface {
	Insert(table string, attributes map[string]string) error
	Select(table string, attributes [][3]string) (*sql.Rows, error)
	SelectRow(table string, attributes [][3]string) *sql.Row
	Update(table string, whereAttributes map[string]string, updateAttributes map[string]string) error
	//Delete(table string, attributes map[string]string)
}
