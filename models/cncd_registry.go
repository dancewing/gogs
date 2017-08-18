package models

import "errors"

var (
	errRegistryAddressInvalid  = errors.New("Invalid Registry Address")
	errRegistryUsernameInvalid = errors.New("Invalid Registry Username")
	errRegistryPasswordInvalid = errors.New("Invalid Registry Password")
)

// RegistryService defines a service for managing registries.
type RegistryService interface {
	RegistryFind(*Repository, string) (*Registry, error)
	RegistryList(*Repository) ([]*Registry, error)
	RegistryCreate(*Repository, *Registry) error
	RegistryUpdate(*Repository, *Registry) error
	RegistryDelete(*Repository, string) error
}

// RegistryStore persists registry information to storage.
type RegistryStore interface {
	RegistryFind(*Repository, string) (*Registry, error)
	RegistryList(*Repository) ([]*Registry, error)
	RegistryCreate(*Registry) error
	RegistryUpdate(*Registry) error
	RegistryDelete(*Registry) error
}

// Registry represents a docker registry with credentials.
// swagger:model registry
type Registry struct {
	ID       int64  `json:"id"       `
	RepoID   int64  `json:"-"        `
	Address  string `json:"address"  `
	Username string `json:"username" `
	Password string `json:"password" `
	Email    string `json:"email"    `
	Token    string `json:"token"    `
}

func (t Registry) TableName() string {
	return "cncd_registry"
}

// Validate validates the registry information.
func (r *Registry) Validate() error {
	switch {
	case len(r.Address) == 0:
		return errRegistryAddressInvalid
	case len(r.Username) == 0:
		return errRegistryUsernameInvalid
	case len(r.Password) == 0:
		return errRegistryPasswordInvalid
	default:
		return nil
	}
}

// Copy makes a copy of the registry without the password.
func (r *Registry) Copy() *Registry {
	return &Registry{
		ID:       r.ID,
		RepoID:   r.RepoID,
		Address:  r.Address,
		Username: r.Username,
		Email:    r.Email,
		Token:    r.Token,
	}
}


func RegistryList(repo *Repository) ([]*Registry, error) {
	regs := make([]*Registry, 0)

	if err := x.Where("repo_id = ? ", repo.ID ).Find(&regs); err != nil {
		return nil, err
	}

	return regs, nil
}
