package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

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
		return nil, err
	}

	return &cards, nil
}

func (c *CardsDatabase) GetAll(ctx context.Context) (*[]model.Card, error) {
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

	return &cards, nil
}

func (c *CardsDatabase) CreateInvitation(ctx context.Context, inviterID int, card model.Card) error {
	l := logging.LoggerFromContext(ctx)

	tx, err := c.db.Begin()
	if err != nil {
		l.Error("Error when beginning the transaction", zap.Error(err))

		return err
	}

	query := fmt.Sprintf(`	INSERT INTO %s(card_id, inviter_id) VALUES ($1, $2) RETURNING id;`,
		constants.InvitationsTable)
	var id int

	err = tx.QueryRow(query, card.ID, inviterID).Scan(&id)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	if id < 1 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}

	querySecond := fmt.Sprintf(`INSERT INTO %s(invitation_id) VALUES ($1) RETURNING id;`,
		constants.InvitationResponsesTable)
	var responseID int

	err = tx.QueryRow(querySecond, id).Scan(&responseID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	return nil
}

func (c *CardsDatabase) GetInvitationsByUserID(ctx context.Context, userID int) (*[]model.InvitationCard, error) {
	var invitationCards []model.InvitationCard

	query := fmt.Sprintf(`SELECT user_id, skills, position, description, status, created_at, updated_at 
	FROM %s c
    INNER JOIN %s cr ON c.id = cr.card_id
    INNER JOIN %s ci on cr.id = ci.invitation_id
    WHERE inviter_id = $1`, constants.CardsTable, constants.InvitationsTable, constants.InvitationResponsesTable)

	rows, err := c.db.Queryx(query, userID)

	for rows.Next() {
		var card model.Card
		var invitationCard model.InvitationCard

		err = rows.Scan(&card.UserID,
			&card.Skills,
			&card.Position,
			&card.Description,
			&invitationCard.Status,
			&card.CreatedAt,
			&card.UpdatedAt)
		if err != nil {
			return nil, err
		}

		invitationCard.InviterID = card.UserID
		invitationCard.Card = &card
		invitationCards = append(invitationCards, invitationCard)
	}

	return &invitationCards, nil
}
