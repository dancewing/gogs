package models

import (
	"bytes"
	"io"
	"io/ioutil"

	"fmt"

	"github.com/kataras/iris/core/errors"
)

func logFind(proc *Proc) (*logData, error) {

	logs := make([]*logData, 0)
	err := x.Where("proc_id = ? ", proc.ID).Limit(1, 0).Find(&logs)
	if err != nil {
		return nil, err
	}

	if len(logs) > 0 {
		return logs[0], nil
	}
	return nil, errors.New(fmt.Sprintf("not found log data with job id %d", proc.ID))
}

func LogFind(proc *Proc) (io.ReadCloser, error) {

	data, err := logFind(proc)

	buf := bytes.NewBuffer(data.Data)

	return ioutil.NopCloser(buf), err
}

func LogSave(proc *Proc, r io.Reader) (err error) {

	data, err := logFind(proc)

	if err != nil {
		data = &logData{ProcID: proc.ID}
		data.Data, _ = ioutil.ReadAll(r)
		_, err = x.InsertOne(data)
		return err
	} else {
		data.Data, _ = ioutil.ReadAll(r)
		_, err = x.ID(data.ID).AllCols().Update(data)
		return err
	}

}

type logData struct {
	ID     int64
	ProcID int64
	Data   []byte `xml:"TEXT"`
}

func (t logData) TableName() string {
	return "cncd_log_data"
}
