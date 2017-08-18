package models

import (
	"errors"
	"fmt"
)

// swagger:model build
type Build struct {
	ID          int64  `json:"id"            `
	RepoID      int64  `json:"-"             `
	ConfigID    int64  `json:"-"             `
	Number      int    `json:"number"        `
	Parent      int    `json:"parent"        `
	Event       string `json:"event"         `
	Status      string `json:"status"        `
	Error       string `json:"error"         `
	Enqueued    int64  `json:"enqueued_at"   `
	Created     int64  `json:"created_at"    `
	Started     int64  `json:"started_at"    `
	Finished    int64  `json:"finished_at"   `
	Deploy      string `json:"deploy_to"     `
	Commit      string `json:"commit"        `
	Branch      string `json:"branch"        `
	Ref         string `json:"ref"           `
	Refspec     string `json:"refspec"       `
	Remote      string `json:"remote"        `
	Title       string `json:"title"         `
	Message     string `json:"message"       `
	Timestamp   int64  `json:"timestamp"     `
	Sender      string `json:"sender"        `
	Author      string `json:"author"        `
	Avatar      string `json:"author_avatar" `
	Email       string `json:"author_email"  `
	Link        string `json:"link_url"      `
	Signed      bool   `json:"signed"        ` // deprecate
	Verified    bool   `json:"verified"      ` // deprecate
	Reviewer    string `json:"reviewed_by"   `
	Reviewed    int64  `json:"reviewed_at"   `
	IsDelivered bool
	Procs       []*Proc `json:"procs,omitempty" xorm:"-"`
	Files       []*File `json:"files,omitempty" xorm:"-"`
}

func (t Build) TableName() string {
	return "cncd_build"
}

// Trim trims string values that would otherwise exceed
// the database column sizes and fail to insert.
func (b *Build) Trim() {
	if len(b.Title) > 1000 {
		b.Title = b.Title[:1000]
	}
	if len(b.Message) > 2000 {
		b.Message = b.Message[:2000]
	}
}

func CreateBuild(build *Build, procs ...*Proc) (err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	_, err = sess.Insert(build)

	if err != nil {
		return nil
	}

	_, err = sess.Insert(procs)

	if err != nil {
		return nil
	}

	return sess.Commit()
}

func GetBuild(id int64) (*Build, error) {
	build := new(Build)
	has, err := x.Id(id).Get(build)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(fmt.Sprintf("Build with ID %d not exisit ", id))
	}
	return build, nil
}

func UpdateBuild(build *Build) error {
	if _, err := x.ID(build.ID).AllCols().Update(build); err != nil {
		return err
	}
	return nil
}

func GetBuildLastBefore(repo *Repository, branch string, num int64) (*Build, error) {

	builds := make([]*Build, 0)

	if err := x.Where("repo_id = ? and branch = ? and id < ?", repo.ID, branch, num).Limit(1, 0).Find(&builds); err != nil {
		return nil, err
	}

	if len(builds) > 0 {
		return builds[0], nil
	}
	return nil, nil
}

func CountBuild(repositoryID int64) int64 {
	sess := x.Where("id > 0")

	if repositoryID > 0 {
		sess.And("repo_id = ?", repositoryID)
	}

	count, _ := sess.Count(new(Build))
	//if err != nil {
	//	log.Error(4, "CountBuild: %v", err)
	//}
	return count
}

func ListBuilds(opts *BuildOptions) ([]*Build, error) {
	sess := x.Where("repo_id=?", opts.RepositoryID).Desc("created")
	if opts.Page <= 0 {
		opts.Page = 1
	}
	sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize)

	builds := make([]*Build, 0, opts.PageSize)
	return builds, sess.Find(&builds)
}

type BuildOptions struct {
	RepositoryID int64
	Page         int
	PageSize     int
}
