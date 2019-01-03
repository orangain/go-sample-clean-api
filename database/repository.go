package database

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
	"github.com/orangain/clean-api/domain"
)

type FilmSqlRepository struct {
	db *sql.DB
}

func NewFilmSqlRepository(db *sql.DB) domain.FilmRepository {
	return &FilmSqlRepository{db: db}
}

func (r *FilmSqlRepository) GetFilms() ([]*domain.Film, error) {
	rows, err := r.db.Query(`
		SELECT film_id, title, description, release_year, language_id
		FROM film`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	films, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (r *FilmSqlRepository) GetFilm(id int) (*domain.Film, error) {
	rows, err := r.db.Query(`
		SELECT film_id, title, description, release_year, language_id
		FROM film
		WHERE film_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	films, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(films) < 1 {
		return nil, domain.ErrNotFound
	}
	return films[0], nil
}

func scanRows(rows *sql.Rows) ([]*domain.Film, error) {
	films := make([]*domain.Film, 0)
	for rows.Next() {
		var film domain.Film
		err := rows.Scan(
			&film.FilmID,
			&film.Title,
			&film.Description,
			&film.ReleaseYear,
			&film.LanguageID,
		)
		if err != nil {
			return nil, err
		}
		films = append(films, &film)
	}
	return films, nil
}

func (r *FilmSqlRepository) InsertFilm(film *domain.Film) (*domain.Film, error) {
	var id int
	err := transact(r.db, func(tx *sql.Tx) error {
		err := tx.QueryRow(`
			INSERT INTO film (title, description, release_year, language_id)
			VALUES ($1, $2, $3, $4)
			RETURNING film_id`,
			film.Title, film.Description, film.ReleaseYear, film.LanguageID,
		).Scan(&id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return r.GetFilm(id)
}

func (r *FilmSqlRepository) DeleteFilm(id int) error {
	err := transact(r.db, func(tx *sql.Tx) error {
		var deletedID int
		err := tx.QueryRow(`
			DELETE FROM film
			WHERE film_id = $1
			RETURNING film_id`, id,
		).Scan(&deletedID)
		return err
	})
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() == "23" {
			return domain.ErrConflict
		}
	}
	switch err {
	case sql.ErrNoRows:
		return domain.ErrNotFound
	}
	return err
}

func (r *FilmSqlRepository) UpdateFilm(id int, film *domain.Film) (*domain.Film, error) {
	err := transact(r.db, func(tx *sql.Tx) error {
		var updatedID int
		err := tx.QueryRow(`
			UPDATE film SET
				title = $2,
				description = $3,
				release_year = $4,
				language_id = $5
			WHERE film_id = $1
			RETURNING film_id`,
			id, film.Title, film.Description, film.ReleaseYear, film.LanguageID,
		).Scan(&updatedID)
		return err
	})
	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return r.GetFilm(id)
}

func transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	log.Println("Transcation started")
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			log.Println("Transcation rollbacked")
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback()
			log.Println("Transcation rollbacked")
		} else {
			err = tx.Commit()
			log.Println("Transcation commited")
		}
	}()
	err = txFunc(tx)
	return err
}
