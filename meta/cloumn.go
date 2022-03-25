/**
 * @Author: guobob
 * @Description:
 * @File:  cloumn.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:23
 */

package meta

import (
	"generate_data_module/util"
)

type Column struct {
	*util.Property
	ColumnName string
	ColumnIdx  int
}
