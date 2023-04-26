package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get user info by id
// @Tags users
// @Description get user info by id from the system
// @ID get-user
// @Accept  json
// @Produce  json
// @Param id path int true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/users/:id [get]
func (h *Handler) getUserInfo(c *fiber.Ctx) error {
	userId := c.Params("id", "")
	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "user id is incorrect",
		})
	}

	user, err := h.services.Authorization.GetUserByID(c.UserContext(), userId)
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
		"user":    user,
	})
}
