package service

import (
	"songs-lib/internal/domain"
	"songs-lib/internal/repository"
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

type Service struct {
	Songs
	SongDetail
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs:      NewSongsService(repos.Songs),
		SongDetail: NewSongDetailService(repos.SongDetail),
	}
}
