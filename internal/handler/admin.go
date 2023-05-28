package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/validation"
)

// @Summary Get all users
// @Security ApiKeyAuth
// @Tags admin
// @Description Get all users
// @ID get-all-users
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/users [get]
func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	users, err := h.services.Users.GetAllUsers(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"errors": false,
		"count":  len(users),
		"users":  users,
	})
}

// @Summary Get a user by id
// @Security ApiKeyAuth
// @Tags admin
// @Description Get a user by id
// @ID get-user-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/users/:id [get]
func (h *Handler) getUserByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting user by id... ")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	user, err := h.services.Users.GetUserByID(c.UserContext(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"errors": false,
		"user":   user,
	})
}

// @Summary Update a user by id
// @Security ApiKeyAuth
// @Tags admin
// @Description Update user by id
// @ID update-user-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "user id"
// @Param request body dto.UpdateUser true "user information"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/users/:id [put]
func (h *Handler) updateUser(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("updating a user...")

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "id is required",
		})
	}

	var input dto.UpdateUser
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	if err := h.services.Users.UpdateUser(c.UserContext(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"errors":  false,
		"message": "User updated",
	})
}

// @Summary Delete a user by id
// @Security ApiKeyAuth
// @Tags admin
// @Description Delete a user by id
// @ID delete-user-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} string "User deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/users/:id [delete]
func (h *Handler) deleteUser(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("deleting a user...")

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "id is required",
		})
	}

	err := h.services.Users.DeleteUser(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"errors":  false,
		"message": "User deleted",
	})
}

// @Summary Create a contest
// @Security ApiKeyAuth
// @Tags admin
// @Description Request for creating a contest
// @ID create-contest
// @Accept  json
// @Produce  json
// @Param request body dto.CreateContest true "contest information"
// @Success 200 {integer} map[string]interface{} "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/coding-lab/admin/contests [post]
func (h *Handler) createContest(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("creating a contest...")

	input := dto.CreateContest{}
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

	err = h.services.Contests.Create(c.UserContext(), input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"errors":  false,
		"message": "Contest created",
	})
}

// @Summary Update a contest by id
// @Security ApiKeyAuth
// @Tags admin
// @Description Update data of contest
// @ID update-contest-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "contest id"
// @Param request body dto.UpdateContest true "contest information"
// @Success 200 {object} map[string]interface{} "Contest updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/contests/:id [put]
func (h *Handler) updateContest(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a contest... ")

	var input dto.UpdateContest

	contestID := c.Params("id", "")
	if contestID == "" {
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

	err := h.services.Contests.Update(c.UserContext(), contestID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Contest updated"})
}

// @Summary Delete a contest by id
// @Security ApiKeyAuth
// @Tags admin
// @Description Delete a contest by id
// @ID delete-contest-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "contest id"
// @Success 200 {object} string "Contest deleted"
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/admin/contests/:id [delete]
func (h *Handler) deleteContest(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a contest... ")

	contestID := c.Params("id", "")
	if contestID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	err := h.services.Contests.Delete(c.UserContext(), contestID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{
		"errors":  false,
		"message": "Contest deleted",
	})
}
