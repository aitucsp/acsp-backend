package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

// Create creates a new card in the database.
func (c *CardsDatabase) Create(ctx context.Context, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", card.UserID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (user_id, position, skills, description) 
								 VALUES ($1, $2, $3, $4) RETURNING id`,
		constants.CardsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(card.UserID, card.Position, pq.Array(card.Skills), card.Description)
	if err != nil {
		l.Error("Error when executing the card creating statement", zap.Error(err))

		return err
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the id of the card", zap.Error(err))

		return err
	}

	if id == 0 {
		l.Error("Error when getting the id of the card", zap.Error(apperror.ErrCreatingCard))

		return apperror.ErrCreatingCard
	}

	return nil
}

// Update updates a card in the database.
func (c *CardsDatabase) Update(ctx context.Context, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", card.ID), zap.Int("userID", card.UserID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET position = $1, skills = $2, description = $3, updated_at = now() 
								  WHERE id = $4 AND user_id = $5`,
		constants.CardsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(card.Position, card.Skills, card.Description, card.ID, card.UserID)
	if err != nil {
		l.Error("Error when executing the card updating statement", zap.Error(err))

		return err
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of result", zap.Error(err))

		return err
	}

	if id == 0 {
		l.Error("Error when getting the id of the card", zap.Error(apperror.ErrUpdatingCard))

		return apperror.ErrUpdatingCard
	}

	return nil
}

// Delete deletes a card in the database.
func (c *CardsDatabase) Delete(ctx context.Context, userID int, cardID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", cardID), zap.Int("userID", userID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND user_id = $2`, constants.CardsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(cardID, userID)
	if err != nil {
		return err
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of the result", zap.Error(err))

		return err
	}

	if id == 0 {
		l.Error("Error when getting the id of the card", zap.Error(apperror.ErrDeletingCard))

		return apperror.ErrDeletingCard
	}

	return nil
}

// GetByID gets a card by id.
func (c *CardsDatabase) GetByID(ctx context.Context, cardID int) (*model.Card, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", cardID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var card model.Card

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", constants.CardsTable)

	err := c.db.QueryRowContext(ctx, query, cardID).Scan(
		&card.ID,
		&card.UserID,
		&card.Position,
		&card.Skills,
		&card.Description,
		&card.CreatedAt,
		&card.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		l.Error("Error when getting card by id", zap.Error(err))

		return &model.Card{}, errors.Wrap(err, "Error when getting card by id")
	case err != nil:
		l.Error("Error when getting card by id", zap.Error(err))

		return &model.Card{}, errors.Wrap(err, "Error when getting card by id")
	default:
		return &card, nil
	}
}

func (c *CardsDatabase) GetByIdAndUserID(ctx context.Context, userID, cardID int) (model.Card, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", cardID), zap.Int("userID", userID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var card model.Card

	query := fmt.Sprintf(`SELECT 
									c.id, 
									c.user_id, 
									c.position, 
									c.skills, 
									c.description, 
									c.created_at, 
									c.updated_at 
								FROM %s c WHERE id = $1 AND user_id = $2`,
		constants.CardsTable)

	err := c.db.QueryRowContext(ctx, query, cardID, userID).Scan(
		&card.ID,
		&card.UserID,
		&card.Position,
		&card.Skills,
		&card.Description,
		&card.CreatedAt,
		&card.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		l.Error("Error when getting card by id", zap.Error(err))

		return model.Card{}, err
	case err != nil:
		l.Error("Error when getting card by id", zap.Error(err))

		return model.Card{}, errors.Wrap(err, "Error when getting card by id")
	default:
		return card, nil
	}
}

// GetAllByUserID gets all cards by user id.
func (c *CardsDatabase) GetAllByUserID(ctx context.Context, userID int) (*[]model.Card, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("cardID", userID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var cards []model.Card

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", constants.CardsTable)

	err := c.db.SelectContext(ctx, &cards, query, userID)
	if err != nil {
		l.Error("Error when getting all cards by user id", zap.Error(err))

		return nil, errors.Wrap(err, "Error when getting all cards by user id")
	}

	return &cards, nil
}

// GetAll gets all cards.
func (c *CardsDatabase) GetAll(ctx context.Context) ([]model.Card, error) {
	l := logging.LoggerFromContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var cards []model.Card

	query := fmt.Sprintf(`SELECT c.*, u.id, u.email, u.name, u.image_url FROM %s c INNER JOIN %s u ON c.user_id = u.id`,
		constants.CardsTable, constants.UsersTable)

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		l.Error("Error when querying get all applicants in database", zap.Error(err))

		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			l.Error("Error when closing the rows", zap.Error(err))
		}
	}(rows)

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
			&card.Author.Name,
			&card.Author.ImageURL)
		if err != nil {
			l.Error("Error when scanning the card in database", zap.Error(err))

			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

// CreateInvitation creates an invitation in the database.
func (c *CardsDatabase) CreateInvitation(ctx context.Context, inviterID int, card model.Card) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("inviterID", inviterID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s(card_id, inviter_id) VALUES ($1, $2) RETURNING id;`,
		constants.CardInvitationsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(card.ID, inviterID)
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

// GetInvitationsByUserID gets all invitations by user id.
func (c *CardsDatabase) GetInvitationsByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID))

	var invitationCards []model.InvitationCard

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
		constants.CardInvitationsTable)

	rows, err := c.db.QueryContext(ctx, query, userID)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			l.Error("Error when closing the rows", zap.Error(err))

			logging.LoggerFromContext(ctx).Error("Error when closing the rows", zap.Error(err))
		}
	}(rows)

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
			l.Error("Error when scanning the card in database", zap.Error(err))

			return nil, errors.Wrap(err, "Error when scanning the card in database")
		}

		invitationCard.Card = &card
		invitationCards = append(invitationCards, invitationCard)
	}

	return invitationCards, nil
}

// GetInvitationsByCardID gets all invitations by card id.
func (c *CardsDatabase) GetInvitationsByCardID(ctx context.Context, cardID int) ([]model.InvitationCard, error) {
	var invitationCards []model.InvitationCard

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
		constants.CardInvitationsTable)

	rows, err := c.db.QueryContext(ctx, query, cardID)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.LoggerFromContext(ctx).Error("Error when closing the rows", zap.Error(err))
		}
	}(rows)

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

// GetInvitationByID gets an invitation by user id, card id and invitation id.
func (c *CardsDatabase) GetInvitationByID(ctx context.Context, userID, cardID, invitationID int) (model.InvitationCard, error) {
	var card model.InvitationCard

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
		constants.CardInvitationsTable)

	row := c.db.QueryRowContext(ctx, query, userID, cardID, invitationID)

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

// AcceptCardInvitation accepts a card invitation.
func (c *CardsDatabase) AcceptCardInvitation(ctx context.Context, userID, cardID, invitationID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("userID", userID),
		zap.Int("cardID", cardID),
		zap.Int("invitationID", invitationID),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s ci SET ci.status = $1, ci.updated_at = now()
								  FROM %s c
								  WHERE c.user_id = $2 AND c.id = $3 AND ci.id = $4`,
		constants.CardInvitationsTable, constants.CardsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(constants.AcceptedStatus, userID, cardID, invitationID)
	if err != nil {
		l.Error("Error when executing the card updating statement", zap.Error(err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of result", zap.Error(err))

		return err
	}

	if rows == 0 {
		l.Error("No rows affected", zap.Error(apperror.ErrAnsweringCard))

		return apperror.ErrAnsweringCard
	}

	return nil
}

// DeclineCardInvitation rejects a card invitation.
func (c *CardsDatabase) DeclineCardInvitation(ctx context.Context, userID, cardID, invitationID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("userID", userID),
		zap.Int("cardID", cardID),
		zap.Int("invitationID", invitationID),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s ci SET ci.status = $1, ci.updated_at = now()
								  FROM %s c
								  WHERE c.user_id = $2 AND c.id = $3 AND ci.id = $4`,
		constants.CardInvitationsTable, constants.CardsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(constants.DeclinedStatus, userID, cardID, invitationID)
	if err != nil {
		l.Error("Error when executing the card invitation declining query", zap.Error(err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of result", zap.Error(err))

		return err
	}

	if rows == 0 {
		l.Error("No rows affected", zap.Error(apperror.ErrAnsweringCard))

		return apperror.ErrAnsweringCard
	}

	return nil
}

func (c *CardsDatabase) GetResponsesByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error) {
	var invitationCards []model.InvitationCard

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
								 WHERE ci.inviter_id = $1`,
		constants.CardsTable,
		constants.CardInvitationsTable)

	rows, err := c.db.QueryContext(ctx, query, userID)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.LoggerFromContext(ctx).Error("Error when closing the rows", zap.Error(err))
		}
	}(rows)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logging.LoggerFromContext(ctx).Error("Error when closing the rows", zap.Error(err))
		}
	}(rows)

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
