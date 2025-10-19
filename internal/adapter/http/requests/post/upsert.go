package post

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UpsertPostRequest struct {
	Title   string `json:"title" form:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" form:"content" validate:"required,min=10"`
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
