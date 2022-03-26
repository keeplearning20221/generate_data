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
	"github.com/generate_data/sql"
	"github.com/generate_data/util"
	"github.com/pingcap/errors"
)

type Column struct {
	*util.Property
	ColumnName string
	ColumnIdx  int
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
		col.Property = new(util.Property)
		err = util.ConvertAssign(&col.ColumnName, v[1])
		if err != nil {
			fmt.Println("get column name fail," + err.Error())
			return err
		}
		err = util.ConvertAssign(&col.ColumnIdx, v[2])
		if err != nil {
			fmt.Println("get column index fail," + err.Error())
			return err
		}
		var dataType string
		err = util.ConvertAssign(&dataType, v[3])
		if err != nil {
			fmt.Println("get column dataType fail," + err.Error())
			return err
		}
		fmt.Println(dataType)
		col.Type = util.ChangeColType(dataType)
		if col.Type == -1 {
			return errors.New(fmt.Sprintf("unsupport type %v", dataType))
		}
		if v[4] == nil {
			col.CharLen = 0
		} else {
			err = util.ConvertAssign(&col.CharLen, v[4])
			if err != nil {
				fmt.Println("get column CharLen fail," + err.Error())
				return err
			}
		}
		if v[5] == nil {
			col.BitLen = 0
		} else {
			err = util.ConvertAssign(&col.BitLen, v[5])
			if err != nil {
				fmt.Println("get column BitLen fail," + err.Error())
				return err
			}
		}
		if v[6] == nil {
			col.Length = 0
		} else {
			err = util.ConvertAssign(&col.Length, v[6])
			if err != nil {
				fmt.Println("get column Length fail," + err.Error())
				return err
			}
		}
		if v[7] == nil {
			col.SuffixLen = 0
		} else {
			err = util.ConvertAssign(&col.SuffixLen, v[7])
			if err != nil {
				fmt.Println("get column SuffixLen fail," + err.Error())
				return err
			}
		}
		t.Columns = append(t.Columns, *col)
	}
	return nil
}
