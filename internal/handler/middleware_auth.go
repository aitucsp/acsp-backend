package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *fiber.Ctx) error {
	header := c.Get(authorizationHeader)

	if header == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "empty auth header",
		})
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid auth header",
		})
	}

	if len(headerParts[1]) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "token is empty",
		})
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "token parsing error",
		})
	}

	c.Set(userCtx, userId)
	return c.Next()
}

func getUserId(c *fiber.Ctx) (string, error) {
	id := c.GetRespHeader(userCtx, "")
	if len(id) == 0 {
		return "", apperror.ErrUserIDNotFound
	}

	return id, nil
}
