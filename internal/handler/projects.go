package handler

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"gorestapi/internal/apperror"
	"gorestapi/internal/dto"
	"gorestapi/internal/logs"
	"gorestapi/internal/validation"
)

// @Summary Create a project
// @Security ApiKeyAuth
// @Tags projects
// @Description Creating a project
// @ID project-id
// @Accept  json
// @Produce  json
// @Param input body model.Project true "project information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects [post]
func (h *Handler) createProject(c *fiber.Ctx) error {
	log.Println("Creating a project... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
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
	if err := validate.Struct(input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	id, err := h.services.Projects.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":    false,
		"message":   nil,
		"projectId": id,
	})
}

// @Summary Get all projects
// @Security ApiKeyAuth
// @Tags projects
// @Description Get all projects from database
// @ID get-all-projects
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects [get]
func (h *Handler) getAllProjects(c *fiber.Ctx) error {
	logs.Log().Info("Getting all projects... ")

	_, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	projects, err := h.services.Projects.GetAll(c.UserContext())
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
		"user_id":  c.GetRespHeader(userCtx, ""),
		"projects": projects,
	})
}

// @Summary Get project by title
// @Security ApiKeyAuth
// @Tags projects
// @Description Get project by title from database
// @ID get-project-by-title
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Project
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:title [get]
func (h *Handler) getProjectByTitle(c *fiber.Ctx) error {
	logs.Log().Info("Getting a project by title... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	title := c.Params("title", "")
	if title == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	project, err := h.services.Projects.GetByTitle(c.UserContext(), userId, title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"project": project,
	})
}

// @Summary Update a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Update a title of project
// @ID update-project-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:id [put]
func (h *Handler) updateProject(c *fiber.Ctx) error {
	log.Println("Updating a project... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateProject
	projectId := c.Params("id", "")
	if projectId == "" {
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

	if err := h.services.Projects.Update(c.UserContext(), userId, projectId, *input.NewTitle); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

// @Summary Delete a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Delete a project from database
// @ID delete-project-by-title
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:id [delete]
func (h *Handler) deleteProject(c *fiber.Ctx) error {
	logs.Log().Info("Deleting a project")

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

	err = h.services.Projects.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{})
}
