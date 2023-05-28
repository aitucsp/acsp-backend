package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/constants"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type CardsService struct {
	cardsRepo repository.Cards
	usersRepo repository.Users
}

func NewCardsService(cardsRepo repository.Cards, usersRepo repository.Users) *CardsService {
	return &CardsService{cardsRepo: cardsRepo, usersRepo: usersRepo}
}

func (c *CardsService) Create(ctx context.Context, userID string, dto dto.CreateCard) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	user, err := c.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "user not found in database")
	}

	card := model.Card{
		UserID:      userId,
		Position:    dto.Position,
		Skills:      dto.Skills,
		Description: dto.Description,
		Author:      user,
	}

	err = c.cardsRepo.Create(ctx, card)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardsService) Update(ctx context.Context, userID string, cardID int, dto dto.UpdateCard) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	user, err := c.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	card := model.Card{
		UserID:      userId,
		Position:    dto.Position,
		Skills:      dto.Skills,
		Description: dto.Description,
		Author:      user,
	}

	err = c.cardsRepo.Update(ctx, card)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardsService) Delete(ctx context.Context, userID string, cardID int) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	isUserExists, err := c.usersRepo.ExistsUserByID(ctx, id)
	if err != nil {
		return err
	}

	if isUserExists {
		err = c.cardsRepo.Delete(ctx, id, cardID)
		if err != nil {
			return errors.Wrap(err, "error when finding user in database")
		}
	} else {
		return errors.Wrap(apperror.ErrUserNotFound, "user not found in database")
	}

	return nil
}

func (c *CardsService) GetAllByUserID(ctx context.Context, userID string) (*[]model.Card, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error converting user id to int")
	}

	isUserExists, err := c.usersRepo.ExistsUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if isUserExists {
		cards, err := c.cardsRepo.GetAllByUserID(ctx, id)
		if err != nil {
			return nil, err
		}

		for _, card := range *cards {
			card.Author = c.getFullURLForUser(card.Author)
		}

		return cards, nil
	} else {
		return nil, errors.Wrap(apperror.ErrUserNotFound, "user not found in database")
	}
}

func (c *CardsService) GetAll(ctx context.Context) ([]model.Card, error) {
	var cards []model.Card

	cards, err := c.cardsRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		card.Author = c.getFullURLForUser(card.Author)
	}

	return cards, nil
}

func (c *CardsService) GetByID(ctx context.Context, cardID int) (*model.Card, error) {
	var card *model.Card

	card, err := c.cardsRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	card.Author = c.getFullURLForUser(card.Author)

	return card, nil
}

func (c *CardsService) CreateInvitation(ctx context.Context, userID string, cardID int) error {
	l := logging.LoggerFromContext(ctx)

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error when converting string to int")
	}

	userExists, err := c.usersRepo.ExistsUserByID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "error when finding user in database")
	}

	if userExists {
		card, err := c.cardsRepo.GetByID(ctx, cardID)
		if err != nil {
			l.Error("Error when getting card from database", zap.Error(err))

			return err
		}

		err = c.cardsRepo.CreateInvitation(ctx, userId, *card)
		if err != nil {
			l.Error("Error when creating an invitation in database", zap.Error(err))

			return err
		}
	} else {
		return errors.Wrap(apperror.ErrUserNotFound, "user not found in database")
	}

	return nil
}

func (c *CardsService) GetInvitationsByUserID(ctx context.Context, userID string) ([]model.InvitationCard, error) {
	userId, _ := strconv.Atoi(userID)
	_, err := c.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	cards, err := c.cardsRepo.GetInvitationsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *CardsService) GetInvitationsByCardID(ctx context.Context, userID, cardID string) ([]model.InvitationCard, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("userID", userID), zap.String("cardID", cardID))

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.InvitationCard{}, errors.Wrap(err, "error converting user id to int")
	}

	cardId, err := strconv.Atoi(cardID)
	if err != nil {
		return []model.InvitationCard{}, errors.Wrap(err, "error converting card id to int")
	}

	// get card by id and user id
	card, err := c.cardsRepo.GetByIdAndUserID(ctx, userId, cardId)
	if err != nil {
		l.Error("Error when getting card from database", zap.Error(err))

		return []model.InvitationCard{}, err
	}

	return c.cardsRepo.GetInvitationsByCardID(ctx, card.ID)
}

func (c *CardsService) GetInvitationByID(ctx context.Context, userID, cardID, invitationID string) (model.InvitationCard, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return model.InvitationCard{}, errors.Wrap(err, "error converting user id to int")
	}

	cardId, err := strconv.Atoi(cardID)
	if err != nil {
		return model.InvitationCard{}, errors.Wrap(err, "error converting card id to int")
	}

	invitationId, err := strconv.Atoi(invitationID)
	if err != nil {
		return model.InvitationCard{}, errors.Wrap(err, "error converting invitation id to int")
	}

	return c.cardsRepo.GetInvitationByID(ctx, userId, cardId, invitationId)
}

func (c *CardsService) AcceptInvitation(ctx context.Context, userID, cardID, invitationID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	cardId, err := strconv.Atoi(cardID)
	if err != nil {
		return errors.Wrap(err, "error converting card id to int")
	}

	invitationId, err := strconv.Atoi(invitationID)
	if err != nil {
		return errors.Wrap(err, "error converting invitation id to int")
	}

	return c.cardsRepo.AcceptCardInvitation(ctx, userId, cardId, invitationId)
}

func (c *CardsService) DeclineInvitation(ctx context.Context, userID, cardID, invitationID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	cardId, err := strconv.Atoi(cardID)
	if err != nil {
		return errors.Wrap(err, "error converting card id to int")
	}

	invitationId, err := strconv.Atoi(invitationID)
	if err != nil {
		return errors.Wrap(err, "error converting invitation id to int")
	}

	return c.cardsRepo.DeclineCardInvitation(ctx, userId, cardId, invitationId)
}

func (c *CardsService) GetResponsesByUserID(ctx context.Context, userID string) ([]model.InvitationCard, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.InvitationCard{}, errors.Wrap(err, "error converting user id to int")
	}

	return c.cardsRepo.GetResponsesByUserID(ctx, userId)
}

// getFullURLForUser function gets an article and changes its image_url to a full url
func (c *CardsService) getFullURLForUser(user model.User) model.User {
	user.ImageURL = constants.BucketName + "." +
		constants.EndPoint + "/" +
		constants.ArticlesImagesFolder +
		user.ImageURL

	return user
}
