package server

import (
	"fmt"
)

type Error struct {
	originalErr error
}

func (e Error) Error() string {
	return fmt.Sprintf("Game server error: %v", e.originalErr)
}

func (e Error) Unwrap() error {
	return e.originalErr
}

func wrapServerErr(err error) error {
	return Error{originalErr: err}
}
