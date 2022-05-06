/**
 * @Author: guobob
 * @Description:
 * @File:  generateint.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 23:05
 */

package util

import (
	"fmt"
	"math/rand"
)

func RandCNString(a *Property) (string, error) {
	if (len(a.DefaultVal)) != 0 {
		strnum := len(a.DefaultVal)
		return a.DefaultVal[rand.Intn(strnum)], nil
	}
	bytes := make([]rune, a.Length)
	var start int
	var end int
	start = 0
	end = 0
	if len(a.StartKey) > 0 {
		start = len([]rune(a.StartKey))
	}
	if len(a.EndKey) > 0 {
		end = len([]rune(a.EndKey))
		if end > a.Length {
			err := fmt.Errorf("startkey and endkey long then length")
			return "", err
		}
	}
	if start+end > a.Length {
		err := fmt.Errorf("startkey and endkey long then length")
		return "", err
	}
	if a.CharFormat == nil {
		for i := start; i < a.Length-end; i++ {
			b := rand.Intn(40869-19968) + 19968
			bytes[i] = rune(b)
		}
	} else {
		for i := start; i < a.Length-end; i++ {
			num := len(a.CharFormat)
			b := a.CharFormat[rand.Intn(num)]
			bytes[i] = rune(b)
		}
	}
	return a.StartKey + string(bytes[start:a.Length-end]) + a.EndKey, nil
}
