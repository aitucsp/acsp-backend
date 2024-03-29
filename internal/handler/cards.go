package handler

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/validation"
)

// @Summary Create a card
// @Security ApiKeyAuth
// @Tags cards
// @Description Method for creating a card
// @ID card-id
// @Accept  json
// @Produce  json
// @Param request body dto.CreateCard true "card information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards [post]
func (h *Handler) createCard(c *fiber.Ctx) error {
	log.Println("Creating a card... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	input := dto.CreateCard{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	err = h.services.Cards.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created card",
	})
}

// @Summary Update a card by id
// @Security ApiKeyAuth
// @Tags cards
// @Description Update data of card
// @ID update-card-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Param request body dto.UpdateCard true "card information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id} [put]
func (h *Handler) updateCard(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Updating a card... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateCard

	cardID, err := c.ParamsInt("id", -1)
	if cardID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
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

	err = h.services.Cards.Update(c.UserContext(), userID, cardID, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Card updated"})
}

// @Summary Delete a card by id
// @Security ApiKeyAuth
// @Tags cards
// @Description Delete a card from database
// @ID delete-card-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Card ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id} [delete]
func (h *Handler) deleteCard(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Deleting a card")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	id, err := c.ParamsInt("id", -1)
	if id == -1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Cards.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "Card deleted"})
}

// @Summary Get all cards of a user
// @Security ApiKeyAuth
// @Tags cards
// @Description Get all cards of a user
// @ID get-all-cards-by-user-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards [get]
func (h *Handler) getAllCardsByUserID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all cards... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cards, err := h.services.Cards.GetAllByUserID(c.UserContext(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(*cards),
		"user_id": c.GetRespHeader(userCtx, ""),
		"cards":   cards,
	})
}

// @Summary Get all cards
// @Security ApiKeyAuth
// @Tags cards
// @Description Get all cards
// @ID get-all-cards
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/applicants [get]
func (h *Handler) getAllCards(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all cards... ")

	applicants, err := h.services.Cards.GetAll(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(applicants),
		"cards":   applicants,
	})
}

// @Summary Get applicant card by id
// @Security ApiKeyAuth
// @Tags cards
// @Description Get applicant card by id
// @ID get-applicant-card-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "applicant id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/applicants/{id} [get]
func (h *Handler) getCardByID(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting card by id... ")

	cardID, err := c.ParamsInt("id", -1)
	if cardID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	card, err := h.services.Cards.GetByID(c.UserContext(), cardID)
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
		"card":    card,
	})
}

// @Summary Create an invitation
// @Security ApiKeyAuth
// @Tags cards
// @Description Invite a user by a card id
// @ID create-an-invitation
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Success 200 {integer} string "success"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id}/invitations [post]
func (h *Handler) createInvitation(c *fiber.Ctx) error {
	log.Println("Creating a card... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cardID, err := c.ParamsInt("id", -1)
	if cardID == -1 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Cards.CreateInvitation(c.UserContext(), userID, cardID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":  false,
		"message": "Created an invitation",
	})
}

// @Summary Get invitations of user
// @Security ApiKeyAuth
// @Tags cards
// @Description Get invitations of a user's cards
// @ID get-card-invitations
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/invitations [get]
func (h *Handler) getInvitations(c *fiber.Ctx) error {
	l := logging.LoggerFromContext(c.UserContext())
	l.Info("Getting all cards... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	invitations, err := h.services.Cards.GetInvitationsByUserID(c.UserContext(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(invitations),
		"cards":   invitations,
	})
}

// @Summary Get invitation by id and card id
// @Security ApiKeyAuth
// @Tags cards
// @Description Get invitation by id and card id
// @ID get-invitation-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Param invitationID path int true "invitation id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id}/invitations/{invitationID} [get]
func (h *Handler) getInvitationByID(c *fiber.Ctx) error {
	log.Println("Creating a card... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cardID := c.Params("id", "")
	if cardID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	invitationID := c.Params("invitationID", "")
	if invitationID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	invitation, err := h.services.Cards.GetInvitationByID(c.UserContext(), userID, cardID, invitationID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"card":    invitation,
	})
}

// @Summary Get all invitations by card
// @Security ApiKeyAuth
// @Tags cards
// @Description Get all invitations by card id
// @ID get-invitations-by-card-id
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id}/invitations [get]
func (h *Handler) getInvitationsByCardID(c *fiber.Ctx) error {
	log.Println("Getting invitation by card id... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cardID := c.Params("id", "")
	if cardID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	invitations, err := h.services.Cards.GetInvitationsByCardID(c.UserContext(), userID, cardID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(invitations),
		"cards":   invitations,
	})
}

// @Summary Accept invitation by invitation id and card id
// @Security ApiKeyAuth
// @Tags cards
// @Description Accept invitation by invitation id and card id
// @ID accept-invitation-by-id-and-card-id
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Param invitationID path int true "invitation id"
// @Param request body dto.AnswerInvitation true "accept invitation feedback"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id}/invitations/{invitationID}/accept [post]
func (h *Handler) acceptInvitation(c *fiber.Ctx) error {
	log.Println("Accepting an invitation... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cardID := c.Params("id", "")
	if cardID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	invitationID := c.Params("id", "")
	if invitationID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	var input dto.AnswerInvitation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Cards.AcceptInvitation(c.UserContext(), userID, cardID, invitationID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Successfully accepted invitation",
	})
}

// @Summary Decline invitation by invitation id and card id
// @Security ApiKeyAuth
// @Tags cards
// @Description Decline invitation by invitation id and card id
// @ID decline-invitation-by-id-and-card-id
// @Accept  json
// @Produce  json
// @Param id path int true "card id"
// @Param invitationID path int true "invitation id"
// @Param request body dto.AnswerInvitation true "decline invitation feedback"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/{id}/invitations/{invitationID}/decline [post]
func (h *Handler) declineInvitation(c *fiber.Ctx) error {
	log.Println("Creating a card... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	cardID := c.Params("id", "")
	if cardID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	invitationID := c.Params("id", "")
	if invitationID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	var input dto.AnswerInvitation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Cards.DeclineInvitation(c.UserContext(), userID, cardID, invitationID, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": "Successfully declined invitation",
	})
}

// @Summary Get responses of a user
// @Security ApiKeyAuth
// @Tags cards
// @Description Get all responses of a user
// @ID get-responses-by-user-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/code-connection/cards/responses [get]
func (h *Handler) getResponses(c *fiber.Ctx) error {
	log.Println("Getting invitation by card id... ")

	userID, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	responses, err := h.services.Cards.GetResponsesByUserID(c.UserContext(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":  false,
		"message": nil,
		"count":   len(responses),
		"cards":   responses,
	})
}
