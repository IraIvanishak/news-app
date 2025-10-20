package domain

import (
	"time"
)

type Post struct {
	ID        string    `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"createdat" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedat" json:"updatedAt"`
}
