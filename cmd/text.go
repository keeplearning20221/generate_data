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
	"github.com/generate_data/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewTextCommand() *cobra.Command {

	var (
		dsn        string
		tables     string
		conFile    string
		fieldTerm  string
		lineTerm   string
		outputPath string
		count      int64
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

			if util.GConfig != nil {
				maxFileSize, err = util.GConfig.GetMaxFileSize()
				if err != nil {
					return err
				}
				filePrefix = util.GConfig.GetfilePrefix()
				outputPath = util.GConfig.GetOutputfile()
				tables = util.GConfig.GetTables()
				count, err = util.GConfig.GetRowCount()
				if err != nil {
					return err
				}
			}
			err = meta.GetTableInfo(tables, dsn, cfg, fieldTerm, lineTerm, log)
			if err != nil {
				log.Error("get meta data fail" + err.Error())
				return err
			}

			var i int64 = 0
			for _, v := range meta.Gmeta {
				tf := output.NewTableFiles(false, maxFileSize, outputPath, filePrefix)
				fmt.Println(count)
				for i = 0; i < count; i++ {
					record, err := v.GenerateRecordData()
					if err != nil {
						return err
					}
					err = tf.WriteData(v.DBName, v.TableName, []byte(record))
					if err != nil {
						return err
					}
				}
				err = tf.Sync()
				if err != nil {
					log.Error("write data fail " + err.Error())
					return err
				}
				tf.Close()
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dsn, "dsn", "d", "", "meta data  server dsn")
	cmd.Flags().StringVarP(&tables, "table", "t", "", "table list , test.t")
	cmd.Flags().StringVarP(&fieldTerm, "fieldterm", "f", "\t", "data filed terminated by ")
	cmd.Flags().StringVarP(&lineTerm, "lineterm", "l", "\n", "data record terminated by ")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "./", "out file path")
	cmd.Flags().StringVarP(&filePrefix, "filePrefix", "p", " ", "file name prefix")
	cmd.Flags().Int64VarP(&count, "filesize", "n", 100, "genereate data row count")
	cmd.Flags().StringVarP(&conFile, "config", "c", "", "config output name ")
	return cmd
}
