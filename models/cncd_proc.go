package models

import (
	"fmt"

	"github.com/kataras/iris/core/errors"
)

// ProcStore persists process information to storage.
type ProcStore interface {
	ProcLoad(int64) (*Proc, error)
	ProcFind(*Build, int) (*Proc, error)
	ProcChild(*Build, int, string) (*Proc, error)
	ProcList(*Build) ([]*Proc, error)
	ProcCreate([]*Proc) error
	ProcUpdate(*Proc) error
	ProcClear(*Build) error
}

// Proc represents a process in the build pipeline.
// swagger:model proc
type Proc struct {
	ID       int64             `json:"id"                   `
	BuildID  int64             `json:"build_id"             `
	PID      int               `json:"pid"                  `
	PPID     int               `json:"ppid"                 `
	PGID     int               `json:"pgid"                 `
	Name     string            `json:"name"                 `
	State    string            `json:"state"                `
	Error    string            `json:"error,omitempty"      xorm:"TEXT"`
	ExitCode int               `json:"exit_code"            `
	Started  int64             `json:"start_time,omitempty" `
	Stopped  int64             `json:"end_time,omitempty"   `
	Machine  string            `json:"machine,omitempty"    `
	Platform string            `json:"platform,omitempty"   `
	Environ  map[string]string `json:"environ,omitempty"    xorm:"JSON"`
	Children []*Proc           `json:"children,omitempty"   xorm:"-"`
}

func (t Proc) TableName() string {
	return "cncd_proc"
}

// Running returns true if the process state is pending or running.
func (p *Proc) Running() bool {
	return p.State == StatusPending || p.State == StatusRunning
}

// Failing returns true if the process state is failed, killed or error.
func (p *Proc) Failing() bool {
	return p.State == StatusError || p.State == StatusKilled || p.State == StatusFailure
}

// Tree creates a process tree from a flat process list.
func Tree(procs []*Proc) []*Proc {
	var (
		nodes  []*Proc
		parent *Proc
	)
	for _, proc := range procs {
		if proc.PPID == 0 {
			nodes = append(nodes, proc)
			parent = proc
			continue
		} else {
			parent.Children = append(parent.Children, proc)
		}
	}
	return nodes
}

func ProcCreate(procs []*Proc) (err error) {

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	for _, p := range procs {
		_, err = sess.InsertOne(p)
		if err != nil {
			return nil
		}
	}

	return sess.Commit()
}

func ProcList(build *Build) ([]*Proc, error) {

	procs := make([]*Proc, 0)

	if err := x.Where("build_id = ? ", build.ID).Find(&procs); err != nil {
		return nil, err
	}

	return procs, nil
}

func ProcUpdate(proc *Proc) (err error) {
	_, err = x.ID(proc.ID).AllCols().Update(proc)
	return err
}

func ProcLoad(id int64) (*Proc, error) {
	proc := new(Proc)
	has, err := x.Id(id).Get(proc)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(fmt.Sprintf("Proc with ID %d not exisit ", id))
	}
	return proc, nil
}

func ProcChild(build *Build, pid int, child string) (*Proc, error) {
	procs := make([]*Proc, 0)

	err := x.Where("build_id = ? and ppid = ? and name = ?", build.ID, pid, child).Limit(1, 0).Find(&procs)

	if err != nil {
		return nil, err
	}

	if len(procs) > 0 {
		return procs[0], nil
	}
	return nil, errors.New(fmt.Sprintf("Proc not found with name %s ", child))
}
