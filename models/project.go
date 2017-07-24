// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/process"
)

type CreateProjectOptions struct {
	Name          string
	Description   string
	Gitignores    string
	License       string
	Readme        string
	IsPrivate     bool
	IsMirror      bool
	AutoInit      bool
	CreateGitRepo bool
}

// CreateProject creates a project for given user or organization. project can be lazy initialize a repository
func CreateProject(doer, owner *User, opts CreateProjectOptions) (_ *Repository, err error) {
	if !owner.CanCreateRepo() {
		return nil, errors.ReachLimitOfRepo{owner.RepoCreationNum()}
	}

	repo := &Repository{
		OwnerID:        owner.ID,
		Owner:          owner,
		Name:           opts.Name,
		LowerName:      strings.ToLower(opts.Name),
		Description:    opts.Description,
		IsPrivate:      opts.IsPrivate,
		EnableWiki:     true,
		EnableIssues:   true,
		EnablePulls:    true,
		GitInitialized: opts.CreateGitRepo,
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	if err = createProject(sess, doer, owner, repo); err != nil {
		return nil, err
	}

	// No need for init mirror.
	if !opts.IsMirror && opts.CreateGitRepo {
		err = initializeGit(sess, doer, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("CreateRepository 'git update-server-info': %s", err.Error())
		}
	}

	return repo, sess.Commit()
}

//ListUserTopProjects TODO
func ListUserTopProjects(owner *User) ([]*Repository, error) {
	return GetUserRepositories(&UserRepoOptions{
		UserID: owner.ID,
	})
}

func (r *Repository) InitializeGit(doer *User, opts CreateProjectOptions)(err error ) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	r.LoadAttributes()

	r.GitInitialized = true
	sess.Cols("git_initialized").Update(r)

	err = initializeGit(sess, doer, r.Owner, r, opts)
	if err != nil {
		return fmt.Errorf("initializeGit 'git update-server-info': %s", err.Error())
	}

	return sess.Commit()
}

func initializeGit(sess *xorm.Session, doer, owner *User, repo *Repository, opts CreateProjectOptions) (err error) {

	repoPath := RepoPath(owner.Name, repo.Name)
	createRepoOptions := CreateRepoOptions{
		Name:        opts.Name,
		Description: opts.Description,
		Gitignores:  opts.Gitignores,
		License:     opts.License,
		Readme:      opts.Readme,
		IsPrivate:   opts.IsPrivate,
		IsMirror:    opts.IsMirror,
		AutoInit:    opts.AutoInit,
	}

	if err = initRepository(sess, repoPath, doer, repo, createRepoOptions); err != nil {
		RemoveAllWithNotice("Delete repository for initialization failure", repoPath)
		return fmt.Errorf("initRepository: %v", err)
	}

	_, stderr, err := process.ExecDir(-1,
		repoPath, fmt.Sprintf("CreateRepository 'git update-server-info': %s", repoPath),
		"git", "update-server-info")
	if err != nil {
		return fmt.Errorf("CreateRepository 'git update-server-info': %s", stderr)
	}
	return nil
}

func createProject(e *xorm.Session, doer, owner *User, repo *Repository) (err error) {
	if err = IsUsableRepoName(repo.Name); err != nil {
		return err
	}

	has, err := isRepositoryExist(e, owner, repo.Name)
	if err != nil {
		return fmt.Errorf("IsRepositoryExist: %v", err)
	} else if has {
		return ErrRepoAlreadyExist{owner.Name, repo.Name}
	}

	if _, err = e.Insert(repo); err != nil {
		return err
	}

	owner.NumRepos++
	// Remember visibility preference.
	owner.LastRepoVisibility = repo.IsPrivate
	if err = updateUser(e, owner); err != nil {
		return fmt.Errorf("updateUser: %v", err)
	}

	// Give access to all members in owner team.
	if owner.IsOrganization() {
		t, err := owner.getOwnerTeam(e)
		if err != nil {
			return fmt.Errorf("getOwnerTeam: %v", err)
		} else if err = t.addRepository(e, repo); err != nil {
			return fmt.Errorf("addRepository: %v", err)
		}
	} else {
		// Organization automatically called this in addRepository method.
		if err = repo.recalculateAccesses(e); err != nil {
			return fmt.Errorf("recalculateAccesses: %v", err)
		}
	}

	if err = watchRepo(e, owner.ID, repo.ID, true); err != nil {
		return fmt.Errorf("watchRepo: %v", err)
	} else if err = newRepoAction(e, doer, owner, repo); err != nil {
		return fmt.Errorf("newRepoAction: %v", err)
	}

	return repo.loadAttributes(e)
}
