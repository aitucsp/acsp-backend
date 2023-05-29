package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"acsp/internal/apperror"
	"acsp/internal/logging"
)

// @Summary Get contest by id
// @Security ApiKeyAuth
// @Tags contests
// @Description Get contest by id
// @ID get-contest-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "contest id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/contests/{id} [get]
func (h *Handler) getContest(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting contest by id... ")

	contestID := c.Params("id", "")
	if contestID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": errors.Wrap(apperror.ErrInvalidParameter, "contest id is not correct"),
		})
	}

	contest, err := h.services.Contests.GetByID(c.UserContext(), contestID)
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
		"contest": contest,
	})
}

// @Summary Get all contests
// @Security ApiKeyAuth
// @Tags contests
// @Description Get all contests
// @ID get-all-contests
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/contests [get]
func (h *Handler) getAllContests(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all contests... ")

	contests, err := h.services.Contests.GetAll(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(contests),
		"cards":   contests,
	})
}
