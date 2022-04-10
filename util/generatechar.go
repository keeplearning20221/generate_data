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
	"time"
)

func RandString(a *Property) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if (len(a.DefaultVal)) != 0 {
		strnum := len(a.DefaultVal)
		return a.DefaultVal[r.Intn(strnum)], nil
	}
	bytes := make([]byte, a.CharLen)
	var i int
	var end int
	i = 0
	end = 0
	if len(a.StartKey) > 0 {
		i = len(a.StartKey)
		for j := 0; j < i; j++ {
			bytes[j] = byte(a.StartKey[j])
		}
	}
	if len(a.EndKey) > 0 {
		end = len(a.EndKey)
		if end > a.CharLen {
			err := fmt.Errorf("startkey and endkey long then CharLen")
			return "", err
		}

		for j := a.CharLen - end; j < a.CharLen; j++ {
			bytes[j] = byte(a.EndKey[end+j-a.CharLen])
		}
	}
	if i+end > a.CharLen {
		err := fmt.Errorf("startkey and endkey long then CharLen")
		return "", err
	}
	if a.CharFormat == nil {
		for ; i < a.CharLen-end; i++ {
			b := r.Intn(48) + 42
			bytes[i] = byte(b)
		}
	} else {
		for ; i < a.CharLen-end; i++ {
			num := len(a.CharFormat)
			b := a.CharFormat[r.Intn(num)]
			bytes[i] = byte(b)
		}
	}
	return string(bytes), nil
}
