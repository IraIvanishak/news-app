package handlers

import (
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type UnsplashHandler struct {
	service ports.UnsplashService
}

// NewUnsplashHandler creates a new Unsplash handler
func NewUnsplashHandler(service ports.UnsplashService) *UnsplashHandler {
	return &UnsplashHandler{
		service: service,
	}
}

// SearchPhotos handles photo search requests
func (h *UnsplashHandler) SearchPhotos(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	photos, err := h.service.SearchPhotos(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Transform to response format
	photoResponses := make([]fiber.Map, 0, len(photos))
	for _, photo := range photos {
		photoResponses = append(photoResponses, fiber.Map{
			"id":               photo.ID,
			"thumbUrl":         photo.ThumbURL,
			"regularUrl":       photo.RegularURL,
			"description":      photo.Description,
			"photographer":     photo.PhotographerName,
			"photographerUrl":  photo.PhotographerURL,
			"downloadLocation": photo.DownloadLocation,
		})
	}

	return c.JSON(fiber.Map{
		"photos": photoResponses,
	})
}
