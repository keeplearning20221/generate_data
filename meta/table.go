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
	if len(t.Columns) == 0 {
		t.PrepareSQL = ""
		return
	}
	fmt.Println(len(t.Columns))
	sql := "insert into "
	key := fmt.Sprintf("%v.%v", t.DBName, t.TableName)
	sql = sql + key + " ("
	valstr := " "
	i := 0
	for ; i < len(t.Columns)-1; i++ {
		sql = sql + t.Columns[i].ColumnName + ","
		valstr = valstr + "?,"
	}
	sql = sql + t.Columns[len(t.Columns)-1].ColumnName + ") values ("
	valstr = valstr + "?);"
	sql = sql + valstr
	t.PrepareSQL = sql

}

func (table *Table) Generate_table_data() (buff []byte, err error) {
	var out string

	columnslen := len(table.Columns)
	fmt.Println(columnslen)
	for i := 0; i < columnslen; i++ {

		str, err := table.Columns[i].Property.GenerateColumnData()
		if err != nil {
			return nil, err
		}
		out = out + str
		out = out + ","

	}
	out = out[:len(out)-1] + "\n"
	return []byte(out), nil
}

func Generate_tables_data(gmeta *map[string]Table) (err error) {
	for table_name, table := range Gmeta {
		fmt.Println(table_name)
		for i := 0; i < 10; i++ {
			out, err := table.Generate_table_data()
			if err == nil {
				fmt.Printf(string(out))
			} else {
				return err

			}

		}
	}

	return nil
}

func (t *Table) GenerateRecordData() (string, error) {
	record := ""
	i := 0
	for ; i < len(t.Columns)-1; i++ {
		val, err := t.Columns[i].GenerateColumnData()
		if err != nil {
			return record, err
		}
		record = record + val + t.FiledTerminate
	}
	val, err := t.Columns[i].GenerateColumnData()
	if err != nil {
		return record, err
	}
	record = record + val + t.LineTerminate
	return record, nil

}
