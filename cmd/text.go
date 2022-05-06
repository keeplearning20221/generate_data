/**
 * @Author: guobob
 * @Description:
 * @File:  text.go
 * @Version: 1.0.0
 * @Date: 2022/3/25 09:39
 */

package cmd

import (
	"fmt"

	"github.com/generate_data/meta"
	"github.com/generate_data/output"
	"github.com/generate_data/sigLimit"
	"github.com/generate_data/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func write_file(tf *output.TableFiles, v *meta.Table, filenum uint64, id uint64, count uint64) error {
	var num uint64
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
		record, err := v.GenerateRecordData(id, incremInfo)
		if err != nil {
			return err
		}
		err = tf.WriteData(v.DBName, v.TableName, []byte(record), id)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTextCommand() *cobra.Command {

	var (
		dsn        string
		tables     string
		conFile    string
		fieldTerm  string
		lineTerm   string
		outputPath string
		count      uint64
		filePrefix string
	)
	cmd := &cobra.Command{
		Use:   "text",
		Short: "generate data in csv ",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			log := zap.L().Named("csv-data")
			cfg, err := util.ParseDSN(dsn)
			if err != nil {
				return err
			}

			if conFile != "" {
				err = util.NewConfig(conFile)
				if err != nil {
					fmt.Println(err)
					return err
				}

				fmt.Println(util.GConfig)
			} else {
				fmt.Println("no config file")
			}
			var maxFileSize uint64 = 100 * 1024 * 1024
			var maxFileNum uint64 = 10000

			if util.GConfig != nil {
				maxFileSize, err = util.GConfig.GetMaxFileSize()
				if err != nil {
					return err
				}
				//convert MB to  Byte
				maxFileSize = maxFileSize * 1024 * 1024
				maxFileNum, err = util.GConfig.GetMaxFileNum()
				if err != nil {
					return err
				}
				filePrefix = util.GConfig.GetfilePrefix()
				outputPath = util.GConfig.GetOutputfile()
				tables = util.GConfig.GetTables()
				count, err = util.GConfig.GetRowcount()
				if err != nil {
					return err
				}

			}
			err = meta.GetTableInfo(tables, dsn, cfg, fieldTerm, lineTerm, log)
			if err != nil {
				log.Error("get meta data fail" + err.Error())
				return err
			}

			err = consolidateConfigAndMeta()
			if err != nil {
				return err
			}
			fmt.Println("-----------------------------")
			for _, v := range meta.Gmeta {
				for _, vv := range v.Columns {
					fmt.Println(vv, vv.DefaultVal, vv.TypeGen, vv.StartValue, vv.EndValue)
				}
			}
			fmt.Println("-----------------------------")
			for _, v := range meta.Gmeta {

				fmt.Println(count)
				s := sigLimit.NewSigLimit(10)
				var i uint64 = 0
				for i = 0; i <= count/maxFileNum; i++ {
					s.Add()
					go func(i uint64) {
						tf := output.NewTableFiles(false, maxFileSize, maxFileNum, outputPath, filePrefix)
						defer s.Done()
						defer tf.Close()
						err := write_file(tf, &v, maxFileNum, i, count)
						if err != nil {
							log.Error("write file fail" + err.Error())

						}

					}(i)
				}
				s.Wait()
				// for i = 0; i < count; i++ {
				// 	record, err := v.GenerateRecordData()
				// 	if err != nil {
				// 		return err
				// 	}
				// 	err = tf.WriteData(v.DBName, v.TableName, []byte(record), 0)
				// 	if err != nil {
				// 		return err
				// 	}
				// }
				// err = tf.Sync()
				// if err != nil {
				// 	log.Error("write data fail " + err.Error())
				// 	return err
				// }
				//go write_file(tf, &v, maxFileNum, i, count)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dsn, "dsn", "d", "", "meta data  server dsn")
	cmd.Flags().StringVarP(&tables, "table", "t", "", "table list , test.t")
	cmd.Flags().StringVarP(&fieldTerm, "fieldterm", "f", ",", "data filed terminated by ")
	cmd.Flags().StringVarP(&lineTerm, "lineterm", "l", "\n", "data record terminated by ")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "./", "out file path")
	cmd.Flags().StringVarP(&filePrefix, "filePrefix", "p", " ", "file name prefix")
	cmd.Flags().Uint64VarP(&count, "filesize", "n", 100, "genereate data row count")
	cmd.Flags().StringVarP(&conFile, "config", "c", "", "config output name ")
	return cmd
}
