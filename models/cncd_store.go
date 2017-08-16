package models

import (
	"io"
)

type Store interface {
	// GetUser gets a user by unique ID.
	GetUser(int64) (*User, error)

	// GetUserLogin gets a user by unique Login name.
	GetUserLogin(string) (*User, error)

	// GetUserList gets a list of all users in the system.
	GetUserList() ([]*User, error)

	// GetUserCount gets a count of all users in the system.
	GetUserCount() (int, error)

	// CreateUser creates a new user account.
	CreateUser(*User) error

	// UpdateUser updates a user account.
	UpdateUser(*User) error

	// DeleteUser deletes a user account.
	DeleteUser(*User) error

	// GetRepo gets a repo by unique ID.
	GetRepo(int64) (*Repository, error)

	// GetRepoName gets a repo by its full name.
	GetRepoName(string) (*Repository, error)

	// GetRepoCount gets a count of all repositories in the system.
	GetRepoCount() (int, error)

	// CreateRepo creates a new repository.
	CreateRepo(*Repository) error

	// UpdateRepo updates a user repository.
	UpdateRepo(*Repository) error

	// DeleteRepo deletes a user repository.
	DeleteRepo(*Repository) error

	// GetBuild gets a build by unique ID.
	GetBuild(int64) (*Build, error)

	// GetBuildNumber gets a build by number.
	GetBuildNumber(*Repository, int) (*Build, error)

	// GetBuildRef gets a build by its ref.
	GetBuildRef(*Repository, string) (*Build, error)

	// GetBuildCommit gets a build by its commit sha.
	GetBuildCommit(*Repository, string, string) (*Build, error)

	// GetBuildLast gets the last build for the branch.
	GetBuildLast(*Repository, string) (*Build, error)

	// GetBuildLastBefore gets the last build before build number N.
	GetBuildLastBefore(*Repository, string, int64) (*Build, error)

	// GetBuildList gets a list of builds for the repository
	GetBuildList(*Repository) ([]*Build, error)

	// GetBuildQueue gets a list of build in queue.
	GetBuildQueue() ([]*Feed, error)

	// GetBuildCount gets a count of all builds in the system.
	GetBuildCount() (int, error)

	// CreateBuild creates a new build and jobs.
	CreateBuild(*Build, ...*Proc) error

	// UpdateBuild updates a build.
	UpdateBuild(*Build) error

	//
	// new functions
	//

	UserFeed(*User) ([]*Feed, error)

	RepoList(*User) ([]*Repository, error)
	RepoListLatest(*User) ([]*Feed, error)
	RepoBatch([]*Repository) error

	PermFind(user *User, repo *Repository) (*Perm, error)
	PermUpsert(perm *Perm) error
	PermBatch(perms []*Perm) error
	PermDelete(perm *Perm) error
	PermFlush(user *User, before int64) error

	ConfigLoad(int64) (*Config, error)
	ConfigFind(*Repository, string) (*Config, error)
	ConfigFindApproved(*Config) (bool, error)
	ConfigCreate(*Config) error

	SenderFind(*Repository, string) (*Sender, error)
	SenderList(*Repository) ([]*Sender, error)
	SenderCreate(*Sender) error
	SenderUpdate(*Sender) error
	SenderDelete(*Sender) error

	SecretFind(*Repository, string) (*Secret, error)
	SecretList(*Repository) ([]*Secret, error)
	SecretCreate(*Secret) error
	SecretUpdate(*Secret) error
	SecretDelete(*Secret) error

	RegistryFind(*Repository, string) (*Registry, error)
	RegistryList(*Repository) ([]*Registry, error)
	RegistryCreate(*Registry) error
	RegistryUpdate(*Registry) error
	RegistryDelete(*Registry) error

	ProcLoad(int64) (*Proc, error)
	ProcFind(*Build, int) (*Proc, error)
	ProcChild(*Build, int, string) (*Proc, error)
	ProcList(*Build) ([]*Proc, error)
	ProcCreate([]*Proc) error
	ProcUpdate(*Proc) error
	ProcClear(*Build) error

	LogFind(*Proc) (io.ReadCloser, error)
	LogSave(*Proc, io.Reader) error

	FileList(*Build) ([]*File, error)
	FileFind(*Proc, string) (*File, error)
	FileRead(*Proc, string) (io.ReadCloser, error)
	FileCreate(*File, io.Reader) error

	TaskList() ([]*Task, error)
	TaskInsert(*Task) error
	TaskDelete(string) error
}
