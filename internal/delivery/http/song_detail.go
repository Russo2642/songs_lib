package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"songs-lib/internal/domain"
	"strconv"
)

// @Summary Update song detail
// @Description Update song detail by ID
// @Tags songDetail
// @Accept json
// @Produce json
// @Param id path int true "Song Detail ID"
// @Param input body domain.SongDetailUpdateInput true "Update song detail input"
// @Success 200 {object} statusResponse "Success"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 422 {object} errorResponse "Unprocessable Entity"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /song_details/{id} [put]
func (h *Handler) updateSongDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.SongDetailUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.services.SongDetail.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
