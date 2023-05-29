package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"acsp/internal/dto"
)

// @Summary Get user info by id
// @Security ApiKeyAuth
// @Tags users
// @Description get user info by id from the system
// @ID get-user
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/users/profile [get]
func (h *Handler) getUserProfile(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "user id is incorrect",
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
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

// @Summary Upload user image
// @Security ApiKeyAuth
// @Tags users
// @Description upload user image to the system
// @ID upload-user-image
// @Accept  json
// @Produce  json
// @Param file formData file true "user image"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/users/image [post]
func (h *Handler) uploadUserImage(c *fiber.Ctx) error {
	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.S3BucketService.UploadFile(c.UserContext(), userID, file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Users.UpdateUserImageURL(c.UserContext(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Successfully uploaded image",
	})
}

// @Summary Update user info by id
// @Security ApiKeyAuth
// @Tags users
// @Description Update user info by id from the system
// @ID update-user-information
// @Accept  json
// @Produce  json
// @Param input body dto.UpdateUser true "user information"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/users/profile [put]
func (h *Handler) updateUserProfile(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "user id is incorrect",
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateUser
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": "invalid input body",
		})
	}

	err = h.services.Users.UpdateUser(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Successfully updated user information",
	})
}
