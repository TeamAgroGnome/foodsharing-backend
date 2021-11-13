package domain

import "time"

type ID uint64

type Object struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt *time.Time
}
