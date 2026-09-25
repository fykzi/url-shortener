package model

import "time"

type Link struct {
	ID          int
	ShortCode   string
	OriginalURL string
	Clicks      int

	CreatedAt time.Time
	UpdatedAt time.Time
}
