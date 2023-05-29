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
// @ID create-course-module
// @Accept  json
// @Produce  json
// @Param request body dto.CreateCourseModule true "module information"
// @Param id path int true "course id"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules [post]
func (h *Handler) createCourseModule(c *fiber.Ctx) error {
	log.Println("Creating a module for a course... ")

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

	input := dto.CreateCourseModule{}
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

	err = h.services.CourseModules.Create(c.UserContext(), courseID, input)
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

// @Summary Update a course module by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Update data of a course module
// @ID update-course-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Param request body dto.UpdateCourseModule true "module information"
// @Success 200 {object} string "Course updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID} [put]
func (h *Handler) updateCourseModule(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a course module... ")

	var input dto.UpdateCourseModule

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

	err = h.services.CourseModules.Update(c.UserContext(), courseID, moduleID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Course module is updated"})
}

// @Summary Delete a course module by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Delete a course module by id
// @ID delete-course-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Param moduleID path int true "module id"
// @Success 200 {object} string "Course module deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules/{moduleID} [delete]
func (h *Handler) deleteCourseModule(c *fiber.Ctx) error {
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

	err = h.services.CourseModules.Delete(c.UserContext(), courseID, moduleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Course module deleted"})
}

// @Summary Get a course module by id
// @Security ApiKeyAuth
// @Tags courses
// @Description Get a course module by id
// @ID get-course-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "course id"
// @Param moduleID path int true "module id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/course/{id}/modules/{moduleID} [get]
func (h *Handler) getCourseModuleByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting course module by id... ")

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

	module, err := h.services.CourseModules.GetByID(c.UserContext(), moduleID)
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
		"module":  module,
	})
}

// @Summary Get all course modules
// @Security ApiKeyAuth
// @Tags courses
// @Description Get all courses modules
// @ID get-all-course-modules
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/courses/{id}/modules [get]
func (h *Handler) getAllCourseModules(c *fiber.Ctx) error {
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

	modules, err := h.services.CourseModules.GetAll(c.UserContext(), courseID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(modules),
		"modules": modules,
	})
}
