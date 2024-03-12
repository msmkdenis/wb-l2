package repository

import (
	"context"
	"sync"
	"time"

	"dev11/internal/calendarerr"
	"dev11/internal/model"
)

type Calendar struct {
	storage map[uint64]model.Event
	mu      sync.RWMutex
}

func NewCalendar() *Calendar {
	return &Calendar{
		storage: make(map[uint64]model.Event),
		mu:      sync.RWMutex{},
	}
}

func (c *Calendar) Insert(ctx context.Context, e *model.Event) (*model.Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[e.ID] = *e
	return e, nil
}

func (c *Calendar) Update(ctx context.Context, e *model.Event) (*model.Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	event, ok := c.storage[e.ID]
	if !ok {
		return nil, calendarerr.ErrEventNotFound
	}

	if e.CreatorID != event.CreatorID {
		return nil, calendarerr.ErrEventOperationNotAllowed
	}

	c.storage[e.ID] = *e
	return e, nil
}

func (c *Calendar) Delete(ctx context.Context, id, creatorID uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	event, ok := c.storage[id]
	if !ok {
		return calendarerr.ErrEventNotFound
	}

	if event.CreatorID != creatorID {
		return calendarerr.ErrEventOperationNotAllowed
	}

	delete(c.storage, id)
	return nil
}

func (c *Calendar) SelectForDay(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	events := make([]model.Event, 0)
	for _, v := range c.storage {
		if v.CreatorID == creatorId && v.Date.Day() == date.Day() {
			events = append(events, v)
		}
	}

	return events, nil
}

func (c *Calendar) SelectForMonth(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	events := make([]model.Event, 0)
	for _, v := range c.storage {
		if v.CreatorID == creatorId && v.Date.Month() == date.Month() {
			events = append(events, v)
		}
	}

	return events, nil
}

func (c *Calendar) SelectForYear(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	events := make([]model.Event, 0)
	for _, v := range c.storage {
		if v.CreatorID == creatorId && v.Date.Year() == date.Year() {
			events = append(events, v)
		}
	}

	return events, nil
}
