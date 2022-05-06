/**
 * @Author: guobob
 * @Description:
 * @File:  column.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:23
 */

package meta

import (
	"fmt"
	"strconv"

	"github.com/generate_data/sql"
	"github.com/generate_data/util"
	"github.com/pingcap/errors"
)

type Column struct {
	*util.Property
	ColumnName string
	ColumnIdx  int
	Ignore     bool
}

/*
 GetColumnFromMetaData
 table_name
 column_name
 ORDINAL_POSITION
 DATA_TYPE
 CHARACTER_MAXIMUM_LENGTH
 CHARACTER_OCTET_LENGTH
 NUMERIC_PRECISION
 NUMERIC_SCALE
 DATETIME_PRECISION
*/
func GetColumnFromMetaData(s *sql.SQLHandle, t *Table) error {
	var err error
	fmt.Println(s.SqlRes)
	for _, v := range s.SqlRes {
		col := new(Column)
		col.Ignore = false
		col.Property = new(util.Property)
		col.ColumnName = v[1]

		if len(v[2]) == 0 {
			col.ColumnIdx = 0
		} else {
			col.ColumnIdx, err = strconv.Atoi(v[2])
			if err != nil {
				fmt.Println("get column index fail," + err.Error())
				return err
			}
		}

		if len(v[3]) == 0 {
			return errors.New(fmt.Sprintf("unsupport type "))
		}
		col.Type = util.ChangeColType(v[3])
		if col.Type == -1 {
			return errors.New(fmt.Sprintf("unsupport type %v", v[3]))
		}
		if len(v[4]) == 0 {
			col.CharLen = 0
		} else {
			col.CharLen, err = strconv.Atoi(v[4])
			if err != nil {
				fmt.Println("get column CharLen fail," + err.Error())
				return err
			}
		}
		if len(v[5]) == 0 {
			col.BitLen = 0
		} else {
			col.BitLen, err = strconv.Atoi(v[5])
			if err != nil {
				fmt.Println("get column BitLen fail,", col.ColumnIdx, err.Error())
				return err
			}
		}
		if len(v[6]) == 0 {
			col.Length = 0
		} else {
			col.Length, err = strconv.Atoi(v[6])
			if err != nil {
				fmt.Println("get column Length fail,", v[6], col.ColumnIdx, err.Error())
				return err
			}
		}
		if len(v[7]) == 0 {
			col.SuffixLen = 0
		} else {
			col.SuffixLen, err = strconv.Atoi(v[7])
			if err != nil {
				fmt.Println("get column SuffixLen fail," + err.Error())
				return err
			}
		}

		if len(v[8]) == 0 {
			col.TypeGen = 1
		} else {
			col.TypeGen, err = strconv.Atoi(v[8])
			if err != nil {
				fmt.Println("get column covertype fail," + err.Error())
				return err
			}
		}
		t.Columns = append(t.Columns, *col)
		fmt.Println(t.Columns)
	}
	return nil
}
