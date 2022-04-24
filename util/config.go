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
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/pingcap/errors"
)

var GConfig *Config
var SM sync.Mutex

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

func NewConfig(filename string) error {
	var t Tomels
	err := t.ParseConfig(filename)
	if err != nil {
		return err
	}
	var c Config
	err = c.ConvertTomelsToConfig(&t)
	if err != nil {
		return err
	}
	SM.Lock()
	defer SM.Unlock()
	GConfig = &c
	return nil
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
	Cols       map[string][]string
	Tables     map[string]bool
	Joins      map[string][]string
	Checks     map[string][]string
}

func getColTables(mtables map[string]bool, col string) (string, bool) {
	colstr := strings.ToLower(col)
	var b = false
	var tablename string
	for k := range mtables {
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
	joins := make(map[string][]string)
	checks := make(map[string][]string)

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

	relationships := strings.Split(t.Join.relationship, ",")
	for _, v := range relationships {
		pair := strings.Split(v, "/")
		if len(pair) != 2 {
			return errors.New("invalid relationship")
		}
		key := strings.ToLower(strings.TrimSpace(pair[0]))
		val := strings.ToLower(strings.TrimSpace(pair[1]))
		vv, ok := joins[key]
		if !ok {
			s := make([]string, 1)
			s[0] = val
			joins[key] = s
		} else {
			vv := append(vv, val)
			joins[key] = vv
		}
	}
	check := strings.Split(t.Check.rule, ";")
	for _, v := range check {
		col := strings.Split(v, ":")
		if len(col) != 2 {
			return errors.New("invalid check")
		} else {
			key := strings.ToLower(strings.TrimSpace(col[0]))
			vals := strings.Split(col[1], ",")
			//The second rule overrides the previous one
			checks[key] = vals
		}

	}
	c.Cols = mcols
	c.Tables = mtables
	c.Joins = joins
	c.Checks = checks
	fmt.Println(mtables)
	fmt.Println(mcols)
	return nil
}
