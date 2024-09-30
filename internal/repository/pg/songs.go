package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"songs-lib/internal/domain"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

type SongsPostgres struct {
	db *sqlx.DB
}

func NewSongs(db *sqlx.DB) *SongsPostgres {
	return &SongsPostgres{db: db}
}

func (r *SongsPostgres) apiRequest(ctx context.Context, group, song string) (*domain.SongDetail, error) {
	logrus.Debugf("Making API request for group: %s, song: %s", group, song)

	apiURL := os.Getenv("API_URL")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		logrus.Errorf("Failed to create new request: %v", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("group", group)
	q.Add("song", song)
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("Failed to execute request: %v", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Bad status from API: %s", resp.Status)
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var songDetail domain.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		logrus.Errorf("Failed to decode response body: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully fetched song detail: %+v", songDetail)
	return &songDetail, nil
}

func (r *SongsPostgres) Create(songs domain.Songs) (int, error) {
	logrus.Infof("Creating song: %+v", songs)

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var songID int
	createSongsQuery := fmt.Sprintf("INSERT INTO %s (\"group\", song) VALUES($1, $2) RETURNING id",
		songsTable)
	row := tx.QueryRow(createSongsQuery, songs.Group, songs.Song)
	if err := row.Scan(&songID); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("failed to rollback transaction: %s", rollbackErr)
		}
		return 0, err
	}

	songDetail, err := r.apiRequest(context.Background(), songs.Group, songs.Song)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("failed to rollback transaction after API called: %s", rollbackErr)
		}
		return 0, err
	}

	songDetail.ID = songID

	if err := r.saveSongDetail(tx, songDetail); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("failed to rollback transaction after saving song detail: %s", rollbackErr)
		}
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("failed to rollback transaction after commit failure: %s", rollbackErr)
		}
		return 0, err
	}

	logrus.Infof("Song created successfully with ID: %d", songID)
	return songID, nil
}

func (r *SongsPostgres) saveSongDetail(tx *sql.Tx, songDetail *domain.SongDetail) error {
	logrus.Debugf("Saving song detail: %+v", songDetail)

	insertQuery := fmt.Sprintf("INSERT INTO %s (song_id, release_date, text, link) VALUES ($1, $2, $3, $4)",
		songsDetailTable)
	_, err := tx.Exec(insertQuery, songDetail.SongID, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
	if err != nil {
		logrus.Errorf("Failed to save song detail to the database: %v", err)
	}

	return err
}

func (r *SongsPostgres) GetAll(group, song, text, releaseDate, link string, page, limit int) ([]domain.SongsWithDetails, error) {
	logrus.Infof("Fetching songs with parameters - Group: %s, Song: %s, Text: %s, ReleaseDate: %s, Link: %s, Page: %d, Limit: %d",
		group, song, text, releaseDate, link, page, limit)

	query := fmt.Sprintf(`SELECT s.id, s."group", s.song, sd.id, sd.song_id, sd.release_date, sd.text, sd.link
		FROM %s s LEFT JOIN %s sd ON s.id = sd.song_id WHERE 1=1`,
		songsTable, songsDetailTable)
	args := []interface{}{}
	argID := 1

	if group != "" {
		query += fmt.Sprintf(" AND s.\"group\" LIKE $%d", argID)
		args = append(args, "%"+group+"%")
		argID++
	}

	if song != "" {
		query += fmt.Sprintf(" AND s.song LIKE $%d", argID)
		args = append(args, "%"+song+"%")
		argID++
	}

	if text != "" {
		query += fmt.Sprintf(" AND sd.text LIKE $%d", argID)
		args = append(args, "%"+text+"%")
		argID++
	}

	if releaseDate != "" {
		query += fmt.Sprintf(" AND sd.release_date = $%d", argID)
		args = append(args, "%"+releaseDate+"%")
		argID++
	}

	if link != "" {
		query += fmt.Sprintf(" AND sd.link LIKE $%d", argID)
		args = append(args, "%"+link+"%")
		argID++
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	logrus.Debugf("Executing query: %s with args: %v", query, args)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var result []domain.SongsWithDetails
	for rows.Next() {
		var songsWithDetails domain.SongsWithDetails
		if err := rows.Scan(&songsWithDetails.Song.ID, &songsWithDetails.Song.Group, &songsWithDetails.Song.Song,
			&songsWithDetails.Detail.ID, &songsWithDetails.Detail.SongID, &songsWithDetails.Detail.ReleaseDate,
			&songsWithDetails.Detail.Text, &songsWithDetails.Detail.Link); err != nil {
			logrus.Errorf("Failed to scan row: %v", err)
			return nil, err
		}
		result = append(result, songsWithDetails)
	}

	logrus.Infof("Fetched %d songs", len(result))
	return result, nil
}

func (r *SongsPostgres) Update(songID int, input domain.SongsUpdateInput) error {
	logrus.Infof("Updating song with ID: %d, Input: %+v", songID, input)

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Group != nil {
		setValues = append(setValues, fmt.Sprintf(`"group"=$%d`, argId))
		args = append(args, *input.Group)
		argId++
	}

	if input.Song != nil {
		setValues = append(setValues, fmt.Sprintf("song=$%d", argId))
		args = append(args, *input.Song)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d",
		songsTable, setQuery, argId)
	args = append(args, songID)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		logrus.Infof("Successfully updated song with ID: %d", songID)
		return err
	}

	logrus.Infof("Successfully updated song with ID: %d", songID)
	return nil
}

func (r *SongsPostgres) Delete(songID int) error {
	logrus.Infof("Deleting song with ID: %d", songID)

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", songsTable)
	_, err := r.db.Exec(query, songID)
	if err != nil {
		logrus.Errorf("Failed to delete song: %v", err)
		return err
	}

	logrus.Infof("Successfully deleted song with ID: %d", songID)
	return nil
}
