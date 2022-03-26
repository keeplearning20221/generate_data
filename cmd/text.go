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
	"github.com/generate_data/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewTextCommand() *cobra.Command {

	var (
		dsn    string
		tables string
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
			fmt.Println(tables)
			err = meta.GetTableInfo(tables, dsn, cfg, log)
			if err != nil {
				log.Error("get meta data fail" + err.Error())
				return err
			}
			fmt.Println(meta.Gmeta)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dsn, "dsn", "d", "", "meta data  server dsn")
	cmd.Flags().StringVarP(&tables, "table", "t", "", "table list , test.t")
	return cmd
}
