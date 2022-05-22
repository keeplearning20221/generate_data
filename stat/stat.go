package stat

import (
	"errors"
	"fmt"
	"sync"
)

const (
	NOT_START int = iota
	RUNING
	END
)

type tableStatus struct {
	state    int
	rowCount uint64
}

var (
	gStat map[string]*tableStatus
	mu    sync.Mutex
)

func init() {
	gStat = make(map[string]*tableStatus)
}

// add table to
func AddTable(tableName string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := gStat[tableName]; !ok {
		gStat[tableName] = &tableStatus{NOT_START, 0}
		return nil
	} else {
		return errors.New(fmt.Sprintf("%v is exist", tableName))
	}
}

func ChangeTableStat(tableName string, state int) error {
	mu.Lock()
	defer mu.Unlock()
	v, ok := gStat[tableName]
	if !ok {
		return errors.New(fmt.Sprintf("%v is not exist", tableName))
	}
	v.state = state
	v.rowCount = 0

	return nil
}

func AddTableRows(tableName string, addRows uint64) error {
	PrintStatis()
	mu.Lock()
	defer mu.Unlock()
	v, ok := gStat[tableName]
	if !ok {
		return errors.New(fmt.Sprintf("%v is not exist", tableName))
	}
	if v.state != RUNING {
		return errors.New(fmt.Sprintf("%v stat is not  running,%v", tableName, v.state))
	} else {
		v.rowCount += addRows
		return nil
	}

}

func PrintStatis() {
	mu.Lock()
	defer mu.Unlock()
	var dones int = 0
	var runnings int = 0
	for _, v := range gStat {
		if v.state == END {
			dones++
		} else if v.state == RUNING {
			runnings++
		}
	}
	fmt.Println(fmt.Sprintf("generate data progress(sum,runing,done):%v/%v/%v", len(gStat), runnings, dones))
	for k, v := range gStat {
		if v.state == runnings {
			fmt.Println(fmt.Sprintf("%v:%v", k, v.rowCount))
		}
	}
}
