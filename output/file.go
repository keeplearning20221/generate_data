/**
 * @Author: guobob
 * @Description:
 * @File:  output.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 23:00
 */

package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pingcap/errors"
)

type WriteFile struct {
	fileName    string
	filePath    string
	filePrefix  string
	tableName   string
	dbName      string
	fileNo      uint64
	maxFileSize uint64
	maxFileNum  uint64
	currentNum  uint64
	sync        bool
	pos         uint64
	fp          *os.File
	buff        []byte
	writechan   chan []byte
	quitchan    chan int
}

func newWriteFile(tf *TableFiles, tableName, dbName string, fileNo uint64) *WriteFile {
	wf := &WriteFile{
		sync:        tf.sync,
		maxFileSize: tf.maxFileSize,
		maxFileNum:  tf.maxFileNum,
		filePrefix:  tf.filePrefix,
		filePath:    tf.filePath,
		tableName:   tableName,
		dbName:      dbName,
		currentNum:  0,
		fileNo:      fileNo,
		writechan:   make(chan []byte, 100),
		quitchan:    make(chan int, 1),
	}
	go writeFileSync(wf)
	return wf
}

func writeFileSync(wf *WriteFile) {
	for {
		select {
		case wf.buff = <-wf.writechan:
			err := wf.write()
			if err != nil {
				//TODO: Handle errors
				fmt.Println("write data fail", err)
				os.Exit(1)
			}
		case <-wf.quitchan:
			return
		}
	}
}

func (wf *WriteFile) WriteAsync(buff []byte) {
	wf.writechan <- buff
}

func (wf *WriteFile) openFile() error {
	var err error
	if len(wf.filePath) == 0 || len(wf.fileName) == 0 {
		err := errors.New("path or filename len is 0")
		return err
	}

	fn := wf.filePath + "/" + wf.fileName
	wf.fp, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	return err
}

func (wf *WriteFile) getFileNo() {
	wf.fileNo = wf.fileNo + 1
}
func (wf *WriteFile) generateFileName() {
	if len(strings.TrimSpace(wf.filePrefix)) == 0 {
		wf.fileName = fmt.Sprintf("%v.%v.%v.csv", wf.dbName, wf.tableName, wf.fileNo)
	} else {
		wf.fileName = fmt.Sprintf("%v.%v.%v.%v.csv", wf.filePrefix, wf.dbName, wf.tableName, wf.fileNo)
	}

}

func (wf *WriteFile) close() {
	for {
		if len(wf.writechan) > 0 {
			continue
		} else {
			break
		}
	}
	time.Sleep(100 * time.Millisecond)
	wf.quitchan <- 1
	err := wf.fp.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (wf *WriteFile) write() error {
	length, err := wf.fp.Write(wf.buff)
	if length < len(wf.buff) || err != nil {
		return err
	}
	if wf.sync {
		err = wf.fp.Sync()
		if err != nil {
			return err
		}
	}
	wf.pos = wf.pos + uint64(length)
	wf.currentNum++
	return err
}

func (wf *WriteFile) checkIfNeedChangeFile() bool {
	return (wf.pos >= wf.maxFileSize && wf.currentNum >= wf.maxFileNum) && false
}
