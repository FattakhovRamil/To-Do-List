package storage

import (
	"errors"
)

var (
	ErrIncorrectDataFormat = errors.New("incorrect data format")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound = errors.New("task not found")
	OKCreated = 201

)
