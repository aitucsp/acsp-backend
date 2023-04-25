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

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
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
			auth.Post("/refresh", h.refreshToken)
			auth.Post("/logout", h.userIdentity, h.logout)
		}

		// Define user routes with authentication middleware (userIdentity) for all routes
		scholar := rest.Group("/scholar", h.userIdentity)
		{
			articles := scholar.Group("/articles", h.userIdentity)
			{
				articles.Post("/", h.createArticle)                       // create an article
				articles.Get("/", h.getAllArticles)                       // get all articles
				articles.Get("?userID=:userID", h.getAllArticlesByUserID) // get all articles of user
				articles.Get("/:id", h.getArticleByID)                    // get article by id
				articles.Put("/:id", h.updateArticle)                     // update an article
				articles.Delete("/:id", h.deleteArticle)                  // delete an article

				comments := articles.Group("/:id/comments", h.userIdentity)
				{
					comments.Post("/", h.commentArticle)                         // comment an article
					comments.Get("/", h.getCommentsByArticleID)                  // get all comments by article id
					comments.Post("/:commentID/replies", h.replyToCommentByID)   // reply to comment by id and article id
					comments.Get("/:commentID/replies", h.getRepliesByCommentID) // get all replies by comment id and article id
					comments.Post("/:commentID/upvote", h.upvoteCommentByID)     // like comment by id and article id
					comments.Post("/:commentID/downvote", h.downvoteCommentByID) // dislike comment by id and article id
					comments.Get("/:commentID/votes", h.getVotesByCommentID)     // get all upvotes by comment id and article id
				}

			}

			materials := scholar.Group("/materials", h.userIdentity)
			{
				materials.Post("/", h.createMaterial)                // create a material
				materials.Get("/", h.getAllMaterials)                // get all materials
				materials.Get("/:userID", h.getAllMaterialsByUserID) // get all materials of user
				materials.Get("/:id", h.getMaterialByID)             // get material by id
				materials.Put("/:id", h.updateMaterial)              // update a material
				materials.Delete("/:id", h.deleteMaterial)           // delete a material
			}
		}

		// Define code connection routes
		codeConnection := rest.Group("/code-connection", h.userIdentity)
		{
			applicants := codeConnection.Group("/applicants")
			{
				applicants.Get("/", h.getAllCards) // get all applicants list
			}

			cards := codeConnection.Group("/cards", h.userIdentity)
			{
				cards.Get("/", h.getAllCardsByUserID)              // get all cards of user
				cards.Post("/", h.createCard)                      // create a card
				cards.Get("/:id", h.getCardByID)                   // get an information of card by id
				cards.Put("/:id", h.updateCard)                    // update card information by id
				cards.Delete("/:id", h.deleteCard)                 // delete card by id
				cards.Get("/:id/invitations", h.getInvitations)    // get all invitations of a user
				cards.Post("/:id/invitations", h.createInvitation) // send an invitation to a user
			}

			// invitations := codeConnection.Group("/invitations")
			// {
			// 	// TODO: Accept and decline functionality
			// 	invitations.Post("/accept")  // Accept an invitation
			// 	invitations.Post("/decline") // Decline an invitation
			// }
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
