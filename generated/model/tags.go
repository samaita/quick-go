package model

import "time"

// Tag maps to the tags table.
type Tag struct {
	Id int64 `json:"id" db:"id,pk"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
