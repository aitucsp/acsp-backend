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

// @Summary Create a lesson for a module of course
// @Security ApiKeyAuth
// @Tags courses
// @Description Method for creating a lesson for a module of a course
// @ID create-course-module-lesson
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Param request body dto.CreateCourseModuleLesson true "lesson information"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID}/lessons [post]
func (h *Handler) createCourseLesson(c *fiber.Ctx) error {
	log.Println("Creating a lesson for a module of a course... ")

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

	moduleID, err := c.ParamsInt("moduleID", -1)
	if moduleID == -1 {
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

	input := dto.CreateCourseModuleLesson{}
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

	err = h.services.ModuleLessons.Create(c.UserContext(), moduleID, input)
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

// @Summary Update a lesson of a module of a course by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Update data of a lesson of a module of a course by id
// @ID update-course-module-lesson-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Param lessonID path int true "lesson id"
// @Param request body dto.UpdateCourseModuleLesson true "lesson information"
// @Success 200 {object} string "Lesson updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID}/lessons/{lessonID} [put]
func (h *Handler) updateLesson(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a course module... ")

	var input dto.UpdateCourseModuleLesson

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

	moduleID, err := c.ParamsInt("moduleID", -1)
	if moduleID == -1 {
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

	err = h.services.ModuleLessons.Update(c.UserContext(), moduleID, lessonID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Module lesson is updated"})
}

// @Summary Delete a course module lesson by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Delete a course module by id
// @ID delete-course-module-lesson-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Param lessonID path int true "lesson id"
// @Success 200 {object} string "Course module deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID}/lessons/{lessonID} [delete]
func (h *Handler) deleteLesson(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a project")

	courseID, err := c.ParamsInt("id", -1)
	if courseID == -1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	moduleID, err := c.ParamsInt("id", -1)
	if moduleID == -1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	lessonID, err := c.ParamsInt("id", -1)
	if lessonID == -1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.ModuleLessons.Delete(c.UserContext(), moduleID, lessonID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Course module lesson is deleted"})
}

// @Summary Get a course module lesson by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Get a course module lesson by id
// @ID get-course-module-lesson-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Param lessonID path int true "lesson id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/course/{id}/modules/{moduleID}/lessons/{lessonID} [get]
func (h *Handler) getLessonByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting lesson of a module by id... ")

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

	moduleID, err := c.ParamsInt("moduleID", -1)
	if moduleID == -1 {
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

	lesson, err := h.services.ModuleLessons.GetByID(c.UserContext(), lessonID)
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
		"lesson":  lesson,
	})
}

// @Summary Get all lessons of a module
// @Security ApiKeyAuth
// @Tags courses
// @Description Get all course module lessons
// @ID get-all-course-module-lessons
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID}/lessons [get]
func (h *Handler) getAllLessonsByModuleID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all course modules... ")

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

	moduleID, err := c.ParamsInt("moduleID", -1)
	if moduleID == -1 {
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

	lessons, err := h.services.ModuleLessons.GetAll(c.UserContext(), moduleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(lessons),
		"lessons": lessons,
	})
}
