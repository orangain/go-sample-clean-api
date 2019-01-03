package web

import (
	"net/http"
	"strconv"

	"github.com/orangain/clean-api/domain"

	"github.com/labstack/echo"
)

type FilmEchoHandlers struct {
	repo domain.FilmRepository
}

func SetupFilmEchoHandlers(e *echo.Echo, repo domain.FilmRepository) {
	h := &FilmEchoHandlers{repo: repo}
	e.GET("/films", h.GetFilms)
	e.GET("/films/:id", h.GetFilm)
	e.POST("/films", h.CreateFilm)
	e.DELETE("/films/:id", h.DeleteFilm)
	e.PUT("/films/:id", h.UpdateFilm)
}

func (h *FilmEchoHandlers) GetFilms(c echo.Context) error {
	films, err := h.repo.GetFilms()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, films)
}

func (h *FilmEchoHandlers) GetFilm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.ErrNotFound
	}
	film, err := h.repo.GetFilm(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, film)
}

func (h *FilmEchoHandlers) CreateFilm(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Content-Type must be application/json")
	}
	var film domain.Film
	err := c.Bind(&film)
	if err != nil {
		return err
	}
	err = c.Validate(film)
	if err != nil {
		return err
	}
	c.Logger().Info(film)
	f, err := h.repo.InsertFilm(&film)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Location", c.Echo().URI(h.GetFilm, f.FilmID))
	return c.JSON(http.StatusCreated, f)
}

func (h *FilmEchoHandlers) DeleteFilm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.ErrNotFound
	}
	err = h.repo.DeleteFilm(id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *FilmEchoHandlers) UpdateFilm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return domain.ErrNotFound
	}
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Content-Type must be application/json")
	}

	var film domain.Film
	err = c.Bind(&film)
	if err != nil {
		return err
	}
	err = c.Validate(film)
	if err != nil {
		return err
	}
	f, err := h.repo.UpdateFilm(id, &film)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, f)
}
