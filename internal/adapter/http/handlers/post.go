package handlers

import (
	"errors"
	"time"

	data_errors "github.com/IraIvanishak/news-app/internal/adapter/http/errors"
	"github.com/IraIvanishak/news-app/internal/adapter/http/requests/post"
	"github.com/IraIvanishak/news-app/internal/core/domain"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("error", fiber.Map{
			"Error": "Invalid post ID format",
		})
	}

	post, err := h.service.Find(c.Context(), objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
			"Error": "Post not found",
		})
	}

	return c.Render("posts/item", post)
}

func (h *PostHandler) RenderList(c *fiber.Ctx) error {
	filter := domain.ListFilter{
		Limit:  int64(c.QueryInt("limit", 10)),
		Offset: int64(c.QueryInt("offset", 0)),
	}

	posts, err := h.service.List(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": "Failed to retrieve posts",
		})
	}

	totalCount, err := h.service.Count(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": "Failed to retrieve post count",
		})
	}

	nextOffset := filter.Offset + filter.Limit
	prevOffset := filter.Offset - filter.Limit
	hasMore := nextOffset < totalCount

	return c.Render("posts/list", fiber.Map{
		"Posts":       posts,
		"HasPrevious": filter.Offset > 0,
		"HasNext":     hasMore,
		"NextOffset":  nextOffset,
		"PrevOffset":  prevOffset,
		"Limit":       filter.Limit,
		"TotalCount":  totalCount,
	})
}

func (h *PostHandler) RenderCreateForm(c *fiber.Ctx) error {
	return c.Render("posts/create", fiber.Map{
		"Title": "Create New Post",
	})
}

func (h *PostHandler) RenderEditForm(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("error", fiber.Map{
			"Error": "Invalid post ID format",
		})
	}

	post, err := h.service.Find(c.Context(), objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
			"Error": "Post not found",
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
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.service.Store(c.Context(), postModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
			"Error": "Failed to create post",
		})
	}

	return h.RenderList(c)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("error", fiber.Map{
			"Error": "Invalid post ID format",
		})
	}

	req := post.NewUpsertPostRequest(c)
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Render("posts/edit", fiber.Map{
			"ID":    idStr,
			"Error": err.Error(),
		})
	}

	postModel := &domain.Post{
		ID:        objectID,
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
		return c.Status(fiber.StatusInternalServerError).Render("posts/edit", fiber.Map{
			"ID":    idStr,
			"Error": "Failed to update post",
		})
	}

	return h.RenderList(c)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("error", fiber.Map{
			"Error": "Invalid post ID format",
		})
	}

	if err := h.service.Delete(c.Context(), objectID); err != nil {
		if errors.Is(err, data_errors.ErrPostNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete post")
	}

	return c.SendString("Post deleted successfully")
}
