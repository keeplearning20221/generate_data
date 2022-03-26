/**
 * @Author: guobob
 * @Description:
 * @File:  table.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:58
 */

package meta

import "github.com/generate_data/util"

type Table struct {
	TableID         uint
	TableName       string
	DBName          string
	Columns         []Column
	PersistenceType int //0:file ;1:database
	PrepareSQL      string
	Record          string
	FiledTerminate  string
	LineTerminate   string
}

func NewTable(tableName,dbName string) *Table{
	return &Table{
		TableID: util.GetTableID(),
		DBName:dbName,
		TableName: tableName,
	}
}




func (t *Table) GeneratePrepareSQL() (string, error) {

	return "", nil
}
