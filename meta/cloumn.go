/**
 * @Author: guobob
 * @Description:
 * @File:  cloumn.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:23
 */

package meta

const (
	INT int = iota
	DECIMAL
	TIMESTAMP
	DATETIME
	STRING
	EUME
)



const (
	VOLATILE int = iota
	UNVOLATILE
)

type Property struct {
	Type       int
	Length     int //total len
	SuffixLen  int //point len
	DefaultVal []string
	StartKey   string
	EndKey     string
	CharFormat []byte //1~9 x
}

type Column struct {
	*Property
	ColumnName string
	ColumnIdx  int
}



func (p *Property) GenerateData() (interface{}, error) {
	switch p.Type{
	case INT:
		
	default:
		return 0,nil
	}
	return "", nil
}
