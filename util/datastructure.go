/**
 * @Author: guobob
 * @Description:
 * @File:  column.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:23
 */

package util

import (
	"fmt"

	"github.com/pingcap/errors"
)

type Property struct {
	Type       int
	Length     int //total len
	SuffixLen  int //point len
	CharLen    int
	BitLen     int
	StartKey   string
	EndKey     string
	StartValue int64
	EndValue   int64
	NowValue   int64
	CharFormat []byte //1~9 x
	TypeGen    int    //1:random 2: range
	// If defaultVal is not empty and TypeGen is random,
	// the value of DefaultVal is randomly fetched and populated
	DefaultVal []string
}

func (p *Property) GenerateColumnData() (string, error) {

	switch p.Type {
	case INT:
		switch p.TypeGen {
		case 1:
			num, res := Randint(p)
			return fmt.Sprintf("%v", num), res
		case 2:
			num, res := Incrementint(p)
			return fmt.Sprintf("%v", num), res
		default:
			return "", errors.New("unsupport type gen")
		}

	case STRING:
		switch p.TypeGen {
		case 1:
			str, res := RandString(p)
			return str, res
		case 2:
			str, res := IncrementString(p)
			return str, res
		default:
			return "", errors.New("unsupport type gen")
		}
	case STRINGCN:
		str, res := RandCNString(p)
		return str, res
	case DECIMAL:
		decimal, res := Randdecimal(p)
		return fmt.Sprintf("%v", decimal), res
	case DATETIME:
		datetimestr, res := Randdatetime(p)
		return datetimestr, res
	default:
		return "", errors.New("unsupport type")
	}

}
