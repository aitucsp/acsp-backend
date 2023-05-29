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

// @Summary Create a discipline
// @Security ApiKeyAuth
// @Tags projects
// @Description Method for creating a discipline
// @ID create-discipline
// @Accept  json
// @Produce  json
// @Param request body dto.CreateDiscipline true "discipline information"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines [post]
func (h *Handler) createDiscipline(c *fiber.Ctx) error {
	log.Println("Creating a discipline... ")

	input := dto.CreateDiscipline{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Disciplines.Create(c.UserContext(), input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created discipline",
	})
}

// @Summary Update a discipline by id
// @Security ApiKeyAuth
// @Tags disciplines
// @Description Update data of a discipline by id
// @ID update-discipline-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Param request body dto.UpdateDiscipline true "discipline information"
// @Success 200 {object} string "Discipline updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines/{id} [put]
func (h *Handler) updateDiscipline(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a discipline... ")

	var input dto.UpdateDiscipline

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

	err = h.services.Disciplines.Update(c.UserContext(), disciplineID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Discipline updated"})
}

// @Summary Delete a discipline by id
// @Security ApiKeyAuth
// @Tags disciplines
// @Description Delete a discipline by id
// @ID delete-discipline-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Success 200 {object} string "Discipline deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/discipline/{id} [delete]
func (h *Handler) deleteDiscipline(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a discipline")

	id, err := c.ParamsInt("id", -1)
	if id == -1 {
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

	err = h.services.Disciplines.Delete(c.UserContext(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Discipline deleted"})
}

// @Summary Get a discipline by id
// @Security ApiKeyAuth
// @Tags disciplines
// @Description Get a discipline by id
// @ID get-discipline-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "discipline id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/discipline/{id} [get]
func (h *Handler) getDisciplineByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting discipline by id... ")

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

	discipline, err := h.services.Disciplines.GetByID(c.UserContext(), disciplineID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":     false,
		"message":    nil,
		"user_id":    c.GetRespHeader(userCtx, ""),
		"discipline": discipline,
	})
}

// @Summary Get all disciplines
// @Security ApiKeyAuth
// @Tags disciplines
// @Description Get all disciplines
// @ID get-all-disciplines
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/disciplines [get]
func (h *Handler) getAllDisciplines(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all disciplines... ")

	disciplines, err := h.services.Disciplines.GetAll(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(disciplines),
		"projects": disciplines,
	})
}
