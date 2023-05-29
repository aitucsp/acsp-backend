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

// @Summary Create a module for a course
// @Security ApiKeyAuth
// @Tags courses
// @Description Method for creating a module of a course
// @ID create-lesson-comment
// @Accept  json
// @Produce  json
// @Param request body dto.CreateLessonComment true "comment information"
// @Param id path int true "course id"
// @Param lessonID path int true "lesson id"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/lessons/{lessonID}/comments [post]
func (h *Handler) commentLesson(c *fiber.Ctx) error {
	log.Println("Commenting a lesson... ")

	// get project id from request
	courseID, err := c.ParamsInt("id", -1)
	if courseID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	lessonID, err := c.ParamsInt("lessonID", -1)
	if lessonID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	input := dto.CreateLessonComment{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
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

	err = h.services.LessonComments.Create(c.UserContext(), lessonID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created module for a course",
	})
}

// @Summary Get a course module by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Get a course module by id
// @ID get-lesson-comments-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param lessonID path int true "lesson id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/lessons/{lessonID}/comments [get]
func (h *Handler) getLessonCommentsByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting lesson comments by id... ")

	courseID, err := c.ParamsInt("id", -1)
	if courseID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	lessonID, err := c.ParamsInt("lessonID", -1)
	if lessonID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	comments, err := h.services.LessonComments.GetAll(c.UserContext(), lessonID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"user_id":  c.GetRespHeader(userCtx, ""),
		"comments": comments,
	})
}
