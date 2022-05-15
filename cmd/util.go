package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/generate_data/meta"
	"github.com/generate_data/util"
	"github.com/go-sql-driver/mysql"
)

type baseConfig struct {
	dsn            string
	cfg            *mysql.Config
	tables         string
	fieldTerm      string
	lineTerm       string
	outputPath     string
	count          uint64
	filePrefix     string
	maxFileSize    uint64
	maxFileNum     uint64
	threadPoolSize int
}

func (bc *baseConfig) getVals() error {
	var err error

	bc.maxFileSize, err = util.GConfig.GetMaxFileSize()
	if err != nil {
		return err
	}

	bc.maxFileNum, err = util.GConfig.GetMaxFileNum()
	if err != nil {
		return err
	} else if bc.maxFileNum == 0 {
		//defaule value 10000
		bc.maxFileNum = 10000
	}

	bc.threadPoolSize, err = util.GConfig.GetThreadPoolSize()
	if err != nil {
		return err
	} else if bc.threadPoolSize == 0 {
		bc.threadPoolSize = 10
	}
	bc.dsn = util.GConfig.GetDSN()
	bc.cfg, err = util.ParseDSN(bc.dsn)
	if err != nil {
		return err
	}
	bc.filePrefix = util.GConfig.GetfilePrefix()
	bc.outputPath = util.GConfig.GetOutputfile()
	err = util.CheckDirValid(bc.outputPath)
	if err != nil {
		return err
	}
	bc.tables = util.GConfig.GetTables()
	if len(strings.TrimSpace(bc.tables)) == 0 {
		return errors.New("tables name  len is zero")
	}
	bc.count, err = util.GConfig.GetRowcount()
	if err != nil {
		return err
	}
	if bc.count == 0 {
		return errors.New("generate row count is zero")
	}
	bc.fieldTerm = util.GConfig.GetFieldTerm()
	if bc.fieldTerm == "" {
		bc.fieldTerm = ","
	}
	bc.lineTerm = util.GConfig.GetLineTerm()
	if bc.lineTerm == "" {
		bc.lineTerm = "\n"
	}

	return err
}

//@SetColGenerateDataFlag
//Flag whether the column needs to generate data
func setColGenerateDataFlag(wg *sync.WaitGroup) {
	defer wg.Done()
	for tableName, cols := range util.GConfig.IgnoreCols {
		t := meta.Gmeta[tableName]
		for _, v := range t.Columns {
			for _, v1 := range cols {
				if v.ColumnName == v1 {
					v.Ignore = true
				}
			}
		}
	}

}

//@SetColProperty
func setColProperty(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	fmt.Println(util.GConfig.Checks)
	for tableName, t := range meta.Gmeta {
		for _, v := range t.Columns {
			key := fmt.Sprintf("%v.%v", tableName, v.ColumnName)
			if check, ok := util.GConfig.Checks[key]; !ok {
				continue
			} else {
				fmt.Println(check[0])
				if strings.Contains(check[0], "~") {
					//fmt.Println("range columns is ", key)
					//check whether it is an interval
					//TODO : Handle errors
					vals := strings.Split(check[0], "~")
					if len(vals) != 2 {
						fmt.Println(" only one interval can be specified ")
					}
					v.TypeGen = 2 //range
					//v.DefaultVal = make([]string, 2)
					//v.DefaultVal[0] = vals[0]
					//v.DefaultVal[1] = vals[1]

					v.StartValue, err = strconv.ParseInt(vals[0], 10, 64)
					if err != nil {
						fmt.Println("convert start value to int64 fail ", vals[0])
						v.StartValue = 0
					}
					v.EndValue, err = strconv.ParseInt(vals[1], 10, 64)
					if err != nil {
						fmt.Println("convert end value to int64 fail ", vals[1])
						v.EndValue = 1 << 62
					}
					continue
				} else if strings.Contains(check[0], ",") {
					//fmt.Println("list columns is ", key)
					//check whether it is an list
					//TODO :Handle errors
					v.TypeGen = 1 //random val from list
					vals := strings.Split(check[0], ",")
					v.DefaultVal = vals
					continue
				}
			}
		}
	}
}

//@consolidateConfigAndMeta
//Consolidate table structure and configuration information
func consolidateConfigAndMeta() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go setColGenerateDataFlag(&wg)
	wg.Add(1)
	go setColProperty(&wg)
	wg.Wait()
	return nil
}
