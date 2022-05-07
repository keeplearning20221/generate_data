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

	"github.com/fufuok/random"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandString_1(a *Property) (string, error) {

	if (len(a.DefaultVal)) != 0 {
		strnum := len(a.DefaultVal)
		return a.DefaultVal[rand.Intn(strnum)], nil
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
		randlen := rand.Intn(a.CharLen-3) + 3
		for i := 0; i < randlen; i++ {

			bytes[i] = byte(rand.Intn(26) + 97)
		}
		return string(bytes[0:randlen]), nil
	}

	if a.CharFormat == nil {
		for ; i < a.CharLen-end; i++ {
			b := rand.Intn(26) + 97
			bytes[i] = byte(b)
		}
	} else {
		for ; i < a.CharLen-end; i++ {
			num := len(a.CharFormat)
			b := a.CharFormat[rand.Intn(num)]
			bytes[i] = byte(b)
		}
	}
	return string(bytes), nil
}

func RandString(a *Property) (string, error) {

	if (len(a.DefaultVal)) != 0 {
		strnum := len(a.DefaultVal)
		return a.DefaultVal[random.FastIntn(strnum)], nil
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
		randlen := random.FastIntn(a.CharLen-3) + 3
		for i := 0; i < randlen; i++ {

			bytes[i] = byte(random.FastIntn(26) + 97)
		}
		return string(bytes[0:randlen]), nil
	}

	if a.CharFormat == nil {
		for ; i < a.CharLen-end; i++ {
			b := random.FastIntn(26) + 97
			bytes[i] = byte(b)
		}
	} else {
		for ; i < a.CharLen-end; i++ {
			num := len(a.CharFormat)
			b := a.CharFormat[random.FastIntn(num)]
			bytes[i] = byte(b)
		}
	}
	return string(bytes), nil
}

func IncrementString(a *Property, increm_info *Incrementinfo) (string, error) {

	if increm_info.NowValue < increm_info.StartValue {
		increm_info.NowValue = increm_info.StartValue
		str := strconv.FormatInt(increm_info.NowValue, 10)
		if len(str) > a.CharLen {
			err := fmt.Errorf("nowvalue long then CharLen")
			return "", err
		}
		return str, nil
	}

	if a.EndValue != 0 && increm_info.NowValue > a.EndValue {
		err := fmt.Errorf("string nowvalue is out of range")
		return "", err
	}
	increm_info.NowValue++
	str := strconv.FormatInt(increm_info.NowValue, 10)
	if len(str) > a.CharLen {
		err := fmt.Errorf("nowvalue long then CharLen")
		return "", err
	}
	return str, nil
}
