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

// @Summary Create a material
// @Security ApiKeyAuth
// @Tags materials
// @Description Method for creating a material for a user by id in the database by user id
// @ID create-material
// @Accept  json
// @Produce  json
// @Param input body dto.CreateMaterial true "material information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/materials [post]
func (h *Handler) createMaterial(c *fiber.Ctx) error {
	log.Println("Creating a material... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	input := dto.CreateMaterial{}
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

	err = h.services.Materials.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created material",
	})
}

// @Summary Get all materials by user id
// @Security ApiKeyAuth
// @Tags materials
// @Description Get all materials of user
// @ID get-all-materials-by-user-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/materials/:userID [get]
func (h *Handler) getAllMaterialsByUserID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all materials... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	materials, err := h.services.Materials.GetByUserID(c.UserContext(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(*materials),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"articles": materials,
	})
}

// GetAllByUserID @Summary Get all materials
// @Security ApiKeyAuth
// @Tags materials
// @Description Get all materials of all users
// @ID get-all-materials
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/materials [get]
func (h *Handler) getAllMaterials(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all materials... ")

	materials, err := h.services.Materials.GetAll(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(materials),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"articles": materials,
	})
}

// @Summary Get material by id and user id
// @Security ApiKeyAuth
// @Tags materials
// @Description Get materials by id and user id
// @ID get-material-by-id-and-user-id
// @Accept  json
// @Produce  json
// @Param id path string true "material id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/materials/:id [get]
func (h *Handler) getMaterialByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting material by id... ")

	materialID := c.Params("id", "")
	if materialID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	material, err := h.services.Materials.GetByID(c.UserContext(), materialID)
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
		"material": material,
	})
}

// @Summary Update a material by id
// @Security ApiKeyAuth
// @Tags materials
// @Description Update data of material by id
// @ID update-material-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "material id"
// @Param input body dto.UpdateMaterial true "material information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/materials/:id [put]
func (h *Handler) updateMaterial(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a material... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateMaterial

	materialID := c.Params("id", "")
	if materialID == "" {
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

	err = h.services.Materials.Update(c.UserContext(), materialID, userID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Material is updated"})
}

// @Summary Delete a material by id
// @Security ApiKeyAuth
// @Tags materials
// @Description Delete a material from database by id
// @ID delete-material-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "material id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/scholar/material/:id [delete]
func (h *Handler) deleteMaterial(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a material")

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

	err = h.services.Materials.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Material is deleted"})
}
