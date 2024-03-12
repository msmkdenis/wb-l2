package mapper

import (
	"errors"
	"strconv"
	"time"

	"dev11/internal/delivery/http/dto"
	"dev11/internal/model"
)

func MapToEvent(e dto.EventCreateRequest) (*model.Event, error) {
	id, err := strconv.ParseUint(e.ID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid event id, must be positive integer")
	}

	creatorId, err := strconv.ParseUint(e.CreatorID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid creator id, must be positive integer")
	}

	date, err := time.Parse("2006-01-02", e.Date)
	if err != nil {
		return nil, errors.New("invalid event date, must be in format 'YYYY-MM-DD'")
	}

	return &model.Event{
		ID:          id,
		CreatorID:   creatorId,
		Name:        e.Name,
		Description: e.Description,
		Date:        date,
	}, nil
}
