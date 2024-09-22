package dto

import (
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func (userdto *UserDTO) FromDbUser(dbUser database.User) {
	userdto.ID = dbUser.ID
	userdto.CreatedAt = dbUser.CreatedAt
	userdto.UpdatedAt = dbUser.UpdatedAt
	userdto.Name = dbUser.Name
	userdto.ApiKey = dbUser.ApiKey
}
