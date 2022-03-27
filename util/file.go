/**
 * @Author: guobob
 * @Description:
 * @File:  output.go
 * @Version: 1.0.0
 * @Date: 2022/3/26 19:24
 */

package util

import (
	"os"
	"syscall"
)

func CheckFileExistAndPrivileges(fileName string) bool {
	err := syscall.Access(fileName, syscall.F_OK)
	return !os.IsNotExist(err)
}
