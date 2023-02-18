package handler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/utils"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.User true "account info"
// @Success 200 {integer} string "message"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/auth/sign-up [post]
func (h *Handler) signUp(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Signing up... ")

	var input dto.CreateUser
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if validationErr := validate.Struct(&input); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": apperror.ErrBadInputBody,
		})
	}

	err := h.services.Authorization.CreateUser(ctx.UserContext(), input)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
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
// @Success 200 {string} service.TokenDetails "token"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/auth/sign-in [post]
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
			"error":   true,
			"message": err.Error(),
		})
	}

	err = h.services.Authorization.SaveRefreshToken(ctx.UserContext(), tokenPair.UserID, tokenPair)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(tokenPair)
}

// @Description Renew access and refresh tokens.
// @Summary renew access and refresh tokens
// @Tags auth
// @ID refresh-token
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func (h *Handler) refreshToken(ctx *fiber.Ctx) error {
	now := time.Now().Unix()

	// Get claims from JWT
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	expiresAccessToken := claims.Expires

	if now > expiresAccessToken {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, token has expired",
		})
	}

	var newToken string

	// Checking received data from JSON body.
	if err := ctx.BodyParser(newToken); err != nil {
		// Return, if JSON data is not correct.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(newToken)
	if err != nil {
		// Return status 400 and error message.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if now < expiresRefreshToken {
		// Define user ID.
		userID := claims.UserID

		// Get user by ID.
		foundedUser, err := h.services.Authorization.GetUserByID(ctx.UserContext(), userID)
		if err != nil {
			// Return, if user not found.
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "user with the given ID is not found",
			})
		}

		// Generate JWT Access & Refresh tokens.
		tokenPair, err := h.services.Authorization.GenerateTokenPair(ctx.UserContext(), foundedUser.Email, foundedUser.Password)
		if err != nil {
			// Return status 500 and token generation error.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Save refresh token to Redis.
		errRedis := h.services.SaveRefreshToken(ctx.UserContext(), userID, tokenPair)
		if errRedis != nil {
			// Return status 500 and Redis connection error.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   errRedis.Error(),
			})
		}

		return ctx.JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"tokens": fiber.Map{
				"access_token":  tokenPair.AccessToken,
				"refresh_token": tokenPair.RefreshToken,
			},
		})
	} else {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}
}

// @Description Logout user and delete refresh token from cache.
// @Summary Logout user
// @Tags auth
// @ID logout-user
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /api/v1/auth/logout [post]
func (h *Handler) logout(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Logging out... ")

	userId, err := getUserId(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	err = h.services.Authorization.DeleteRefreshToken(ctx.UserContext(), userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{})
}
