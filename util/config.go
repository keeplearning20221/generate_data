/**
 * @Author: guobob
 * @Description:
 * @File:  const.go
 * @Version: 1.0.0
 * @Date: 2022/3/25 14:40
 */

package util

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pingcap/errors"
)

type OutputInfo struct {
	path       string
	fileprefix string
}

type BaseInfo struct {
	tables  string
	columns string
}

type JoinInfo struct {
	relationship string
}

type CheckInfo struct {
	rule string
}
type Tomels struct {
	Output OutputInfo
	Base   BaseInfo
	Join   JoinInfo
	Check  CheckInfo
}

func (t *Tomels) ParseConfig(filename string) error {
	if _, err := toml.DecodeFile(filename, t); err != nil {
		return err
	}
	return nil
}

type Config struct {
	OutPutPath string
	Fileprefix string
	Cols       map[string]string
}

func getColTables(mtables map[string]bool, col string) (string, bool) {
	colstr := strings.ToLower(col)
	var b = false
	var tablename string
	for k, _ := range mtables {
		b = strings.HasPrefix(colstr, k)
		if b {
			tablename = k
			break
		}
	}

	return tablename, b
}

func (c *Config) ConvertTomelsToConfig(t *Tomels) error {
	mtables := make(map[string]bool)
	mcols := make(map[string][]string)
	c.OutPutPath = t.Output.path
	c.Fileprefix = t.Output.fileprefix
	tables := strings.Split(t.Base.tables, ",")

	for _, tablename := range tables {
		tbname := strings.ToLower(tablename)
		mtables[tbname] = true
	}
	cols := strings.Split(t.Base.columns, ",")
	for _, col := range cols {
		tbname, isExists := getColTables(mtables, col)
		if !isExists {
			return errors.New(fmt.Sprintf("could not find %v 's table", col))
		}
		if v, ok := mcols[tbname]; !ok {
			columns := make([]string, 1)
			columns[0] = strings.ToLower(col)
			mcols[tbname] = columns
		} else {
			v = append(v, strings.ToLower(col))
			mcols[tbname] = v
		}
	}

	return nil
}
