// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
	"time"
)

//Project for
type Project struct {
	ID          int64
	OwnerID     int64  `xorm:"UNIQUE(s)"`
	Owner       *User  `xorm:"-"`
	LowerName   string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Name        string `xorm:"INDEX NOT NULL"`
	Description string

	ParentID int64
	Parent   *Project `xorm:"-"`

	NumIssues       int
	NumClosedIssues int
	NumOpenIssues   int `xorm:"-"`

	IsPrivate bool

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-"`
	UpdatedUnix int64
}

type CreateProjectOptions struct {
	Name        string
	Description string
	IsPrivate   bool
}

var (
	reservedProjectnames    = []string{"assets", "css", "img", "js", "less", "plugins", "debug", "raw", "install", "api", "avatar", "user", "org", "help", "stars", "issues", "pulls", "commits", "repo", "template", "admin", "new", ".", ".."}
	reservedProjectPatterns = []string{"*.keys"}
)

// CreateProject creates a project for given user
func CreateProject(owner *User, opts CreateProjectOptions) (_ *Project, err error) {

	project := &Project{
		OwnerID:     owner.ID,
		Owner:       owner,
		Name:        opts.Name,
		LowerName:   strings.ToLower(opts.Name),
		Description: opts.Description,
		IsPrivate:   opts.IsPrivate,
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	if err = createProject(sess, owner, project); err != nil {
		return nil, err
	}

	return project, sess.Commit()
}

func createProject(e *xorm.Session, owner *User, project *Project) (err error) {

	if err = IsUsableProjectName(project.Name); err != nil {
		return err
	}

	has, err := isProjectExist(e, owner, project.Name)
	if err != nil {
		return fmt.Errorf("IsProjectExist: %v", err)
	} else if has {
		return fmt.Errorf("IsProjectExist: %v", project.Name)
		//return ErrRepoAlreadyExist{owner.Name, project.Name}
	}

	if _, err = e.Insert(project); err != nil {
		return err
	}

	owner.NumRepos++

	if err = updateUser(e, owner); err != nil {
		return fmt.Errorf("updateUser: %v", err)
	}

	// Give access to all members in owner team.
	if owner.IsOrganization() {
		//t, err := owner.getOwnerTeam(e)
		//if err != nil {
		//	return fmt.Errorf("getOwnerTeam: %v", err)
		//} else if err = t.addRepository(e, project); err != nil {
		//	return fmt.Errorf("addRepository: %v", err)
		//}
	} else {
		// Organization automatically called this in addRepository method.
		//if err = project.recalculateAccesses(e); err != nil {
		//	return fmt.Errorf("recalculateAccesses: %v", err)
		//}
	}

	return nil
}

func IsUsableProjectName(name string) error {
	return isUsableName(reservedProjectnames, reservedProjectPatterns, name)
}

func isProjectExist(e Engine, u *User, projectName string) (bool, error) {
	has, err := e.Get(&Project{
		LowerName: strings.ToLower(projectName),
	})
	return has, err
}


type SearchProjectOptions struct {
	Keyword  string
	OwnerID  int64
	UserID   int64 // When set results will contain all public/private repositories user has access to
	OrderBy  string
	Private  bool // Include private repositories in results
	Page     int
	PageSize int // Can be smaller than or equal to setting.ExplorePagingNum
}

// SearchProjectByName takes keyword and part of project name to search,
// it returns results in given range and number of total results.
func SearchProjectByName(opts *SearchProjectOptions) (repos []*Project, _ int64, _ error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}

	repos = make([]*Project, 0, opts.PageSize)
	sess := x.Alias("prg")
	// Attempt to find repositories that opts.UserID has access to,
	// this does not include other people's private repositories even if opts.UserID is an admin.
	if !opts.Private && opts.UserID > 0 {
		sess.Join("LEFT", "access", "access.repo_id = repo.id").
			Where("(prg.owner_id = ? OR access.user_id = ? OR prg.is_private = ?)", opts.UserID, opts.UserID, false)
	} else {
		// Only return public repositories if opts.Private is not set
		if !opts.Private {
			sess.And("prg.is_private = ?", false)
		}
	}
	if len(opts.Keyword) > 0 {
		sess.And("prg.lower_name LIKE ? OR prg.description LIKE ?", "%"+strings.ToLower(opts.Keyword)+"%", "%"+strings.ToLower(opts.Keyword)+"%")
	}
	if opts.OwnerID > 0 {
		sess.And("prg.owner_id = ?", opts.OwnerID)
	}

	var countSess xorm.Session
	countSess = *sess
	count, err := countSess.Count(new(Project))
	if err != nil {
		return nil, 0, fmt.Errorf("Count: %v", err)
	}

	if len(opts.OrderBy) > 0 {
		sess.OrderBy("prg." + opts.OrderBy)
	}
	return repos, count, sess.Distinct("prg.*").Limit(opts.PageSize, (opts.Page-1)*opts.PageSize).Find(&repos)
}
