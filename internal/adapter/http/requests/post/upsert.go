package post

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UpsertPostRequest struct {
	Title   string `json:"title" form:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" form:"content" validate:"required,min=10"`
}

func (r *UpsertPostRequest) Validate(ctx *fiber.Ctx) error {
	validate := validator.New()

	if err := ctx.BodyParser(r); err != nil {
		return err
	}

	return validate.Struct(r)
}
