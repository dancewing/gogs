package models

import (
	"fmt"

	"github.com/kataras/iris/core/errors"
)

// ConfigStore persists pipeline configuration to storage.
type ConfigStore interface {
	ConfigLoad(int64) (*Config, error)
	ConfigFind(*Repository, string) (*Config, error)
	ConfigFindApproved(*Config) (bool, error)
	ConfigCreate(*Config) error
}

// Config represents a pipeline configuration.
type Config struct {
	ID     int64  `json:"-"    `
	RepoID int64  `json:"-"    `
	Data   string `json:"data" `
	Hash   string `json:"hash" `
}

func (t Config) TableName() string {
	return "cncd_config"
}

func ConfigFind(repo *Repository, sha string) (*Config, error) {

	confs := make([]*Config, 0)

	if err := x.Where("repo_id = ? and hash = ?", repo.ID, sha).Find(&confs); err != nil {
		return nil, err
	}

	if len(confs) > 0 {
		return confs[0], nil
	}

	return nil, errors.New("not found config")
}

func GetConfigByID(id int64) (*Config, error) {
	cnf := new(Config)
	has, err := x.Id(id).Get(cnf)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(fmt.Sprintf("cant load config by ID :%d", id))
	}
	return cnf, nil
}

func ConfigCreate(conf *Config) (err error) {
	_, err = x.InsertOne(conf)
	return err
}
