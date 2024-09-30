package domain

import "time"

type SongDetail struct {
	ID          int       `json:"id" db:"id"`
	SongID      int       `json:"song_id" db:"song_id"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Text        string    `json:"text" db:"text"`
	Link        string    `json:"link" db:"link"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type SongDetailUpdateInput struct {
	ReleaseDate *time.Time `json:"release_date"`
	Text        *string    `json:"text"`
	Link        *string    `json:"link"`
}
