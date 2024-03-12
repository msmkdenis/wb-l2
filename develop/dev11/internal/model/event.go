package model

import (
	"time"
)

type Event struct {
	ID          uint64    `json:"id,omitempty"`
	CreatorID   uint64    `json:"creator_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
}
