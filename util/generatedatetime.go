/**
 * @Author: guobob
 * @Description:
 * @File:  generateint.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 23:05
 */

package util

import (
	"math/rand"
	"time"

	"github.com/fufuok/random"
)

var intmax int = 1 << 32
var timeUnix int64 = time.Now().Unix()

func Randdatetime_1(a *Property) (string, error) {
	//timeUnix := time.Now().Unix()
	randtime := rand.Int63n(timeUnix)
	timeStr := time.Unix(randtime, 0).Format("2006-01-02 15:04:05.000")
	return timeStr, nil
}
func Randdatetime(a *Property) (string, error) {
	highint := random.FastIntn(intmax)
	lowint := random.FastIntn(intmax)
	randtime := int64((highint<<32 + lowint)) % timeUnix
	timeStr := time.Unix(randtime, 0).Format("2006-01-02 15:04:05.000")
	return timeStr, nil
}
