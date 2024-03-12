package calendarerr

import "errors"

var (
	ErrEventNotFound            = errors.New("event not found")
	ErrEventOperationNotAllowed = errors.New("event operation not allowed")
)
