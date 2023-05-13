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
