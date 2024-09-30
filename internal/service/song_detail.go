package service

import (
	"songs-lib/internal/domain"
	"songs-lib/internal/repository"
)

type SongDetailService struct {
	repo repository.SongDetail
}

func NewSongDetailService(repo repository.SongDetail) *SongDetailService {
	return &SongDetailService{repo: repo}
}

func (s *SongDetailService) Update(songDetailID int, input domain.SongDetailUpdateInput) error {
	return s.repo.Update(songDetailID, input)
}
