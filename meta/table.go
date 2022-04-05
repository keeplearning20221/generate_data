/**
 * @Author: guobob
 * @Description:
 * @File:  table.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:58
 */

package meta

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/generate_data/util"
	deci "github.com/shopspring/decimal"
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

func Generate_table_data(table Table) (buff []byte, err error) {
	var out string

	columnslen := len(table.Columns)
	fmt.Println(columnslen)
	for i := 0; i < columnslen; i++ {

		str, err := util.GenerateData(table.Columns[i].Property)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("第%d行数据:", i)
		switch str := str.(type) {
		case string:
			//fmt.Printf(str)
			out = out + str
			out = out + ","
		case int64:
			out = out + strconv.FormatInt(str, 10)
			out = out + ","
		case deci.Decimal:
			out = out + str.String()
			out = out + ","
		default:
			fmt.Println(str)
			err := errors.New("unkown str type")
			return nil, err
		}

	}
	out = out[:len(out)-1] + "\n"
	return []byte(out), nil
}

func Generate_tables_data(gmeta *map[string]Table) (err error) {
	for table_name, table := range Gmeta {
		fmt.Println(table_name)
		for i := 0; i < 10; i++ {
			out, err := Generate_table_data(table)
			if err == nil {
				fmt.Printf(string(out))
			} else {
				return err
			}

		}
	}
	return nil
}
