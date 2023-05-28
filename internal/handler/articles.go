package handler

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/validation"
)

// @Summary Create an article
// @Security ApiKeyAuth
// @Tags articles
// @Description Method for creating an article for a user by id in the database by user id
// @ID create-article
// @Accept  json
// @Produce  json
// @Param request body dto.CreateArticle true "article information"
// @Param file formData file true "article image"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles [post]
func (h *Handler) createArticle(c *fiber.Ctx) error {
	log.Println("Creating an article... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	input := dto.CreateArticle{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	input.Image, err = c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Articles.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created article",
	})
}

// @Summary Get all articles
// @Security ApiKeyAuth
// @Tags articles
// @Description Get all articles in the database
// @ID get-all-articles
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles [get]
func (h *Handler) getAllArticles(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all articles... ")

	articles, err := h.services.Articles.GetAll(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(articles),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"articles": articles,
	})
}

// @Summary Get all articles of user
// @Security ApiKeyAuth
// @Tags articles
// @Description Get all articles of user by user id
// @ID get-all-articles-by-user-id
// @Accept  json
// @Produce  json
// @Param userID query string true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/user [get]
func (h *Handler) getAllArticlesByUserID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all articles... ")

	userId := c.Query("userID")
	if userId == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": "User id is required",
		})
	}

	articles, err := h.services.Articles.GetAllByUserID(c.UserContext(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(articles),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"articles": articles,
	})
}

// @Summary Get article by id and user id
// @Security ApiKeyAuth
// @Tags articles
// @Description Get articles by id and user id
// @ID get-article-by-id-and-user-id
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id} [get]
func (h *Handler) getArticleByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting article by id... ")

	articleID := c.Params("id", "")
	if articleID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	article, err := h.services.Articles.GetByID(c.UserContext(), articleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"user_id": c.GetRespHeader(userCtx, ""),
		"article": article,
	})
}

// @Summary Update an article by id
// @Security ApiKeyAuth
// @Tags articles
// @Description Update data of article
// @ID update-title-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param input body dto.UpdateArticle true "article information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id} [put]
func (h *Handler) updateArticle(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a project... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateArticle

	articleID := c.Params("id", "")
	if articleID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Articles.Update(c.UserContext(), articleID, userID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Article updated"})
}

// @Summary Delete an article by id
// @Security ApiKeyAuth
// @Tags articles
// @Description Delete an article from database
// @ID delete-article-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id} [delete]
func (h *Handler) deleteArticle(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting an article")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	id := c.Params("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	err = h.services.Articles.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Article deleted"})
}

// @Summary Comment an article
// @Security ApiKeyAuth
// @Tags articles
// @Description Method for commenting an article by id and user id
// @ID comment-an-article
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param input body dto.CreateComment true "comment information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments [post]
func (h *Handler) commentArticle(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Commenting an article")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	articleID := c.Params("id", "")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	var input dto.CreateComment
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Articles.CommentByID(c.UserContext(), articleID, userID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Article commented"})
}

// @Summary Get all comments by article id
// @Security ApiKeyAuth
// @Tags articles
// @Description Get all comments of an article
// @ID get-all-comments-by-article-id
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments [get]
func (h *Handler) getCommentsByArticleID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Get all comments by article id... ")

	articleID := c.Params("id")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	comments, err := h.services.Articles.GetCommentsByArticleID(c.UserContext(), articleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(comments),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"comments": comments,
	})
}

// @Summary Reply to comment by article id and comment id
// @Security ApiKeyAuth
// @Tags articles
// @Description Reply to comment of an article
// @ID reply-to-comment-by-article-id-and-comment-id
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param commentID path string true "comment id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments/{commentID}/replies [post]
func (h *Handler) replyToCommentByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Replying to comment... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	articleID := c.Params("id")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	parentCommentID := c.Params("commentID")
	if parentCommentID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	var input dto.ReplyToComment
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Articles.ReplyToCommentByArticleIDAndCommentID(c.UserContext(),
		articleID,
		userID,
		parentCommentID,
		input)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Replied to the comment successfully",
	})
}

// @Summary Get all comments by article id
// @Security ApiKeyAuth
// @Tags articles
// @Description Get all comments of an article
// @ID get-all-replies-of-comment
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param commentID path string true "comment id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments/{commentID}/replies [get]
func (h *Handler) getRepliesByCommentID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Get all replies by comment id... ")

	articleID := c.Params("id")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	commentID := c.Params("commentID")
	if commentID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	comments, err := h.services.Articles.GetRepliesByArticleIDAndCommentID(c.UserContext(), articleID, commentID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(comments),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"comments": comments,
	})
}

// @Summary Upvote comment
// @Security ApiKeyAuth
// @Tags articles
// @Description Upvote comment of an article by article id and comment id
// @ID upvote-comment
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param commentID path string true "comment id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments/{commentID}/upvote [post]
func (h *Handler) upvoteCommentByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Up-voting a comment... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	articleID := c.Params("id")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	commentID := c.Params("commentID")
	if commentID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	err = h.services.Articles.UpvoteCommentByArticleIDAndCommentID(c.UserContext(), articleID, commentID, userID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Up-voted the comment successfully",
	})
}

// @Summary Downvote comment
// @Security ApiKeyAuth
// @Tags articles
// @Description Downvote comment of an article by article id and comment id
// @ID downvote-comment
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param commentID path string true "comment id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments/{commentID}/downvote [post]
func (h *Handler) downvoteCommentByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Down-voting comment id... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	articleID := c.Params("id")
	if articleID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	commentID := c.Params("commentID")
	if commentID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	err = h.services.Articles.DownvoteCommentByArticleIDAndCommentID(c.UserContext(), articleID, commentID, userID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Did upvote the comment successfully",
	})
}

// @Summary Get vote difference of comment
// @Security ApiKeyAuth
// @Tags articles
// @Description Downvote comment of an article by article id and comment id
// @ID get-votes-of-comment
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param commentID path string true "comment id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/articles/{id}/comments/{commentID}/votes [get]
func (h *Handler) getVotesByCommentID(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Getting votes of comment... ")

	articleID := ctx.Params("id")
	if articleID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	commentID := ctx.Params("commentID")
	if commentID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	votes, err := h.services.Articles.GetVotesByArticleIDAndCommentID(ctx.UserContext(), articleID, commentID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"votes":   votes,
	})
}
