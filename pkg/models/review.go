package models

import "github.com/google/uuid"

type Review struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FilmTitle string    `json:"film_title" db:"film_title"`
	Body      string    `json:"body" db:"body"`
	Rating    int       `json:"rating" db:"rating"`
}

type ReviewOutput struct {
	Review
	AuthorUsername string `json:"author_username" db:"author_username"`
}

type ReviewInputForHandler struct {
	FilmTitle string `json:"film_title" db:"film_title"`
	Body      string `json:"body" db:"body"`
	Rating    string `json:"rating" db:"rating"`
}

type ReviewInput struct {
	Body   *string `json:"body" db:"body"`
	Rating *int    `json:"rating" db:"rating"`
}
