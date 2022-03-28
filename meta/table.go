/**
 * @Author: guobob
 * @Description:
 * @File:  table.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:58
 */

package meta

import (
	"fmt"
	"github.com/generate_data/util"
)

type Table struct {
	TableID         uint
	TableName       string
	DBName          string
	Columns         []Column
	PersistenceType int //0:output ;1:database
	PrepareSQL      string
	Record          string
	FiledTerminate  string
	LineTerminate   string
}

func NewTable(tableName, dbName string) *Table {
	return &Table{
		TableID:   util.GetTableID(),
		DBName:    dbName,
		TableName: tableName,
	}
}

func (t *Table) GeneratePrepareSQL() {
	sql := "insert into "
	key := fmt.Sprintf("%v.%v", t.DBName, t.TableName)
	sql = sql + key + " ("
	valstr := " "
	i := 0
	for ; i < len(t.Columns)-1; i++ {
		sql = sql + t.Columns[i].ColumnName + ","
		valstr = "?,"
	}
	sql = sql + t.Columns[len(t.Columns)-1].ColumnName + ") values ("
	valstr = "?"
	sql = sql + valstr
	t.PrepareSQL = sql

}
