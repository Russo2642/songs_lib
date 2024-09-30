package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"songs-lib/internal/domain"
	"strconv"
)

// createSong godoc
// @Summary Create a new song
// @Description Create a new song record by sending an API request with group and song data to retrieve song details
// @Tags songs
// @Accept json
// @Produce json
// @Param input body domain.Songs true "Song data"
// @Success 201 {object} statusResponse "Successfully created song"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /songs [post]
func (h *Handler) createSong(c *gin.Context) {
	var input domain.Songs
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Songs.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"status":  http.StatusCreated,
		"message": "Song created",
	})
}

// getAllSongs godoc
// @Summary Get all songs
// @Description Retrieve a list of songs with optional filters
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group filter"
// @Param song query string false "Song filter"
// @Param text query string false "Text filter"
// @Param release_date query string false "Release date filter"
// @Param link query string false "Link filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit number" default(10)
// @Success 200 {array} domain.SongsWithDetails "List of songs"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /songs [get]
func (h *Handler) getAllSongs(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	text := c.Query("text")
	releaseDate := c.Query("release_date")
	link := c.Query("link")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	songs, err := h.services.Songs.GetAll(group, song, text, releaseDate, link, page, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, songs)
}

// updateSong godoc
// @Summary Update an existing song
// @Description Update a song by ID with the given input data
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param input body domain.SongsUpdateInput true "Updated song data"
// @Success 200 {object} statusResponse "Successfully updated song"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /songs/{id} [put]
func (h *Handler) updateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.SongsUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Songs.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// deleteSong godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} statusResponse "Successfully deleted song"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *Handler) deleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Songs.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
