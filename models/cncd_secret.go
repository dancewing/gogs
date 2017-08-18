package models

import (
	"errors"
	"path/filepath"
)

var (
	errSecretNameInvalid  = errors.New("Invalid Secret Name")
	errSecretValueInvalid = errors.New("Invalid Secret Value")
)

// Secret represents a secret variable, such as a password or token.
// swagger:model registry
type Secret struct {
	ID         int64    `json:"id"              `
	RepoID     int64    `json:"-"               `
	Name       string   `json:"name"            `
	Value      string   `json:"value,omitempty" `
	Images     []string `json:"-"        xorm:"JSON"`
	Events     []string `json:"-"        xorm:"JSON"`
	SkipVerify bool     `json:"-"               `
	Conceal    bool     `json:"-"               `
}

func (t Secret) TableName() string {
	return "cncd_secret"
}

// Match returns true if an image and event match the restricted list.
func (s *Secret) Match(event string) bool {
	if len(s.Events) == 0 {
		return true
	}
	for _, pattern := range s.Events {
		if match, _ := filepath.Match(pattern, event); match {
			return true
		}
	}
	return false
}

// Validate validates the required fields and formats.
func (s *Secret) Validate() error {
	switch {
	case len(s.Name) == 0:
		return errSecretNameInvalid
	case len(s.Value) == 0:
		return errSecretValueInvalid
	default:
		return nil
	}
}

// Copy makes a copy of the secret without the value.
func (s *Secret) Copy() *Secret {
	return &Secret{
		ID:     s.ID,
		RepoID: s.RepoID,
		Name:   s.Name,
		Images: s.Images,
		Events: s.Events,
	}
}

func SecretListBuild(repo *Repository) ([]*Secret, error) {
	secs := make([]*Secret, 0)

	if err := x.Where("repo_id = ? ", repo.ID).Find(&secs); err != nil {
		return nil, err
	}

	return secs, nil
}
