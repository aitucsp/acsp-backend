package handler

import (
	"github.com/gofiber/fiber/v2"
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
		auth.Post("/sign-up", h.signUp)
		auth.Post("/sign-in", h.signIn)
		// auth.Post("/refresh", h.refreshToken)
		// auth.Post("/logout", h.logout)
	}

	// Define API routes
	rest := app.Group("/api/v1")
	{
		articles := rest.Group("/articles", h.userIdentity)
		{
			articles.Post("/", h.createArticle)      // create an article
			articles.Get("/", h.getAllArticles)      // get all articles of user
			articles.Get("/:id", h.getArticleByID)   // get article by id
			articles.Put("/:id", h.updateArticle)    // update an article
			articles.Delete("/:id", h.deleteArticle) // delete an article

			comments := articles.Group("/:id/comments")
			{
				comments.Post("/", h.commentArticle)                        // comment an article
				comments.Get("/", h.getCommentsByArticleID)                 // get all comments by article id
				comments.Post("/:comment-id/replies", h.replyToCommentByID) // reply to comment by id and article id
				comments.Get("/:comment-id/replies", h.getRepliesByCommentID)
			}
		}

	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
