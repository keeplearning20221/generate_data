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
	DefaultVal []string
	StartKey   string
	EndKey     string
	CharFormat []byte //1~9 x
}

func (p *Property) GenerateColumnData() (string, error) {

	switch p.Type {
	case INT:
		num, res := Randint(p)
		return fmt.Sprintf("%v", num), res
	case STRING:
		str, res := RandString(p)
		return str, res
	case STRINGCN:
		str, res := RandCNString(p)
		return str, res
	case DECIMAL:
		decimal, res := Randdecimal(p)
		return fmt.Sprintf("%v", decimal), res
	default:
		return "", errors.New("unsupport type")
	}

}
