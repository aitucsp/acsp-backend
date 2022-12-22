package handler

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gorestapi/docs"
	"gorestapi/internal/service"

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
		projects := rest.Group("/projects", logger.New(), h.userIdentity)
		{
			projects.Post("/", h.createProject)
			projects.Get("/", h.getAllProjects)
			projects.Get("/:title", h.getProjectByTitle)
			projects.Put("/:id", h.updateProject)
			projects.Delete("/:id", h.deleteProject)
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func (h *Handler) InitPrometheusRoutes(app *fiber.App) *fiber.App {
	app.Get("/api/v1/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return app
}
