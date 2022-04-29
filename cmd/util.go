package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/generate_data/meta"
	"github.com/generate_data/util"
)

//Flag whether the column needs to generate data
func SetColGenerateDataFlag(wg *sync.WaitGroup) {
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
	wg.Done()
}

func SetColProperty(wg *sync.WaitGroup) {
	for tableName, t := range meta.Gmeta {
		for _, v := range t.Columns {
			key := fmt.Sprintf("%v.%v", tableName, v.ColumnName)
			if check, ok := util.GConfig.Checks[key]; !ok {
				continue
			} else {
				if strings.Contains(check[0], "~") {
					//TODO :set column Property
					return
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
