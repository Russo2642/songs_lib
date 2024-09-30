package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"songs-lib/internal/domain"
	"strings"
)

type SongDetailPostgres struct {
	db *sqlx.DB
}

func NewSongDetail(db *sqlx.DB) *SongDetailPostgres {
	return &SongDetailPostgres{db: db}
}

func (r *SongDetailPostgres) Update(songDetailID int, input domain.SongDetailUpdateInput) error {
	logrus.Infof("Updating song detail with ID: %d", songDetailID)

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.ReleaseDate != nil {
		setValues = append(setValues, fmt.Sprintf("release_date=$%d", argId))
		args = append(args, *input.ReleaseDate)
		argId++
	}

	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, *input.Text)
		argId++
	}

	if input.Link != nil {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argId))
		args = append(args, *input.Link)
		argId++
	}

	if len(setValues) == 0 {
		logrus.Warn("No fields to update for song detail ID: %d", songDetailID)
		return fmt.Errorf("no fields to update")
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d",
		songsDetailTable, setQuery, argId)
	args = append(args, songDetailID)

	logrus.Debugf("Executing update query: %s with args: %v", query, args)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		logrus.Errorf("Failed to update song detail with ID %d: %v", songDetailID, err)
		return err
	}

	logrus.Infof("Successfully updated song detail with ID: %d", songDetailID)
	return nil
}
