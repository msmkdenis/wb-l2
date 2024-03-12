package calendarerr

import "errors"

var (
	ErrEventNotFound         = errors.New("event not found")
	ErrEventUpdateNotAllowed = errors.New("update event not allowed")
)
