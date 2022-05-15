package generate

import (
	"github.com/generate_data/meta"
	"github.com/generate_data/output"
	"github.com/generate_data/util"
)

func GenerateDataNormal(tf *output.TableFiles, v *meta.Table, filenum uint64, id uint64, count uint64) error {
	var num uint64
	var record []byte
	var err error
	if count < (id+1)*filenum {
		num = count
	} else {
		num = (id + 1) * filenum
	}
	//var increm_info []util.Incrementinfo
	incremInfo := make([]util.Incrementinfo, len(v.Columns))
	for i := 0; i < len(v.Columns); i++ {
		incremInfo[i].StartValue = v.Columns[i].Property.StartValue + int64(id*filenum)
	}
	for i := id * filenum; i < num; i++ {
		record, err = v.GenerateRecordData(id, incremInfo)
		if err != nil {
			return err
		}
		err = tf.WriteData(v.DBName, v.TableName, record, id)
		if err != nil {
			return err
		}
	}
	return nil
}
