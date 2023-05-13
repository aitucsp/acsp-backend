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

		// Define user routes
		users := rest.Group("/users", h.userIdentity)
		{
			users.Get("/profile", h.getUserProfile)
			users.Post("/image", h.uploadUserImage)
		}

		// Define user routes with authentication middleware (userIdentity) for all routes
		scholar := rest.Group("/scholar", h.userIdentity)
		{
			articles := scholar.Group("/articles")
			{
				articles.Post("/", h.createArticle)             // create an article
				articles.Get("/", h.getAllArticles)             // get all articles
				articles.Get("/user", h.getAllArticlesByUserID) // get all articles of user
				articles.Get("/:id", h.getArticleByID)          // get article by id
				articles.Put("/:id", h.updateArticle)           // update an article
				articles.Delete("/:id", h.deleteArticle)        // delete an article

				comments := articles.Group("/:id/comments")
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

			materials := scholar.Group("/materials")
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

			cards := codeConnection.Group("/cards")
			{
				cards.Get("/", h.getAllCardsByUserID)       // get all cards of user
				cards.Post("/", h.createCard)               // create a card
				cards.Get("/:id", h.getCardByID)            // get an information of card by id
				cards.Put("/:id", h.updateCard)             // update card information by id
				cards.Delete("/:id", h.deleteCard)          // delete card by id
				cards.Get("/invitations", h.getInvitations) // get all invitations of a user

				invitations := cards.Group("/:id/invitations")
				{
					invitations.Get("/", h.getInvitationsByCardID)                  // get all invitations of a card
					invitations.Post("/", h.createInvitation)                       // send an invitation to a user
					invitations.Get("/:invitationID", h.getInvitationByID)          // get an invitation by card id and invitation id
					invitations.Post("/:invitationID/accept", h.acceptInvitation)   // accept an invitation
					invitations.Post("/:invitationID/decline", h.declineInvitation) // decline an invitation
				}
			}
		}

		contests := rest.Group("/contests", h.userIdentity)
		{
			// contests.Post("/", h.Authorize("admin"), h.createContest) // create a contest
			contests.Get("/", h.getAllContests) // get all contests
			contests.Get("/:id", h.getContest)  // get contest by id
		}

		codingLab := rest.Group("/coding-lab", h.userIdentity)
		{
			projects := codingLab.Group("/projects")
			{
				projects.Post("/", h.Authorize("admin"), h.createProject)      // create a project
				projects.Put("/", h.Authorize("admin"), h.updateProject)       // update a project
				projects.Delete("/:id", h.Authorize("admin"), h.deleteProject) // delete a project
				projects.Get("/:id", h.getProjectByID)                         // get project by id
				projects.Get("/", h.getAllProjects)                            // get all projects

				modules := projects.Group("/:id/modules")
				{
					modules.Post("/", h.Authorize("admin"), h.createProjectModule)      // create a module
					modules.Put("/:id", h.Authorize("admin"), h.updateProjectModule)    // update a module
					modules.Delete("/:id", h.Authorize("admin"), h.deleteProjectModule) // delete a module
					modules.Get("/:id", h.getProjectModuleByID)                         // get module by id
					modules.Get("/", h.getAllProjectModules)                            // get all modules
				}
			}
		}

		// Define admin routes
		// admin := rest.Group("/admin", h.userIdentity, h.Authorize("admin"))
		// {
		// 	admin.Get("/users", h.getAllUsers)       // get all users
		// 	admin.Get("/users/:id", h.getUserByID)   // get user by id
		// 	admin.Put("/users/:id", h.updateUser)    // update user by id
		// 	admin.Delete("/users/:id", h.deleteUser) // delete user by id
		// 	admin.Post("/contests", h.createContest)
		// 	admin.Post("/contests/:id", h.updateContest)
		// 	admin.Delete("/contests/:id", h.deleteContest)
		// }
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
