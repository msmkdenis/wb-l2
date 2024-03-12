package handlers

import (
	"errors"
	"strconv"
	"time"

	"dev11/internal/delivery/http/dto"
)

type CalendarValidator struct{}

func NewCalendarValidator() *CalendarValidator {
	return &CalendarValidator{}
}

func (v *CalendarValidator) Validate(e dto.EventCreateRequest) error {
	id, err := strconv.ParseUint(e.ID, 10, 64)
	if err != nil {
		return errors.New("invalid event id, must be positive integer")
	}

	if id == 0 {
		return errors.New("invalid event id, must be positive integer")
	}

	creatorId, err := strconv.ParseUint(e.CreatorID, 10, 64)
	if err != nil {
		return errors.New("invalid creator id, must be positive integer")
	}

	if creatorId == 0 {
		return errors.New("invalid creator id, must be positive integer")
	}

	if e.Name == "" {
		return errors.New("invalid event name, must be not empty")
	}

	if e.Description == "" {
		return errors.New("invalid event description, must be not empty")
	}

	_, err = time.Parse("2006-01-02", e.Date)
	if err != nil {
		return errors.New("invalid event date, must be in format 'YYYY-MM-DD'")
	}

	return nil
}
