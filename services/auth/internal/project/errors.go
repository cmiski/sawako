package project

import "errors"

var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrNameRequired      = errors.New("name is required")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
)
