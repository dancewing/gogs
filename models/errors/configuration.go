package errors

import "fmt"

type ConfigurationNotExist struct {
	Type string
	RepositoryID int64
}

func IsConfigurationNotExist(err error) bool {
	_, ok := err.(ConfigurationNotExist)
	return ok
}

func (err ConfigurationNotExist) Error() string {
	return fmt.Sprintf("Configuration does not exist [type: %d, key: %d]", err.Type, err.RepositoryID)
}
