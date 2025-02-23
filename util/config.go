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
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/pingcap/errors"
)

var GConfig *Config
var SM sync.Mutex

type Output struct {
	Path       string `toml:"path"`
	Fileprefix string `toml:"fileprefix"`
}

type Tables struct {
	Tables        string `toml:"tables"`
	IgnoreColumns string `toml:"ignorecolumns"`
}

type Join struct {
	Relationship string `toml:"relationship"`
}

type Check struct {
	Rule string `toml:"rule"`
}

type Base struct {
	Table          string `toml:"table"`
	Rowcount       string `toml:"rowcount"`
	Peerfilesize   string `toml:"peerfilesize"`
	Peerfilenum    string `toml:"peerfilenum"`
	Threadpoolsize string `toml:"threadpoolsize"`
	DSN            string `toml:"dsn"`
	LineTerm       string `toml:"lineterm"`
	FieldTerm      string `tome:"fieldterm"`
}
type Tomels struct {
	Output *Output
	Tables *Tables
	Join   *Join
	Check  *Check
	Base   *Base
}

func NewConfig(filename string) error {
	fmt.Println(filename)
	t := &Tomels{
		Output: &Output{},
		Tables: &Tables{},
		Join:   &Join{},
		Check:  &Check{},
		Base:   &Base{},
	}
	err := t.parseConfig(filename)
	if err != nil {
		return err
	}
	var c Config
	err = c.convertTomelsToConfig(t)
	if err != nil {
		return err
	}
	SM.Lock()
	defer SM.Unlock()
	GConfig = &c
	return nil
}

func (t *Tomels) parseConfig(filename string) error {
	if _, err := toml.DecodeFile(filename, t); err != nil {
		return err
	}
	fmt.Println(t.Base, t.Output, t.Tables, t.Check, t.Join)
	return nil
}

type JoinConditionst struct {
	ColName string
	Rate    int
}

type Config struct {
	OutPutPath string
	Fileprefix string
	IgnoreCols map[string][]string
	Tables     map[string]bool
	Joins      map[string][]JoinConditionst
	Checks     map[string][]string
	Base       map[string]string
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

func parseBaseInfo(b map[string]string, bi *Base) error {
	retT := reflect.TypeOf(*bi)
	retV := reflect.ValueOf(*bi)
	//获取结构体里的名称级值
	for i := 0; i < retT.NumField(); i++ {
		field := retT.Field(i)
		b[field.Name] = retV.FieldByName(field.Name).String()
	}
	return nil
}

func (c *Config) GetMaxFileSize() (uint64, error) {
	return strconv.ParseUint(c.Base["Peerfilesize"], 10, 64)
}

func (c *Config) GetfilePrefix() string {
	return c.Fileprefix
}

func (c *Config) GetOutputfile() string {
	return c.OutPutPath
}

func (c *Config) GetMaxFileNum() (uint64, error) {
	return strconv.ParseUint(c.Base["Peerfilenum"], 10, 64)
}

func (c *Config) GetDSN() string {
	return c.Base["DSN"]
}

func (c *Config) GetRowcount() (uint64, error) {
	return strconv.ParseUint(c.Base["Rowcount"], 10, 64)
}

func (c *Config) GetThreadPoolSize() (int, error) {
	return strconv.Atoi(c.Base["Threadpoolsize"])
}

func (c *Config) GetFieldTerm() string {
	return c.Base["FieldTerm"]
}

func (c *Config) GetLineTerm() string {
	return c.Base["LineTerm"]
}
func (c *Config) GetTables() string {
	var table_name string
	for k := range c.Tables {
		table_name = table_name + k + ","
	}
	return table_name
}

func (c *Config) convertTomelsToConfig(t *Tomels) error {

	mtables := make(map[string]bool)
	mignorecols := make(map[string][]string)
	joins := make(map[string][]JoinConditionst)
	checks := make(map[string][]string)
	bases := make(map[string]string)
	var err error

	c.OutPutPath = t.Output.Path
	c.Fileprefix = t.Output.Fileprefix
	if len(strings.TrimSpace(t.Tables.Tables)) == 0 {
		return errors.New("no specific tables")
	}
	tables := strings.Split(t.Tables.Tables, ",")
	for _, tablename := range tables {
		tbname := strings.ToLower(tablename)
		mtables[tbname] = true
	}

	fmt.Println("tables is :", mtables)

	if len(strings.TrimSpace(t.Tables.IgnoreColumns)) > 0 {
		cols := strings.Split(t.Tables.IgnoreColumns, ",")
		fmt.Println("cols is :", cols)
		for _, col := range cols {
			tbname, isExists := getColTables(mtables, col)
			if !isExists {
				return errors.New(fmt.Sprintf("could not find %v 's table", col))
			}
			if v, ok := mignorecols[tbname]; !ok {
				columns := make([]string, 1)
				columns[0] = strings.ToLower(col)
				mignorecols[tbname] = columns
			} else {
				v = append(v, strings.ToLower(col))
				mignorecols[tbname] = v
			}
		}
	}

	if len(strings.TrimSpace(t.Join.Relationship)) > 0 {
		relationships := strings.Split(t.Join.Relationship, ",")
		var rate int
		for _, v := range relationships {
			fmt.Println(v)
			pair := strings.Split(v, "/")
			if len(pair) < 2 || len(pair) > 3 {
				fmt.Println(pair)
				return errors.New("invalid relationship")
			}
			key := strings.ToLower(strings.TrimSpace(pair[0]))
			val := strings.ToLower(strings.TrimSpace(pair[1]))
			if len(pair) == 3 {
				rate, err = strconv.Atoi(strings.TrimSpace(pair[2]))
				if err != nil {
					return errors.New(fmt.Sprintf("get rate fail ,can not convert %v to int ,%v ", pair[2], err.Error()))
				}
			} else {
				rate = 1
			}
			vv, ok := joins[key]
			if !ok {
				s := make([]JoinConditionst, 1)
				s[0] = JoinConditionst{val, rate}
				joins[key] = s
			} else {
				vv := append(vv, JoinConditionst{val, rate})
				joins[key] = vv
			}
		}
	}
	if len(strings.TrimSpace(t.Check.Rule)) > 0 {
		check := strings.Split(t.Check.Rule, ";")
		for _, v := range check {
			col := strings.Split(v, ":")
			if len(col) != 2 {
				return errors.New("invalid check")
			} else {
				key_part := strings.Split(strings.TrimSpace(col[0]), ".")
				fmt.Println(col)
				if len(key_part) != 3 {
					return errors.New("invalid columns check")
				}
				key := fmt.Sprintf("%v.%v.%v", strings.ToLower(key_part[0]), strings.ToLower(key_part[1]), key_part[2])
				//vals := strings.Split(col[1], ",")
				//The second rule overrides the previous one
				checks[key] = []string{strings.TrimSpace(col[1])}
			}

		}
	}

	err = parseBaseInfo(bases, t.Base)
	if err != nil {
		return err
	}
	c.Base = bases
	c.IgnoreCols = mignorecols
	c.Tables = mtables
	c.Joins = joins
	c.Checks = checks
	fmt.Println(mtables)
	fmt.Println(mignorecols)
	return nil
}
