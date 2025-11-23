package routes

import (
	"github.com/IraIvanishak/news-app/internal/adapter/http/handlers"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

func PostRoutes(app *fiber.App, handler ports.PostHandlers, unsplashHandler *handlers.UnsplashHandler) {
	// Post routes
	app.Get("/posts", handler.RenderList)
	app.Get("/posts/create", handler.RenderCreateForm)
	app.Post("/posts", handler.Store)
	app.Get("/posts/:id", handler.RenderItem)
	app.Get("/posts/:id/edit", handler.RenderEditForm)
	app.Put("/posts/:id", handler.Update)
	app.Delete("/posts/:id", handler.Delete)

	// Unsplash API routes
	app.Get("/api/photos/search", unsplashHandler.SearchPhotos)
}
