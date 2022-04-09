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
		dsn         string
		tables      string
		conFile     string
		fieldTerm   string
		lineTerm    string
		outputPath  string
		count       int64
		maxFileSize uint64
		filePrefix  string
	)
	cmd := &cobra.Command{
		Use:   "text",
		Short: "generate data in csv ",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			log := zap.L().Named("csv-data")
			/*
				if !util.CheckFileExistAndPrivileges(conFile) {
					return errors.New("config output is not exist or privileges incorrect ")
				}*/
			cfg, err := util.ParseDSN(dsn)
			if err != nil {
				return err
			}
			fmt.Println(tables)
			err = meta.GetTableInfo(tables, dsn, cfg, log)
			if err != nil {
				log.Error("get meta data fail" + err.Error())
				return err
			}

			err = meta.Generate_tables_data(&meta.Gmeta)
			fmt.Println(meta.Gmeta)

			if maxFileSize == 0 {
				maxFileSize = 100 * 1024 * 1024
			}
			var i int64 = 0
			for _, v := range meta.Gmeta {
				//v.GeneratePrepareSQL()
				//fmt.Println(v.PrepareSQL)
				tf := output.NewTableFiles(false, maxFileSize, outputPath, filePrefix)
				for i = 0; i < count; i++ {
					record, err := meta.Generate_table_data(v)
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
	cmd.Flags().StringVarP(&fieldTerm, "fieldterm", "f", "", "data filed terminated by ")
	cmd.Flags().StringVarP(&lineTerm, "lineterm", "l", "", "data record terminated by ")
	cmd.Flags().StringVarP(&lineTerm, "lineterm", "l", "", "data record terminated by ")
	cmd.Flags().StringVarP(&outputPath, "output", "0", "./", "out file path")
	cmd.Flags().StringVarP(&filePrefix, "filePrefix", "fp", "", "file name prefix")
	count = *cmd.Flags().Int64P("count", "c", 0, "genereate data row count")
	cmd.Flags().StringVarP(&conFile, "config", "c", "", "config output name ")
	return cmd
}
