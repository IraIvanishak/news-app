package ports

import (
	"context"

	"github.com/IraIvanishak/news-app/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type PostHandlers interface {
	RenderItem(*fiber.Ctx) error
	RenderList(*fiber.Ctx) error
	RenderCreateForm(*fiber.Ctx) error
	RenderEditForm(*fiber.Ctx) error
	Store(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

type PostService interface {
	Store(context.Context, *domain.Post) error
	Find(context.Context, string) (*domain.Post, error)
	List(context.Context, domain.ListFilter) ([]*domain.Post, error)
	Count(context.Context) (int64, error)
	Update(context.Context, *domain.Post) error
	Delete(context.Context, string) error
}

type PostRepository interface {
	Store(context.Context, *domain.Post) error
	Find(context.Context, string) (*domain.Post, error)
	List(context.Context, domain.ListFilter) ([]*domain.Post, error)
	Count(context.Context) (int64, error)
	Update(context.Context, *domain.Post) error
	Delete(context.Context, string) error
}
