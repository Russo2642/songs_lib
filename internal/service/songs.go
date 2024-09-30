package service

import (
	"songs-lib/internal/domain"
	"songs-lib/internal/repository"
)

type SongsService struct {
	repo repository.Songs
}

func NewSongsService(repo repository.Songs) *SongsService {
	return &SongsService{repo: repo}
}

func (s *SongsService) Create(songs domain.Songs) (int, error) {
	return s.repo.Create(songs)
}

func (s *SongsService) GetAll(group, song, text, releaseDate, link string, page, limit int) ([]domain.SongsWithDetails, error) {
	return s.repo.GetAll(group, song, text, releaseDate, link, page, limit)
}

func (s *SongsService) Delete(songID int) error {
	return s.repo.Delete(songID)
}

func (s *SongsService) Update(songID int, input domain.SongsUpdateInput) error {
	return s.repo.Update(songID, input)
}
