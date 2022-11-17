package nerrs

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("not found")
	ErrUnexpectedValue = errors.New("unexpected value")
	ErrOperationNotAllowed = errors.New("operation not allowed")
)

