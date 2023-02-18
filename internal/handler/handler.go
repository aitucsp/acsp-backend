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

	// Define API routes
	rest := app.Group("/api/v1")
	{
		// Define auth routes
		auth := rest.Group("/auth")
		{
			auth.Post("/sign-up", h.signUp)
			auth.Post("/sign-in", h.signIn)
			// auth.Post("/refresh", h.refreshToken)
			auth.Post("/logout", h.userIdentity, h.logout)
		}

		scholar := rest.Group("/scholar", h.userIdentity)
		{
			articles := scholar.Group("/articles", h.userIdentity)
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

		codeConnection := rest.Group("/code-connection", h.userIdentity)
		{
			applicants := codeConnection.Group("/applicants")
			{
				applicants.Get("/", h.getAllCards)                 // get all applicants list
				applicants.Post("/:id/invite", h.createInvitation) // invite applicant by id
			}

			cards := codeConnection.Group("/cards", h.userIdentity)
			{
				cards.Get("/", h.getAllCardsByUserID)         // get all cards of user
				cards.Post("/", h.createCard)                 // create a card
				cards.Get("/:id", h.getCardByID)              // get an information of card by id
				cards.Put("/:id", h.updateCard)               // update card information by id
				cards.Delete("/:id", h.deleteCard)            // delete card by id
				cards.Post("/:id/invite", h.createInvitation) // send an invitation to a user
			}

			invitations := codeConnection.Group("/invitations")
			{
				invitations.Get("/", h.getInvitations) // Get all invitations of user

				// TODO Accept and decline functionality
				// invitations.Post("/accept")  // Accept an invitation
				// invitations.Post("/decline") // Decline an invitation
			}
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
