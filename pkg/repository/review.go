package repository

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReviewPostgres struct {
	db *sqlx.DB
}

func NewReviewPostgres(db *sqlx.DB) *ReviewPostgres {
	return &ReviewPostgres{db: db}
}

func (r *ReviewPostgres) CreateReview(review models.Review, userId uuid.UUID) error {
	query := fmt.Sprintln(
		`
		INSERT INTO reviews (id, body, film_title, author_id, rating)
		VALUES ($1, $2, $3, $4, $5)
		`,
	)
	_, err := r.db.Exec(query, review.Id, review.Body, review.FilmTitle, userId, review.Rating)
	return err
}

func (r *ReviewPostgres) EditReview(new models.ReviewInput, userId, reviewId uuid.UUID) error {
	fields := make([]string, 0)
	values := make([]interface{}, 0)

	rv := reflect.ValueOf(new)
	trueCounter := 1
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			fieldName := rv.Type().Field(i).Tag.Get("db")
			fields = append(fields, fmt.Sprintf("%s=$%d", fieldName, trueCounter))
			values = append(values, field.Elem().Interface())
			trueCounter += 1
		}
	}

	query := fmt.Sprintf("UPDATE reviews SET %s WHERE author_id = $%d and id = $%d", strings.Join(fields, ", "), trueCounter, trueCounter+1)

	values = append(values, userId, reviewId)

	rows, err := r.db.Exec(query, values...)
	rowsCount, _ := rows.RowsAffected()
	if rowsCount == 0 {
		return errors.New("review not found")
	}
	return err
}

func (r *ReviewPostgres) DeleteReview(reviewId, userId uuid.UUID) error {
	query := fmt.Sprintln(
		`
		DELETE FROM reviews WHERE author_id = $1 AND id = $2 
		`,
	)
	rows, err := r.db.Exec(query, userId, reviewId)
	rowsCount, _ := rows.RowsAffected()
	if rowsCount == 0 {
		return errors.New("review not found")
	}
	return err
}

func (r *ReviewPostgres) GetAllReviewByF(filmTitle string) ([]models.ReviewOutput, error) {
	var output []models.ReviewOutput

	query := fmt.Sprintln(
		`
		SELECT 
		r.id,
		r.body,
		r.rating,
		(select username from users where id=r.author_id) as author_username
		FROM reviews r
		WHERE film_title = $1
		`,
	)
	err := r.db.Select(&output, query, filmTitle)
	return output, err
}

func (r *ReviewPostgres) GetAllReview() ([]models.Film, error) {

	query := fmt.Sprintln(
		`
		SELECT 
		f.title AS title,
		(SELECT sum(rating)/count(film_title)::real FROM reviews GROUP BY film_title) AS avg_rating,
		f.cover_path AS cover_path,
		r.body,
		r.rating,
		r.id,
		(SELECT username FROM users WHERE id = r.author_id) AS author_username
		FROM reviews r
		LEFT JOIN films f ON r.film_title = f.title
		`,
	)
	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Film{}, err
	}
	defer rows.Close()

	films := make(map[string]models.Film)

	for rows.Next() {
		var film models.Film
		var review models.ReviewOutput

		err := rows.Scan(
			&film.Title,
			&film.AvgRating,
			&film.CoverUrl,
			&review.Body,
			&review.Rating,
			&review.Id,
			&review.AuthorUsername,
		)

		if err != nil {
			return []models.Film{}, err
		}

		target, ok := films[film.Title]
		if ok {
			target.Reviews = append(target.Reviews, review)
			films[film.Title] = target

		} else {
			film.Reviews = append(film.Reviews, review)
			films[film.Title] = film
		}
	}

	var output []models.Film
	for _, v := range films {
		output = append(output, v)
		fmt.Println(v)
	}

	return output, err
}
