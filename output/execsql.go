/**
 * @Author: guobob
 * @Description:
 * @File:  execsql.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:57
 */

package output

import (
	"github.com/generate_data/sql"
	"github.com/go-sql-driver/mysql"
	"sync"
)

var stmtID uint64 = 0
var mu sync.Mutex

func getStmtID() uint64 {
	mu.Lock()
	defer mu.Unlock()
	stmtID = stmtID + 1
	return stmtID
}

type TableSql struct {
	//dbName      string
	//tableName   string
	cfg         *mysql.Config
	dsn         string
	handleStmts map[string]uint64
	params      []interface{}
	handle      sql.SQLHandle
}

func (ts *TableSql) WriteData(sql, dbName, tableName string, params []interface{}) error {
	var err error

	return err
}

func (ts *TableSql) Sync() error {
	query := "commit;"
	return ts.handle.Execute(query)
}

func (ts *TableSql) Close() error {
	ts.handle.Quit(false)
	return nil
}
