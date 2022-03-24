/**
 * @Author: guobob
 * @Description:
 * @File:  table.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:58
 */

package meta

import "sync"

var g_meta map[string]Table
var g_mu sync.RWMutex


type Table struct {
	TableID    int
	TableName  string
	DBName     string
	Columns    []Column
	PersistenceType int //0:file ;1:database
	PrepareSQL string
	Record     string
	FiledTerminate string
	LineTerminate  string
}


func (t *Table) GeneratePrepareSQL() (string, error) {

	return "", nil
}