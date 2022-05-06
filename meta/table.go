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

func NewTableWithTable(t *Table, fileNoInc, peerFileNum int64) *Table {
	cols := make([]Column, len(t.Columns))
	for k, v := range t.Columns {
		var p util.Property
		p.BitLen = v.BitLen
		p.CharFormat = make([]byte, len(v.CharFormat))
		copy(p.CharFormat, v.CharFormat)
		p.CharLen = v.CharLen
		p.DefaultVal = make([]string, len(v.DefaultVal))
		copy(p.DefaultVal, v.DefaultVal)
		p.EndKey = v.EndKey
		p.EndValue = v.EndValue
		p.Length = v.Length
		p.NowValue = v.NowValue
		p.StartKey = v.StartKey
		p.StartValue = v.StartValue + fileNoInc*peerFileNum
		p.SuffixLen = v.SuffixLen
		p.Type = v.Type
		p.TypeGen = v.TypeGen
		cols[k] = Column{
			&p,
			v.ColumnName,
			v.ColumnIdx,
			v.Ignore,
		}
	}
	return &Table{
		TableID:         t.TableID,
		TableName:       t.TableName,
		DBName:          t.DBName,
		Columns:         cols,
		PersistenceType: t.PersistenceType,
		PrepareSQL:      t.PrepareSQL,
		Record:          t.Record,
		FiledTerminate:  t.FiledTerminate,
		LineTerminate:   t.LineTerminate,
	}
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

// func (table *Table) Generate_table_data() (buff []byte, err error) {
// 	var out string

// 	columnslen := len(table.Columns)
// 	fmt.Println(columnslen)
// 	for i := 0; i < columnslen; i++ {

// 		str, err := table.Columns[i].Property.GenerateColumnData()
// 		if err != nil {
// 			return nil, err
// 		}
// 		out = out + str
// 		out = out + ","

// 	}
// 	out = out[:len(out)-1] + "\n"
// 	return []byte(out), nil
// }

// func Generate_tables_data(gmeta *map[string]Table) (err error) {
// 	for table_name, table := range Gmeta {
// 		fmt.Println(table_name)
// 		for i := 0; i < 10; i++ {
// 			out, err := table.Generate_table_data()
// 			if err == nil {
// 				fmt.Printf(string(out))
// 			} else {
// 				return err

// 			}

// 		}
// 	}

// 	return nil
// }

func (t *Table) GenerateRecordData(id uint64, increm_info []util.Incrementinfo) (string, error) {
	record := ""
	i := 0
	for ; i < len(t.Columns)-1; i++ {
		val, err := t.Columns[i].GenerateColumnData(&increm_info[i])
		if err != nil {
			return record, err
		}
		record = record + val + t.FiledTerminate
	}
	val, err := t.Columns[i].GenerateColumnData(&increm_info[i])
	if err != nil {
		return record, err
	}
	record = record + val + t.LineTerminate
	return record, nil

}
