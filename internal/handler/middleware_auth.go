package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/logging"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())

	header := c.Get(authorizationHeader)

	if header == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "empty auth header",
		})
	}

	// Split the header into Bearer and the token part
	headerParts := strings.Split(header, " ")

	// Bearer token format
	// https://tools.ietf.org/html/rfc6750#section-2.1
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid auth header",
		})
	}

	// Check if the token is empty or not
	if len(headerParts[1]) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "token is empty",
		})
	}

	// Parse the token and get the user id
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		l.Error("Error when parsing token", zap.Error(err))

		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "token parsing error",
		})
	}

	// Set the user id in the context so that it can be used in other handlers to get the user information
	c.Set(userCtx, userId)

	// Call the next handler
	return c.Next()
}

// Authorize is a middleware that checks if the user has the required role(-s) to access the endpoint
func (h *Handler) Authorize(allowedRoles ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the user id from the context that was set in the userIdentity middleware
		userID, err := getUserId(c)
		if len(userID) == 0 || userID == "" {
			return errors.Wrap(err, apperror.ErrUserIDNotFound.Error())
		}
		fmt.Println(userID)
		// Get the user roles from the database
		roles, err := h.services.Roles.GetUserRoles(c.UserContext(), userID)
		if err != nil {
			return errors.Wrap(err, "Error when getting roles of user")
		}
		fmt.Println(roles)
		// Check if the user has the required role
		roleAllowed := false

		// Loop through the user roles
		for _, role := range roles {
			// Loop through the allowed roles
			for _, allowedRole := range allowedRoles {
				if role.Name == allowedRole {
					roleAllowed = true
					break
				}
			}
		}

		// If the user does not have the required role, return an error
		if !roleAllowed {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "You are not allowed to call this endpoint",
			})
		}

		// Call the next handler if the user has the required role
		return c.Next()
	}
}

// getUserId gets the user id from the context
func getUserId(c *fiber.Ctx) (string, error) {
	// Get the user id from the context that was set in the userIdentity middleware
	id := c.GetRespHeader(userCtx, "")
	if len(id) == 0 || id == "" {
		return "", apperror.ErrUserIDNotFound
	}

	return id, nil
}
