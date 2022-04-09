/**
 * @Author: guobob
 * @Description:
 * @File:  common.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:53
 */

package util

import (
	"fmt"

	//"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
)

func CheckTableValid(tableName string) error {
	if len(tableName) > 128 {
		return errors.New(fmt.Sprintf("table name %v len large 128", tableName))
	}
	pos := strings.Index(tableName, ".")
	if pos == -1 || pos < 2 || pos == len(tableName) {
		return errors.New(fmt.Sprintf("table name invalid"))
	}
	return nil
}

func SpiltTableName(s string) ([]string, error) {
	ss := strings.Split(s, ".")
	if len(ss) != 2 {
		return nil, errors.New(fmt.Sprintf("%v is invalid ", s))
	}
	return ss, nil
}

//GetTableName : get table name from config string
func GetTableName(s string) ([]string, error) {
	var tables []string
	ss := strings.Split(s, ",")

	for _, v := range ss {
		table := strings.TrimSpace(v)
		tables = append(tables, table)
	}

	return tables, nil
}

func ParseDSN(dsn string) (*mysql.Config, error) {

	if len(dsn) == 0 {
		return nil, errors.New("parma dsn len is zero")
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
