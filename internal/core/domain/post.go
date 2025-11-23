package domain

import (
	"time"
)

type Post struct {
	ID        string    `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`

	// Photo metadata from Unsplash
	PhotoURL          string `bson:"photourl,omitempty" json:"photoUrl,omitempty"`
	PhotoAttribution  string `bson:"photoattribution,omitempty" json:"photoAttribution,omitempty"`
	PhotographerName  string `bson:"photographername,omitempty" json:"photographerName,omitempty"`
	PhotographerURL   string `bson:"photographerurl,omitempty" json:"photographerUrl,omitempty"`
	UnsplashPhotoID   string `bson:"unsplashphotoid,omitempty" json:"unsplashPhotoId,omitempty"`

	CreatedAt time.Time `bson:"createdat" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedat" json:"updatedAt"`
}
