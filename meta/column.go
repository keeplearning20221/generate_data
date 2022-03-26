/**
 * @Author: guobob
 * @Description:
 * @File:  column.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:23
 */

package meta

import (
	"github.com/generate_data/sql"
	"github.com/generate_data/util"
)

type Column struct {
	*util.Property
	ColumnName string
	ColumnIdx  int
}


/*
table_name
column_name
ORDINAL_POSITION
DATA_TYPE
NUMERIC_PRECISION
NUMERIC_SCALE
DATETIME_PRECISION
*/
func GetColumnFromMetaData ( s *sql.SQLHandle,t *Table)  error {
	var err error
	for  _,v := range s.SqlRes{
		col := new(Column)
		err = util.ConvertAssign(col.ColumnName,v[1])
		if err != nil {
			return err
		}
		err = util.ConvertAssign(col.ColumnIdx,v[2])
		if err != nil {
			return err
		}
		var dataType string
		err = util.ConvertAssign(dataType,v[1])
		if err != nil {
			return err
		}

	}



	return nil
}
