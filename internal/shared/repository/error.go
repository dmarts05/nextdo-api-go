package repository

import (
	"fmt"
)

type ErrAlreadyExists struct{}

func (e ErrAlreadyExists) Error() string {
	return "already exists"
}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "not found"
}

type ErrFailedToCreate struct {
	Err error
}

func (e ErrFailedToCreate) Error() string {
	return fmt.Sprintf("failed to create: %v", e.Err)
}

type ErrFailedToUpdate struct {
	Err error
}

func (e ErrFailedToUpdate) Error() string {
	return fmt.Sprintf("failed to update: %v", e.Err)
}

type ErrFailedToDelete struct {
	Err error
}

func (e ErrFailedToDelete) Error() string {
	return fmt.Sprintf("failed to delete: %v", e.Err)
}

type ErrFailedToGet struct {
	Err error
}

func (e ErrFailedToGet) Error() string {
	return fmt.Sprintf("failed to get: %v", e.Err)
}
