package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type CardsDatabase struct {
	db *sqlx.DB
}

func NewCardsRepository(db *sqlx.DB) *CardsDatabase {
	return &CardsDatabase{
		db: db,
	}
}

func (c *CardsDatabase) Create(ctx context.Context, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", card.UserID))

	query := fmt.Sprintf(`INSERT INTO %s (user_id, position, skills, description) 
								 VALUES ($1, $2, $3, $4) RETURNING id`,
		constants.CardsTable)

	_, err := c.db.Exec(query, card.UserID, card.Position, pq.Array(card.Skills), card.Description)
	if err != nil {
		l.Error("Error when creating the card in database", zap.Error(err))

		return err
	}

	return nil
}

func (c *CardsDatabase) Update(ctx context.Context, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", card.ID), zap.Int("userID", card.UserID))

	query := fmt.Sprintf(`UPDATE %s SET position = $1, skills = $2, description = $3, updated_at = now() 
								  WHERE id = $4 AND user_id = $5`,
		constants.CardsTable)

	_, err := c.db.Exec(query, card.Position, card.Skills, card.Description, card.ID, card.UserID)
	if err != nil {
		l.Error("Error when updating card in database", zap.Error(err))

		return err
	}

	return nil
}

func (c *CardsDatabase) Delete(ctx context.Context, userID int, cardID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND user_id = $2`, constants.CardsTable)

	_, err := c.db.Exec(query, cardID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardsDatabase) GetByID(ctx context.Context, cardID int) (*model.Card, error) {
	var card model.Card

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", constants.CardsTable)

	row := c.db.QueryRow(query, cardID)

	err := row.Scan(
		&card.ID,
		&card.UserID,
		&card.Position,
		&card.Skills,
		&card.Description,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (c *CardsDatabase) GetAllByUserID(ctx context.Context, userID int) (*[]model.Card, error) {
	var cards []model.Card

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", constants.CardsTable)

	err := c.db.Select(&cards, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Error when getting all cards by user id")
	}

	return &cards, nil
}

func (c *CardsDatabase) GetAll(ctx context.Context) ([]model.Card, error) {
	l := logging.LoggerFromContext(ctx)

	var cards []model.Card

	query := fmt.Sprintf(`SELECT c.*, u.id, u.email, u.name FROM %s c INNER JOIN %s u ON c.user_id = u.id`,
		constants.CardsTable, constants.UsersTable)

	rows, err := c.db.Query(query)
	if err != nil {
		l.Error("Error when querying get all applicants in database", zap.Error(err))

		return nil, err
	}

	for rows.Next() {
		var card model.Card

		err = rows.Scan(&card.ID,
			&card.UserID,
			&card.Position,
			&card.Skills,
			&card.Description,
			&card.CreatedAt,
			&card.UpdatedAt,
			&card.Author.ID,
			&card.Author.Email,
			&card.Author.Name)
		if err != nil {
			l.Error("Error when scanning the card in database", zap.Error(err))

			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (c *CardsDatabase) CreateInvitation(ctx context.Context, inviterID int, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("inviterID", inviterID))

	query := fmt.Sprintf(`INSERT INTO %s(card_id, inviter_id) VALUES ($1, $2) RETURNING id;`,
		constants.InvitationsTable)

	res, err := c.db.Exec(query, card.ID, inviterID)
	if err != nil {
		l.Error("Error when creating the card in database", zap.Error(err))

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Error when getting the number of rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(apperror.ErrRowsAffected, "Error when creating the invitation")
	}

	return nil
}

func (c *CardsDatabase) GetInvitationsByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error) {
	var invitationCards []model.InvitationCard

	query := fmt.Sprintf(`SELECT 
										c.user_id, 
										c.skills, 
										c.position, 
										c.description, 
										c.created_at, 
										c.updated_at,
										ci.inviter_id,
										ci.status
								 FROM %s c 
							     INNER JOIN %s ci ON ci.card_id = c.id
								 WHERE c.user_id = $1`,
		constants.CardsTable,
		constants.InvitationsTable)

	rows, err := c.db.Queryx(query, userID)

	for rows.Next() {
		var card model.Card
		var invitationCard model.InvitationCard

		err = rows.Scan(&card.UserID,
			&card.Skills,
			&card.Position,
			&card.Description,
			&card.CreatedAt,
			&card.UpdatedAt,
			&invitationCard.InviterID,
			&invitationCard.Status)
		if err != nil {
			return nil, errors.Wrap(err, "Error when scanning the card in database")
		}

		invitationCard.Card = &card
		invitationCards = append(invitationCards, invitationCard)
	}

	return invitationCards, nil
}

func (c *CardsDatabase) GetInvitationsByCardID(ctx context.Context, cardID int) ([]model.InvitationCard, error) {
	var invitationCards []model.InvitationCard

	query := fmt.Sprintf(`SELECT 
										c.user_id, 
										c.skills, 
										c.position, 
										c.description, 
										c.created_at, 
										c.updated_at,
										ci.inviter_id,
										ci.status
								 FROM %s c 
							     INNER JOIN %s ci ON ci.card_id = c.id
								 WHERE c.id = $1`,
		constants.CardsTable,
		constants.InvitationsTable)

	rows, err := c.db.Queryx(query, cardID)

	for rows.Next() {
		var card model.Card
		var invitationCard model.InvitationCard

		err = rows.Scan(&card.UserID,
			&card.Skills,
			&card.Position,
			&card.Description,
			&card.CreatedAt,
			&card.UpdatedAt,
			&invitationCard.InviterID,
			&invitationCard.Status)
		if err != nil {
			return nil, errors.Wrap(err, "Error when scanning the card in database")
		}

		invitationCard.Card = &card
		invitationCards = append(invitationCards, invitationCard)
	}

	return invitationCards, nil
}

func (c *CardsDatabase) GetInvitationByID(ctx context.Context, userID, cardID, invitationID int) (model.InvitationCard, error) {
	var card model.InvitationCard

	query := fmt.Sprintf(
		`SELECT 
						c.id,
						c.user_id, 
						c.position, 
						c.skills, 
						c.description, 
						c.created_at, 
						c.updated_at, 
						ci.inviter_id, 
						ci.status
					FROM %s c INNER JOIN %s ci ON c.id = ci.card_id 
						WHERE user_id = $1 AND c.id = $2 AND ci.id = $3`,
		constants.CardsTable,
		constants.InvitationsTable)

	row := c.db.QueryRow(query, userID, cardID, invitationID)

	err := row.Scan(
		&card.Card.ID,
		&card.Card.UserID,
		&card.Card.Position,
		&card.Card.Skills,
		&card.Card.Description,
		&card.Card.CreatedAt,
		&card.Card.UpdatedAt,
		&card.InviterID,
		&card.Status,
	)
	if err != nil {
		return model.InvitationCard{}, err
	}

	return card, nil
}
