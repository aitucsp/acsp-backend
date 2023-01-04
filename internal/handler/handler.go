package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"acsp/docs"

	"acsp/internal/service"

	_ "github.com/swaggo/files"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutesFiber(app *fiber.App) *fiber.App {
	docs.SwaggerInfo.BasePath = "/"

	// Define auth routes
	auth := app.Group("/auth")
	{
		auth.Post("sign-up", h.signUp)
		auth.Post("sign-in", h.signIn)
	}

	// Define API routes
	rest := app.Group("/api/v1")
	{
		articles := rest.Group("/articles", logger.New(), h.userIdentity)
		{
			articles.Post("/", h.createArticle)
			articles.Get("/", h.getAllArticles)
			articles.Put("/:id", h.updateArticle)
			articles.Delete("/:id", h.deleteArticle)
		}

	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
