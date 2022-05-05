/**
 * @Author: guobob
 * @Description:
 * @File:  const.go
 * @Version: 1.0.0
 * @Date: 2022/3/25 14:40
 */

package util

import "sync"

const (
	INT int = iota
	DECIMAL
	TIMESTAMP
	DATETIME
	STRING
	EUME
	STRINGCN
	UNKNOW = -1
)

const (
	VOLATILE int = iota
	UNVOLATILE
)

var TableID uint = 0
var mu sync.RWMutex

func GetTableID() uint {
	mu.Lock()
	defer mu.Unlock()
	TableID++
	return TableID
}

func ChangeColType(colType string) int {
	switch colType {
	case "int", "bigint":
		return INT
	case "decimal":
		return DECIMAL
	case "timestamp":
		return DATETIME
	case "datetime":
		return DATETIME
	case "char", "varchar":
		return STRING
	default:
		return UNKNOW
	}
}
