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
			users.Put("/profile", h.updateUserProfile)
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
				cards.Get("/responses", h.getResponses)     // get all cards of user
				cards.Get("/invitations", h.getInvitations) // get all invitations of a user
				cards.Get("/:id", h.getCardByID)            // get an information of card by id
				cards.Get("/", h.getAllCardsByUserID)       // get all cards of user
				cards.Post("/", h.createCard)               // create a card
				cards.Put("/:id", h.updateCard)             // update card information by id
				cards.Delete("/:id", h.deleteCard)          // delete card by id

				invitations := cards.Group("/:id/invitations")
				{
					invitations.Post("/", h.createInvitation)                       // send an invitation to a user
					invitations.Post("/:invitationID/accept", h.acceptInvitation)   // accept an invitation
					invitations.Get("/:invitationID", h.getInvitationByID)          // get an invitation by card id and invitation id
					invitations.Get("/", h.getInvitationsByCardID)                  // get all invitations of a card
					invitations.Post("/:invitationID/decline", h.declineInvitation) // decline an invitation
				}
			}
		}

		// Define contest routes
		contests := rest.Group("/contests", h.userIdentity)
		{
			// contests.Post("/", h.Authorize("admin"), h.createContest) // create a contest
			contests.Get("/", h.getAllContests) // get all contests
			contests.Get("/:id", h.getContest)  // get contest by id
		}

		// Define course routes
		courses := rest.Group("/courses", h.userIdentity)
		{
			courses.Post("/", h.Authorize("admin"), h.createCourse)      // create a course
			courses.Put("/:id", h.Authorize("admin"), h.updateCourse)    // update a course
			courses.Delete("/:id", h.Authorize("admin"), h.deleteCourse) // delete a course
			courses.Get("/", h.getAllCourses)                            // get all courses
			courses.Get("/:id", h.getCourseByID)                         // get course by id

			modules := courses.Group("/:id/modules")
			{
				modules.Post("/", h.Authorize("admin"), h.createCourseModule)     // create a module
				modules.Put("/:id", h.Authorize("admin"), h.updateCourseModule)   // update a module
				modules.Delete(":id", h.Authorize("admin"), h.deleteCourseModule) // delete a module
				modules.Get("/", h.getAllCourseModules)                           // get all modules
				modules.Get("/:id", h.getCourseModuleByID)                        // get module by id

				lessons := modules.Group("/:id/lessons")
				{
					lessons.Post("/", h.Authorize("admin"), h.createCourseLesson) // create a lesson
					lessons.Put("/:id", h.Authorize("admin"), h.updateLesson)     // update a lesson
					lessons.Delete("/:id", h.Authorize("admin"), h.deleteLesson)  // delete a lesson
					lessons.Get("/", h.getAllLessonsByModuleID)                   // get all lessons
					lessons.Get("/:id", h.getLessonByID)                          // get lesson by id
				}
			}

			lessonComments := courses.Group("/:id/lessons")
			{
				lessonComments.Post("/:id/comments", h.commentLesson)        // comment a lesson
				lessonComments.Get("/:id/comments", h.getLessonCommentsByID) // get all comments of a lesson
			}
		}

		codingLab := rest.Group("/coding-lab", h.userIdentity)
		{
			disciplines := codingLab.Group("/disciplines")
			{
				disciplines.Post("/", h.Authorize("admin"), h.createDiscipline)   // create a discipline
				disciplines.Put("/", h.Authorize("admin"), h.updateDiscipline)    // update a discipline
				disciplines.Delete("/", h.Authorize("admin"), h.deleteDiscipline) // delete a discipline
				disciplines.Get("/", h.getAllDisciplines)                         // get all disciplines
				disciplines.Get("/:id", h.getDisciplineByID)                      // get discipline by id

				projects := disciplines.Group("/:id/projects")
				{
					projects.Post("/", h.Authorize("admin"), h.createProject)             // create a project
					projects.Put("/", h.Authorize("admin"), h.updateProject)              // update a project
					projects.Delete("/:projectID", h.Authorize("admin"), h.deleteProject) // delete a project
					projects.Get("/:projectID", h.getProjectByID)                         // get project by id
					projects.Get("/", h.getAllProjectsByDisciplineID)                     // get all projects

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
		}

		// Define admin routes
		admin := rest.Group("/admin", h.userIdentity, h.Authorize("admin"))
		{
			admin.Get("/users", h.getAllUsers)       // get all users
			admin.Get("/users/:id", h.getUserByID)   // get user by id
			admin.Put("/users/:id", h.updateUser)    // update user by id
			admin.Delete("/users/:id", h.deleteUser) // delete user by id
			admin.Post("/contests", h.createContest)
			admin.Post("/contests/:id", h.updateContest)
			admin.Delete("/contests/:id", h.deleteContest)
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
