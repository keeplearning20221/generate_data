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
	"math"
	"math/rand"
	"strconv"
	"time"
)

func Randint(a *Property) (int64, error) {
	if a.Length > 19 {
		a.Length = 18
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if (len(a.DefaultVal)) != 0 {
		strnum := len(a.DefaultVal)
		num, err := strconv.ParseInt(a.DefaultVal[r.Intn(strnum)], 10, 64)
		return num, err
	}
	bytes := make([]byte, a.Length)
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
		if end > a.Length {
			err := fmt.Errorf("startkey and endkey long then length")
			return 0, err
		}

		for j := a.Length - end; j < a.Length; j++ {
			bytes[j] = byte(a.EndKey[end+j-a.Length])
		}
	}
	if i+end > a.Length {
		err := fmt.Errorf("startkey and endkey long then length")
		return 0, err
	}
	if a.CharFormat == nil {
		for ; i < a.Length-end; i++ {
			b := r.Intn(10) + 48
			bytes[i] = byte(b)
		}
	} else {
		for ; i < a.Length-end; i++ {
			num := len(a.CharFormat)
			b := a.CharFormat[r.Intn(num)]
			bytes[i] = byte(b)
		}
	}
	num, err := strconv.ParseInt(string(bytes), 10, 64)
	return num, err

}

func Incrementint(a *Property) (int64, error) {
	if a.NowValue < a.StartValue {
		a.NowValue = a.StartValue
		if a.NowValue >= int64(math.Pow10(a.Length)) {
			err := fmt.Errorf("nowvalue long then CharLen")
			return 0, err
		}
		return a.NowValue, nil
	}

	if a.EndValue != 0 && a.NowValue > a.EndValue {
		err := fmt.Errorf("int nowvalue is out of range")
		return 0, err
	}
	a.NowValue++
	if a.NowValue >= int64(math.Pow10(a.Length)) {
		err := fmt.Errorf("nowvalue long then CharLen")
		return 0, err
	}
	return a.NowValue, nil
}
