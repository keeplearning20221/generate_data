/**
 * @Author: guobob
 * @Description:
 * @File:  getmeta.go
 * @Version: 1.0.0
 * @Date: 2022/3/25 09:56
 */

package meta

import (
	"fmt"
	"github.com/generate_data/sql"
	"github.com/generate_data/util"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"sync"
)

var Gmeta map[string]Table
var Gmu sync.RWMutex

func init() {
	Gmeta = make(map[string]Table)
}

func GetTableInfo(s string, dsn string, cfg *mysql.Config, log *zap.Logger) error {
	//get table name from config string
	handle := sql.NewSQLHandle(dsn, cfg)
	err := handle.HandShake(cfg.DBName)
	if err != nil {
		return err
	}

	tables, err := util.GetTableName(s)
	if err != nil {
		return err
	}

	fmt.Println("get table name is ", tables)
	for _, v := range tables {
		err = util.CheckTableValid(v)
		if err != nil {
			return err
		}
		ss, err := util.SpiltTableName(v)
		if err != nil {
			return err
		}
		fmt.Println(ss)
		err = sql.GetColumnInfo(handle, ss[0], ss[1])
		if err != nil {
			log.Error("get column info fail ," + err.Error())
			return err
		}
		table := &Table{
			TableID:   util.GetTableID(),
			TableName: ss[1],
			DBName:    ss[0],
		}
		err = GetColumnFromMetaData(handle, table)
		if err != nil {
			log.Error("convert column info fail," + err.Error())
			return err
		}
		Gmeta[v] = *table
	}
	return nil
}
