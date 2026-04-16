package model

import "time"

// Post maps to the posts table.
type Post struct {
	Id int64 `json:"id" db:"id,pk"`
	Title string `json:"title" db:"title"`
	Body string `json:"body" db:"body"`
	AuthorId int64 `json:"author_id" db:"author_id"`
	Published bool `json:"published" db:"published"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
