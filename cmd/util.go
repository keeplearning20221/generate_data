package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/generate_data/meta"
	"github.com/generate_data/util"
)

//@SetColGenerateDataFlag
//Flag whether the column needs to generate data
func SetColGenerateDataFlag(wg *sync.WaitGroup) {
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
func SetColProperty(wg *sync.WaitGroup) {
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
	go SetColGenerateDataFlag(&wg)
	wg.Add(1)
	go SetColProperty(&wg)
	wg.Wait()
	return nil
}
