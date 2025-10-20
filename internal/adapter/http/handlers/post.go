package handlers

import (
	"errors"
	"fmt"
	"time"

	data_errors "github.com/IraIvanishak/news-app/internal/adapter/http/errors"
	"github.com/IraIvanishak/news-app/internal/adapter/http/requests/post"
	"github.com/IraIvanishak/news-app/internal/core/domain"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service ports.PostService
}

var _ ports.PostHandlers = (*PostHandler)(nil)

func NewPostHandler(service ports.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (h *PostHandler) RenderItem(c *fiber.Ctx) error {
	idStr := c.Params("id")

	post, err := h.service.Find(c.Context(), idStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
			"Error": "Post not found",
		})
	}

	return c.Render("posts/item", post)
}

func (h *PostHandler) RenderList(c *fiber.Ctx) error {
	filter := domain.NewListFilter(c.QueryInt("limit", 10), c.QueryInt("offset", 0), c.Query("q", ""))

	posts, err := h.service.List(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": fmt.Errorf("failed to retrieve posts: %w", err).Error(),
		})
	}

	totalCount, err := h.service.Count(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": fmt.Errorf("failed to retrieve post count: %w", err).Error(),
		})
	}

	pagination := domain.NewPagination(filter, totalCount)
	renderData := struct {
		Posts []*domain.Post
		domain.Pagination
		Query string
	}{
		Posts:      posts,
		Pagination: pagination,
		Query:      filter.Query,
	}
	return c.Render("posts/list", renderData)
}

func (h *PostHandler) RenderCreateForm(c *fiber.Ctx) error {
	return c.Render("posts/create", fiber.Map{
		"Title": "Create New Post",
	})
}

func (h *PostHandler) RenderEditForm(c *fiber.Ctx) error {
	idStr := c.Params("id")

	post, err := h.service.Find(c.Context(), idStr)
	if err != nil {
		if errors.Is(err, data_errors.ErrPostNotFound) {
			return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
				"Error": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": fmt.Errorf("failed to retrieve post: %w", err).Error(),
		})
	}

	return c.Render("posts/edit", fiber.Map{
		"ID":      post.ID,
		"Title":   post.Title,
		"Content": post.Content,
	})
}

func (h *PostHandler) Store(c *fiber.Ctx) error {
	req := post.NewUpsertPostRequest(c)
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Render("error", fiber.Map{
			"Error": err.Error(),
		})
	}

	postModel := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.Store(c.Context(), postModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": fmt.Errorf("failed to create post: %w", err).Error(),
		})
	}

	return h.RenderList(c)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")

	req := post.NewUpsertPostRequest(c)
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Render("error", fiber.Map{
			"Error": err.Error(),
		})
	}

	postModel := &domain.Post{
		ID:        idStr,
		Title:     req.Title,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	if err := h.service.Update(c.Context(), postModel); err != nil {
		if errors.Is(err, data_errors.ErrPostNotFound) {
			return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
				"Error": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": fmt.Errorf("failed to update post: %w", err).Error(),
		})
	}

	return h.RenderList(c)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")

	if err := h.service.Delete(c.Context(), idStr); err != nil {
		if errors.Is(err, data_errors.ErrPostNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("failed to delete post: %w", err).Error())
	}

	return c.SendString(fmt.Errorf("post deleted successfully").Error())
}
