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
	"github.com/pingcap/errors"
	"os"
	"time"
)

type TableFiles struct {
	files       map[string]*WriteFile
	filePath    string
	filePrefix  string
	maxFileSize uint64
	sync        bool
}

type WriteFile struct {
	fileName    string
	filePath    string
	filePrefix  string
	tableName   string
	dbName      string
	fileNo      uint64
	maxFileSize uint64
	sync        bool
	pos         uint64
	fp          *os.File
	buff        []byte
}

func newWriteFile(tf *TableFiles, tableName, dbName string) *WriteFile {
	return &WriteFile{
		sync:        tf.sync,
		maxFileSize: tf.maxFileSize,
		filePrefix:  tf.filePrefix,
		filePath:    tf.filePath,
		tableName:   tableName,
		dbName:      dbName,
		fileNo:      0,
	}
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
	wf.fileName = fmt.Sprintf("%v-%v.%v-%6d", wf.filePrefix, wf.tableName, wf.tableName, wf.fileNo)
}

func (wf *WriteFile) close() {
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
	return err
}

func (wf *WriteFile) checkIfNeedChangeFile() bool {
	if wf.pos >= wf.maxFileSize {
		return true
	}
	return false
}

func NewTableFiles(sync bool, maxFileSize uint64, path, filePrefix string) *TableFiles {
	if len(filePrefix) == 0 {
		ts := time.Now()
		filePrefix = fmt.Sprintf("%v%02d%02d", ts.Year(), ts.Month(), ts.Day())
	}
	return &TableFiles{
		sync:        sync,
		maxFileSize: maxFileSize,
		filePath:    path,
		filePrefix:  filePrefix,
		files:       make(map[string]*WriteFile),
	}
}

func (tf *TableFiles) Close() {
	for _, v := range tf.files {
		v.close()
	}
}

func (tf *TableFiles) WriteDataToFile(dbName, tableName string, buff []byte) error {
	var err error
	key := fmt.Sprintf("%v.%v", dbName, tableName)
	v, ok := tf.files[key]
	if !ok {
		wf := newWriteFile(tf, tableName, dbName)
		tf.files[key] = wf
		v = wf
		v.getFileNo()
		v.generateFileName()
		err = v.openFile()
		if err != nil {
			return err
		}
		v.pos = 0
	}
	v.buff = buff
	err = v.write()
	if err != nil {
		return err
	}
	if v.checkIfNeedChangeFile() {
		v.close()
		v.getFileNo()
		v.generateFileName()
		err := v.openFile()
		if err != nil {
			return err
		}
		v.pos = 0
	}
	return err
}
