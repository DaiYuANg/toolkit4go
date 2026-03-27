package repository

import "fmt"

// ErrNotFound reports that a repository entity does not exist.
var ErrNotFound = &repositoryError{"not found"}

// ErrOperationNotSupported reports that the selected backend cannot perform the requested operation.
var ErrOperationNotSupported = &repositoryError{"operation not supported"}

// ErrFieldNotFound reports that the requested entity field does not exist in repository metadata.
var ErrFieldNotFound = &repositoryError{"field not found"}

type repositoryError struct{ msg string }

func (e *repositoryError) Error() string { return "kvx: " + e.msg }

func wrapRepositoryError(err error, action string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	return nil
}

func wrapRepositoryResult[T any](value T, err error, action string) (T, error) {
	if err != nil {
		var zero T
		return zero, fmt.Errorf("%s: %w", action, err)
	}

	return value, nil
}
