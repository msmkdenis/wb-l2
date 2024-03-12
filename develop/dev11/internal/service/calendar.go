package service

import (
	"context"
	"time"

	"dev11/internal/model"
)

type CalendarRepository interface {
	Insert(ctx context.Context, e *model.Event) (*model.Event, error)
	Update(ctx context.Context, e *model.Event) (*model.Event, error)
	Delete(ctx context.Context, id, creatorID uint64) error
	SelectForDay(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
	SelectForMonth(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
	SelectForYear(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
}

type Calendar struct {
	repository CalendarRepository
}

func NewCalendar(r CalendarRepository) *Calendar {
	return &Calendar{
		repository: r,
	}
}

func (c *Calendar) Add(ctx context.Context, e *model.Event) (*model.Event, error) {
	return c.repository.Insert(ctx, e)
}

func (c *Calendar) Update(ctx context.Context, e *model.Event) (*model.Event, error) {
	return c.repository.Update(ctx, e)
}

func (c *Calendar) GetEventsForDay(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	return c.repository.SelectForDay(ctx, creatorId, date)
}

func (c *Calendar) GetEventsForMonth(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	return c.repository.SelectForMonth(ctx, creatorId, date)
}

func (c *Calendar) GetEventsForYear(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error) {
	return c.repository.SelectForYear(ctx, creatorId, date)
}

func (c *Calendar) Delete(ctx context.Context, id, creatorID uint64) error {
	return c.repository.Delete(ctx, id, creatorID)
}
