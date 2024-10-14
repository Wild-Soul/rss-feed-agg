package dto

import (
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

type FeedFollows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (ff *FeedFollows) FromDbFeed(feedFollow database.FeedFollow) {
	ff.ID = feedFollow.ID
	ff.UserID = feedFollow.UserID
	ff.FeedID = feedFollow.FeedID
	ff.CreatedAt = feedFollow.CreatedAt
	ff.UpdatedAt = feedFollow.UpdatedAt
}
