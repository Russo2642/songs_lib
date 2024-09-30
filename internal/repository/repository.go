package repository

import (
	"github.com/jmoiron/sqlx"
	"songs-lib/internal/domain"
	"songs-lib/internal/repository/pg"
)

type Songs interface {
	Create(songs domain.Songs) (int, error)
	GetAll(group, song, text, releaseDate, link string, page, limit int) ([]domain.SongsWithDetails, error)
	Update(songID int, input domain.SongsUpdateInput) error
	Delete(songID int) error
}

type SongDetail interface {
	Update(songDetailID int, input domain.SongDetailUpdateInput) error
}

type Repository struct {
	Songs
	SongDetail
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Songs:      pg.NewSongs(db),
		SongDetail: pg.NewSongDetail(db),
	}
}
