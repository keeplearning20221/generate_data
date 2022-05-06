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
	"strconv"
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
	//If you need to generate a string of more than 3 bytes, randomly generate 3~maximum length characters
	if i == 0 && end == 0 && a.CharLen > 3 {
		randlen := r.Intn(a.CharLen-3) + 3
		for i := 0; i < randlen; i++ {

			bytes[i] = byte(r.Intn(26) + 97)
		}
		return string(bytes[0:randlen]), nil
	}

	if a.CharFormat == nil {
		for ; i < a.CharLen-end; i++ {
			b := r.Intn(26) + 97
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

func IncrementString(a *Property) (string, error) {

	if a.NowValue < a.StartValue {
		a.NowValue = a.StartValue
		str := strconv.FormatInt(a.NowValue, 10)
		if len(str) > a.CharLen {
			err := fmt.Errorf("nowvalue long then CharLen")
			return "", err
		}
		return str, nil
	}

	if a.EndValue != 0 && a.NowValue > a.EndValue {
		err := fmt.Errorf("string nowvalue is out of range")
		return "", err
	}
	a.NowValue++
	str := strconv.FormatInt(a.NowValue, 10)
	if len(str) > a.CharLen {
		err := fmt.Errorf("nowvalue long then CharLen")
		return "", err
	}
	return str, nil
