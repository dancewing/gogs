package models

import (
	"io"
	"io/ioutil"
)

// File represents a pipeline artifact.
type File struct {
	ID      int64  `json:"id"      `
	BuildID int64  `json:"-"       `
	ProcID  int64  `json:"proc_id" `
	PID     int    `json:"pid"     `
	Name    string `json:"name"    `
	Size    int    `json:"size"    `
	Mime    string `json:"mime"    `
	Time    int64  `json:"time"    `
	Passed  int    `json:"passed"  `
	Failed  int    `json:"failed"  `
	Skipped int    `json:"skipped" `
}

type fileData struct {
	ID      int64  ``
	BuildID int64  ``
	ProcID  int64  ``
	PID     int    ``
	Name    string ``
	Size    int    ``
	Mime    string ``
	Time    int64  ``
	Passed  int    ``
	Failed  int    ``
	Skipped int    ``
	Data    []byte ``
}

func (t File) TableName() string {
	return "cncd_file"
}

func (t fileData) TableName() string {
	return "cncd_file_data"
}

func FileCreate(file *File, r io.Reader) error {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	f := fileData{
		ID:      file.ID,
		BuildID: file.BuildID,
		ProcID:  file.ProcID,
		PID:     file.PID,
		Name:    file.Name,
		Size:    file.Size,
		Mime:    file.Mime,
		Time:    file.Time,
		Passed:  file.Passed,
		Failed:  file.Failed,
		Skipped: file.Skipped,
		Data:    d,
	}
	_, err = x.Insert(&f)

	return err
}
