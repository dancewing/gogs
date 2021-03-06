package models

import (
	"errors"
)

var (
	errEnvironNameInvalid  = errors.New("Invalid Environment Variable Name")
	errEnvironValueInvalid = errors.New("Invalid Environment Variable Value")
)

// EnvironService defines a service for managing environment variables.
type EnvironService interface {
	EnvironList(*Repository) ([]*Environ, error)
}

// EnvironStore persists environment information to storage.
type EnvironStore interface {
	EnvironList(*Repository) ([]*Environ, error)
}

// Environ represents an environment variable.
// swagger:model environ
type Environ struct {
	ID    int64  `json:"id"              `
	Name  string `json:"name"            `
	Value string `json:"value,omitempty" `
}

func (t Environ) TableName() string {
	return "cncd_environ"
}

func EnvironList() ([]*Environ, error) {
	environs := make([]*Environ, 0)

	if err := x.Find(&environs); err != nil {
		return nil, err
	}

	return environs, nil
}

// Validate validates the required fields and formats.
func (e *Environ) Validate() error {
	switch {
	case len(e.Name) == 0:
		return errEnvironNameInvalid
	case len(e.Value) == 0:
		return errEnvironValueInvalid
	default:
		return nil
	}
}

// Copy makes a copy of the environment variable without the value.
func (e *Environ) Copy() *Environ {
	return &Environ{
		ID:   e.ID,
		Name: e.Name,
	}
}
