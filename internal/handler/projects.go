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

// @Summary Create a project
// @Security ApiKeyAuth
// @Tags projects
// @Description Method for creating a project for a user by id in the database by user id
// @ID create-project
// @Accept  json
// @Produce  json
// @Param request body dto.CreateProject true "project information"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/:id/projects [post]
func (h *Handler) createProject(c *fiber.Ctx) error {
	log.Println("Creating a project... ")

	// get discipline id
	disciplineID, err := c.ParamsInt("id", -1)
	if disciplineID == -1 {
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

	input := dto.CreateProject{}
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

	err = h.services.Projects.Create(c.UserContext(), disciplineID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created project of discipline",
	})
}

// @Summary Update a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Update data of project
// @ID update-project-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Param projectID path int true "project id"
// @Param request body dto.UpdateProject true "project information"
// @Success 200 {object} string "Project updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/:id/projects/:projectID [put]
func (h *Handler) updateProject(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a project... ")

	var input dto.UpdateProject

	disciplineID, err := c.ParamsInt("id", -1)
	if disciplineID == -1 {
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

	projectID, err := c.ParamsInt("projectID", -1)
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

	err = h.services.Projects.Update(c.UserContext(), disciplineID, projectID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Project of discipline updated"})
}

// @Summary Delete a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Delete a project by id
// @ID delete-project-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Param projectID path int true "project id"
// @Success 200 {object} string "Project deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/:id/projects/:projectID [delete]
func (h *Handler) deleteProject(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a project")

	disciplineID, err := c.ParamsInt("id", -1)
	if disciplineID == -1 {
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

	projectID, err := c.ParamsInt("projectID", -1)
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

	err = h.services.Projects.Delete(c.UserContext(), disciplineID, projectID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Project deleted"})
}

// @Summary Get a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Get a project by id
// @ID get-project-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Param projectID path int true "project id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/:id/projects/:projectID [get]
func (h *Handler) getProjectByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting project by id... ")

	disciplineID, err := c.ParamsInt("id", -1)
	if disciplineID == -1 {
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

	projectID, err := c.ParamsInt("projectID", -1)
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

	project, err := h.services.Projects.GetByID(c.UserContext(), disciplineID, projectID)
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
		"project": project,
	})
}

// @Summary Get all projects by discipline id
// @Security ApiKeyAuth
// @Tags projects
// @Description Get all projects of discipline
// @ID get-all-projects-by-discipline-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/:id/projects [get]
func (h *Handler) getAllProjectsByDisciplineID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all projects by discipline... ")

	disciplineID, err := c.ParamsInt("id", -1)
	if disciplineID == -1 {
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

	projects, err := h.services.GetAllByDisciplineID(c.UserContext(), disciplineID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(projects),
		"projects": projects,
	})
}
