package dto

import (
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Feed) FromDbFeed(feed database.Feed) {
	f.Name = feed.Name
	f.Url = feed.Url
	f.Id = feed.ID
	f.CreatedAt = feed.CreatedAt
	f.UpdatedAt = feed.UpdatedAt
}
