package controllers

import (
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getLeaderboardRatingRequest struct {
	Top int `query:"top" validate:"numeric,max=100,min=1" default:"10"`
}

type getLeaderboardRatingResponse struct {
	Names   []string `json:"names" validate:"required"`
	Ratings []int    `json:"ratings" validate:"required"`
}

func (h *Handlers) GetLeaderboardRating(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[getLeaderboardRatingRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	players, err := h.playerService.GetTopPlayersByRating(ctx, req.Top)

	names := make([]string, len(players))
	ratings := make([]int, len(players))
	for i, p := range players {
		names[i] = p.Name
		ratings[i] = p.Rating
	}

	resp := getLeaderboardRatingResponse{
		Names:   names,
		Ratings: ratings,
	}

	return c.JSON(http.StatusOK, resp)
}
