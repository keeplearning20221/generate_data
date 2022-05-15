/**
 * @Author: guobob
 * @Description:
 * @File:  text.go
 * @Version: 1.0.0
 * @Date: 2022/3/25 09:39
 */

package cmd

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/generate_data/generate"
	"github.com/generate_data/meta"
	"github.com/generate_data/output"
	"github.com/generate_data/sigLimit"
	"github.com/generate_data/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewTextCommand() *cobra.Command {

	var (
		conFile string
		bc      baseConfig
	)
	cmd := &cobra.Command{
		Use:   "text",
		Short: "generate data in csv ",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("CPU线程数:", runtime.NumCPU())

			var err error
			log := zap.L().Named("csv-data")

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

			if util.GConfig != nil {
				err = bc.getVals()
				if err != nil {
					return err
				}
			} else {
				return errors.New("Configuration files need to be specified ")
			}
			err = meta.GetTableInfo(bc.tables, bc.dsn, bc.cfg, bc.fieldTerm, bc.lineTerm, log)
			if err != nil {
				log.Error("get meta data fail" + err.Error())
				return err
			}

			/*
				fmt.Println("-------------1111----------------")
				for _, v := range meta.Gmeta {
					for _, vv := range v.Columns {
						fmt.Println(vv)
					}
				}
				fmt.Println("-------------1111----------------")
			*/
			err = consolidateConfigAndMeta()
			if err != nil {
				return err
			}
			/*
				fmt.Println("-----------------------------")
				for _, v := range meta.Gmeta {
					for _, vv := range v.Columns {
						fmt.Println(vv, vv.DefaultVal, vv.TypeGen, vv.StartValue, vv.EndValue)
					}
				}
				fmt.Println("-----------------------------")
			*/

			fmt.Println("-----------begin--", time.Now().String(), "------------")
			for _, v := range meta.Gmeta {
				fmt.Println(bc.count)
				s := sigLimit.NewSigLimit(bc.threadPoolSize, zap.L().Named(fmt.Sprintf("%v.%v", v.DBName, v.TableName)))
				var i uint64 = 0
				for i = 0; i <= bc.count/bc.maxFileNum; i++ {
					s.Add()
					go func(i uint64) error {
						tf := output.NewTableFiles(false, bc.maxFileSize, bc.maxFileNum, bc.outputPath, bc.filePrefix)
						defer s.Done()
						defer tf.Close()
						err := generate.GenerateDataNormal(tf, &v, bc.maxFileNum, i, bc.count)
						if err != nil {
							log.Error("write file fail" + err.Error())
							return err
						}
						return nil
					}(i)
				}
				s.Wait()
				fmt.Println("-----------end--", time.Now().String(), "------------")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&conFile, "config", "c", "", "config output name ")
	return cmd
}
