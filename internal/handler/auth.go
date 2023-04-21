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
// @Description signing up a new user to the system
// @ID sign-up
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
// @Description signing in a user to the system and returning a token pair (access and refresh) to the client for further requests
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} service.TokenDetails "token pair"
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

	// Validate input data. If data is not valid, return error. Otherwise, continue.
	validate := validator.New()
	if validationErr := validate.Struct(&input); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": apperror.ErrBadInputBody,
		})
	}

	// Generate token pair
	tokenPair, err := h.services.Authorization.GenerateTokenPair(ctx.UserContext(), input.Email, input.Password)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Save refresh token to database
	err = h.services.Authorization.SaveRefreshToken(ctx.UserContext(), tokenPair.UserID, tokenPair)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Return token pair
	return ctx.Status(http.StatusOK).JSON(tokenPair)
}

// @Summary renew access and refresh tokens
// @Description Renew access and refresh tokens by refresh token from request body and set new expiration time for refresh token in database
// @Tags auth
// @ID refresh-token
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {string} service.TokenDetails "token pair"
// @Failure 400,401,404 {object} map[string]interface{} "error" "message"
// @Failure 500 {object} map[string]interface{} "error" "message"
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func (h *Handler) refreshToken(ctx *fiber.Ctx) error {
	now := time.Now().Unix()

	// Get claims from JWT
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Check if token is expired
	expiresAccessToken := claims.Expires

	// If token is expired, return status 401 and error message.
	if now > expiresAccessToken {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized, token has expired",
		})
	}

	var newToken string

	// Get new token from request body.
	if err := ctx.BodyParser(newToken); err != nil {
		// Return, if JSON data is not correct.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user
	// and check if it is valid.
	expiresRefreshToken, err := utils.ParseRefreshToken(newToken)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Check if refresh token is expired
	if now < expiresRefreshToken {
		// Define user ID
		userID := claims.UserID

		// Get user by ID
		foundedUser, err := h.services.Authorization.GetUserByID(ctx.UserContext(), userID)
		if err != nil {
			// Return, if user not found.
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user with the given ID is not found",
			})
		}

		// Generate new token pair for user with new expiration time
		tokenPair, err := h.services.Authorization.GenerateTokenPair(ctx.UserContext(), foundedUser.Email, foundedUser.Password)
		if err != nil {
			// Return status 500 and token generation error.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		// Save refresh token to Redis database with new expiration time
		// and check if it is valid. If not, return error.
		err = h.services.SaveRefreshToken(ctx.UserContext(), userID, tokenPair)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		// Return new token pair
		return ctx.JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"tokens": fiber.Map{
				"access_token":  tokenPair.AccessToken,
				"refresh_token": tokenPair.RefreshToken,
			},
		})
	} else {
		// Return status 401 and error message, if refresh token is expired.
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized, your session was ended earlier",
		})
	}
}

// @Summary logout user
// @Description Logout user and delete refresh token from cache by user ID
// @Tags auth
// @ID logout
// @Accept json
// @Produce json
// @Success 200 {string} status "ok"
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /api/v1/auth/logout [post]
func (h *Handler) logout(ctx *fiber.Ctx) error {
	l := logging.LoggerFromContext(ctx.UserContext())
	l.Info("Logging out... ")

	// Get user ID from context
	userId, err := getUserId(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Delete refresh token from cache by user ID
	err = h.services.Authorization.DeleteRefreshToken(ctx.UserContext(), userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 and message "ok" if everything is ok.
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "ok"})
}
