package post

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UpsertPostRequest struct {
	Title   string `json:"title" form:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" form:"content" validate:"required,min=10"`

	// Optional photo metadata from Unsplash
	PhotoURL          string `json:"photoUrl" form:"photoUrl" validate:"omitempty,url"`
	PhotoAttribution  string `json:"photoAttribution" form:"photoAttribution"`
	PhotographerName  string `json:"photographerName" form:"photographerName"`
	PhotographerURL   string `json:"photographerUrl" form:"photographerUrl" validate:"omitempty,url"`
	UnsplashPhotoID   string `json:"unsplashPhotoId" form:"unsplashPhotoId"`
	DownloadLocation  string `json:"downloadLocation" form:"downloadLocation"`
}

func NewUpsertPostRequest(ctx *fiber.Ctx) *UpsertPostRequest {
	var req UpsertPostRequest
	if err := ctx.BodyParser(&req); err != nil {
		return nil
	}
	return &req
}

func (r *UpsertPostRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
