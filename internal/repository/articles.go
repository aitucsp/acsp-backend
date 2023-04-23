package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type ArticlesDatabase struct {
	db *sqlx.DB
}

func NewArticlesRepository(db *sqlx.DB) *ArticlesDatabase {
	return &ArticlesDatabase{
		db: db,
	}
}

func (a *ArticlesDatabase) Create(ctx context.Context, article model.Article) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", article.ID))

	query := fmt.Sprintf("INSERT INTO %s (user_id, topic, description) VALUES ($1, $2, $3) RETURNING id",
		constants.ArticlesTable)

	_, err := a.db.Exec(query, article.Author.ID, article.Topic, article.Description)
	if err != nil {
		l.Error("Error when creating the article in database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) Update(ctx context.Context, article model.Article) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", article.ID))

	query := fmt.Sprintf("UPDATE %s SET topic = $1, description = $2, updated_at = now() WHERE id = $3 AND user_id = $4",
		constants.ArticlesTable)

	_, err := a.db.Exec(query, article.Topic, article.Description, article.ID, article.Author.ID)
	if err != nil {
		l.Error("Error when update the article in database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) Delete(ctx context.Context, userID int, articleID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", constants.ArticlesTable)

	_, err := a.db.Exec(query, articleID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesDatabase) GetArticleByIDAndUserID(ctx context.Context, articleID, userID int) (*model.Article, error) {
	var article model.Article

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2",
		constants.ArticlesTable)

	err := a.db.Get(&article, query, articleID, userID)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (a *ArticlesDatabase) GetAllByUserID(ctx context.Context, userID int) (*[]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", constants.ArticlesTable)

	err := a.db.Select(&articles, query, userID)
	if err != nil {
		return nil, err
	}

	return &articles, nil
}

func (a *ArticlesDatabase) CreateComment(ctx context.Context, articleID, userID int, comment model.Comment) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", articleID), zap.Int("userID", userID))

	query := fmt.Sprintf("INSERT INTO %s (user_id, article_id, text) VALUES ($1, $2, $3) RETURNING id",
		constants.CommentsTable)

	_, err := a.db.Exec(query, userID, articleID, comment.Text)
	if err != nil {
		l.Error("Error when creating the comment in database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) GetCommentsByArticleID(ctx context.Context, articleID int) ([]model.Comment, error) {
	var comments []model.Comment
	query := fmt.Sprintf(`SELECT c.*, u.id as "user.user_id",
		       u.email AS "user.email",
		       u.name AS "user.name",
		       u.roles AS "user.roles" 
               FROM %s c INNER JOIN %s u ON u.id = c.user_id WHERE article_id = $1`,
		constants.CommentsTable,
		constants.UsersTable)

	rows, err := a.db.Queryx(query, articleID)

	for rows.Next() {
		var comment model.Comment
		var user model.User
		var parentID sql.NullInt64

		err = rows.Scan(&comment.ID,
			&comment.UserID,
			&comment.ArticleID,
			&parentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Upvotes,
			&comment.Downvotes,
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Roles)
		if err != nil {
			return []model.Comment{}, err
		}

		comment.Author = user
		comment.ParentID = int(parentID.Int64)
		comment.VoteDiff = comment.Upvotes - comment.Downvotes
		comments = append(comments, comment)
	}

	return comments, nil
}

func (a *ArticlesDatabase) ReplyToComment(ctx context.Context, articleID, userID, parentCommentID int, comment model.Comment) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("articleID", articleID),
		zap.Int("userID", userID),
		zap.Int("parentCommentID", parentCommentID),
	)

	query := fmt.Sprintf("INSERT INTO %s (user_id, article_id, parent_id, text) VALUES ($1, $2, $3, $4) RETURNING id",
		constants.CommentsTable)

	_, err := a.db.Exec(query, userID, articleID, parentCommentID, comment.Text)
	if err != nil {
		l.Error("Error when creating the comment in database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) GetRepliesByArticleIDAndCommentID(
	ctx context.Context, articleID,
	userID,
	parentCommentID int) ([]model.Comment, error) {

	var comments []model.Comment

	query := fmt.Sprintf(`SELECT c.*, u.id as "user.user_id",
		       u.email AS "user.email",
		       u.name AS "user.name",
		       u.roles AS "user.roles" 
               FROM %s c 
               INNER JOIN %s u ON u.id = c.user_id 
			   INNER JOIN %s a ON a.id = c.article_id
			   WHERE c.article_id = $1 AND c.parent_id = $2`,
		constants.CommentsTable,
		constants.UsersTable,
		constants.ArticlesTable)

	rows, err := a.db.Queryx(query, articleID, parentCommentID)

	for rows.Next() {
		var comment model.Comment
		var user model.User
		var parentID sql.NullInt64

		err = rows.Scan(&comment.ID,
			&comment.UserID,
			&comment.ArticleID,
			&parentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Roles)
		if err != nil {
			return nil, err
		}

		comment.Author = user
		comment.ParentID = int(parentID.Int64)
		comments = append(comments, comment)
	}

	return comments, nil
}

func (a *ArticlesDatabase) UpvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, commentID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("articleID", articleID),
		zap.Int("userID", userID),
		zap.Int("commentID", commentID),
	)

	tx, err := a.db.Begin()
	if err != nil {
		l.Error("Error when beginning the transaction", zap.Error(err))

		return err
	}

	query := fmt.Sprintf(`UPDATE %s SET upvotes = upvotes + 1 WHERE id = $1 AND article_id = $2 AND user_id = $3;`,
		constants.CommentsTable)
	var id int

	res, err := tx.Exec(query, commentID, articleID, userID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}

	querySecond := fmt.Sprintf("INSERT INTO %s (articleID, comment_id, user_id, vote_type) VALUES ($1, $2, $3, $4) RETURNING id",
		constants.CommentVotesTable)
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

	if responseID < 1 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
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

func (a *ArticlesDatabase) DownvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, commentID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("articleID", articleID),
		zap.Int("userID", userID),
		zap.Int("commentID", commentID),
	)

	tx, err := a.db.Begin()
	if err != nil {
		l.Error("Error when beginning the transaction", zap.Error(err))

		return err
	}

	query := fmt.Sprintf(`UPDATE %s SET upvotes = upvotes - 1 WHERE id = $1 AND article_id = $2 AND user_id = $3;`,
		constants.CommentsTable)
	var id int

	res, err := tx.Exec(query, commentID, articleID, userID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}

	querySecond := fmt.Sprintf("INSERT INTO %s (articleID, comment_id, user_id, vote_type) VALUES ($1, $2, $3, $4) RETURNING id",
		constants.CommentVotesTable)
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

	if responseID < 1 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
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

func (a *ArticlesDatabase) GetVotesByArticleIDAndCommentID(ctx context.Context, articleID, commentID int) (int, error) {
	var upvote, downvote int

	query := fmt.Sprintf("SELECT upvote, downvote FROM %s WHERE id = $1 AND article_id = $2",
		constants.CommentsTable)

	err := a.db.QueryRow(query, commentID, articleID).
		Scan(&upvote, &downvote)
	if err != nil {
		return 0, err
	}

	return upvote - downvote, nil
}

func (a *ArticlesDatabase) HasUserVotedForComment(ctx context.Context, userID, commentID int) (bool, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND comment_id = $2",
		constants.CommentVotesTable)

	err := a.db.QueryRow(query, userID, commentID).
		Scan(&id)
	if err != nil {
		return false, err
	}

	return id > 0, nil
}
