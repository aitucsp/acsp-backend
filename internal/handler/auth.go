package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /auth/sign-up [post]
func (h *Handler) signUp(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Signing up... ")

	var input dto.CreateUser
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if validationErr := validate.Struct(&input); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": apperror.ErrBadInputBody,
		})
	}

	err := h.services.Authorization.CreateUser(ctx.UserContext(), input)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully registered"})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /auth/sign-in [post]
func (h *Handler) signIn(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Signing in...")

	var input signInInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": apperror.ErrBodyParsed,
		})
	}

	tokenPair, err := h.services.Authorization.GenerateTokenPair(ctx.UserContext(), input.Email, input.Password)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.services.Authorization.SaveTokenPair(ctx.UserContext(), tokenPair.UserID, tokenPair)

	return ctx.Status(http.StatusOK).JSON(tokenPair)
}

// @Summary RefreshToken
// @Tags auth
// @Description Refresh a token
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /auth/sign-in [post]
func (h *Handler) refreshToken(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(fiber.Map{})
}

func (h *Handler) logout(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{})
}
