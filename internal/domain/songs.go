package domain

import "time"

type Songs struct {
	ID        int       `json:"id" db:"id"`
	Group     string    `json:"group" db:"group" binding:"required"`
	Song      string    `json:"song" db:"song" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SongsWithDetails struct {
	Song   Songs      `json:"song"`
	Detail SongDetail `json:"detail"`
}

type SongsUpdateInput struct {
	Group *string `json:"group"`
	Song  *string `json:"song"`
}
