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

// @Summary Create a module for a project
// @Security ApiKeyAuth
// @Tags projects
// @Description Method for creating a module of a project
// @ID create-project-module
// @Accept  json
// @Produce  json
// @Param request body dto.CreateProjectModule true "module information"
// @Param id path int true "project id"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/projects/{id}/modules [post]
func (h *Handler) createProjectModule(c *fiber.Ctx) error {
	log.Println("Creating a module for a project... ")

	// get project id from request
	projectID, err := c.ParamsInt("id", -1)
	if projectID == -1 {
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

	input := dto.CreateProjectModule{}
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

	err = h.services.ProjectModules.Create(c.UserContext(), projectID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created module for a project",
	})
}

// @Summary Update a project module by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Update data of a project module
// @ID update-project-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Param moduleID path int true "module id"
// @Param request body dto.UpdateProjectModule true "card information"
// @Success 200 {object} string "Project updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/projects/{id}/modules/{moduleID} [put]
func (h *Handler) updateProjectModule(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a project... ")

	var input dto.UpdateProjectModule

	projectID, err := c.ParamsInt("id", -1)
	if projectID == -1 {
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

	err = h.services.ProjectModules.Update(c.UserContext(), projectID, moduleID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Project module updated"})
}

// @Summary Delete a project module by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Delete a project module by id
// @ID delete-project-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Param moduleID path int true "module id"
// @Success 200 {object} string "Project module deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/projects/{id}/modules/{moduleID} [delete]
func (h *Handler) deleteProjectModule(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a project")

	projectID, err := c.ParamsInt("id", -1)
	if projectID == -1 {
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

	err = h.services.ProjectModules.Delete(c.UserContext(), projectID, moduleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Project module deleted"})
}

// @Summary Get a project module by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Get a project module by id
// @ID get-project-module-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Param moduleID path int true "module id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/projects/{id}/modules/{moduleID} [get]
func (h *Handler) getProjectModuleByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting project module by id... ")

	projectID, err := c.ParamsInt("id", -1)
	if projectID == -1 {
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

	module, err := h.services.ProjectModules.GetByID(c.UserContext(), moduleID)
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

// @Summary Get all project modules
// @Security ApiKeyAuth
// @Tags projects
// @Description Get all project modules
// @ID get-all-project-modules
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/projects/{id}/modules [get]
func (h *Handler) getAllProjectModules(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all project modules... ")

	projectID, err := c.ParamsInt("id", -1)
	if projectID == -1 {
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

	modules, err := h.services.ProjectModules.GetAll(c.UserContext(), projectID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(modules),
		"projects": modules,
	})
}
