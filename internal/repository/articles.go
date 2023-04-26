package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
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

	res, err := a.db.Exec(query, article.Author.ID, article.Topic, article.Description)
	if err != nil {
		l.Error("error when creating the article in database", zap.Error(err))

		return errors.Wrap(err, "error when executing query")
	}

	r, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if r == 0 {
		return errors.Wrap(apperror.ErrCreatingArticle, "error when creating an article in database")
	}

	return nil
}

func (a *ArticlesDatabase) Update(ctx context.Context, article model.Article) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", article.ID))

	query := fmt.Sprintf("UPDATE %s SET topic = $1, description = $2, updated_at = now() WHERE id = $3 AND user_id = $4",
		constants.ArticlesTable)

	res, err := a.db.Exec(query, article.Topic, article.Description, article.ID, article.Author.ID)
	if err != nil {
		l.Error("Error when update the article in database", zap.Error(err))

		return errors.Wrap(err, "error when executing query")
	}

	r, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if r == 0 {
		return errors.Wrap(apperror.ErrUpdatingArticle, "error when update the article in database")
	}

	return nil
}

func (a *ArticlesDatabase) Delete(ctx context.Context, userID int, articleID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", constants.ArticlesTable)

	res, err := a.db.Exec(query, articleID, userID)
	if err != nil {
		return errors.Wrap(err, "error when executing query")
	}

	r, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if r == 0 {
		return errors.Wrap(apperror.ErrDeletingArticle, "error when delete the article in database")
	}

	return nil
}

func (a *ArticlesDatabase) GetAll(ctx context.Context) ([]model.Article, error) {
	var articles []model.Article

	query := fmt.Sprintf("SELECT * FROM %s", constants.ArticlesTable)

	err := a.db.Select(&articles, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing query")
	}

	return articles, nil
}

func (a *ArticlesDatabase) GetArticleByIDAndUserID(ctx context.Context, articleID, userID int) (*model.Article, error) {
	var article model.Article

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2",
		constants.ArticlesTable)

	err := a.db.Get(&article, query, articleID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "error when get article by id and user id")
	}

	return &article, nil
}

func (a *ArticlesDatabase) GetAllByUserID(ctx context.Context, userID int) ([]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", constants.ArticlesTable)

	err := a.db.Select(&articles, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing query")
	}

	return articles, nil
}

func (a *ArticlesDatabase) CreateComment(ctx context.Context, articleID, userID int, comment model.Comment) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", articleID), zap.Int("userID", userID))

	query := fmt.Sprintf("INSERT INTO %s (user_id, article_id, text) VALUES ($1, $2, $3) RETURNING id",
		constants.ArticlesCommentsTable)

	res, err := a.db.Exec(query, userID, articleID, comment.Text)
	if err != nil {
		l.Error("Error when creating the comment in database", zap.Error(err))

		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if r == 0 {
		return errors.Wrap(apperror.ErrCreatingComment, "error when creating the comment in database")
	}

	return nil
}

func (a *ArticlesDatabase) GetCommentsByArticleID(ctx context.Context, articleID int) ([]model.Comment, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", articleID))
	l.Info("Get comments by article id")

	var comments []model.Comment
	query := fmt.Sprintf(`SELECT c.*, u.id as "user.user_id",
		       u.email AS "user.email",
		       u.name AS "user.name",
		       u.roles AS "user.roles" 
               FROM %s c INNER JOIN %s u ON u.id = c.user_id WHERE article_id = $1`,
		constants.ArticlesCommentsTable,
		constants.UsersTable)

	rows, err := a.db.Queryx(query, articleID)
	if err != nil {
		return []model.Comment{}, errors.Wrap(err, "error when executing query")
	}

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
			return []model.Comment{}, errors.Wrap(err, "error when scanning row")
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
		constants.ArticlesCommentsTable)

	res, err := a.db.Exec(query, userID, articleID, parentCommentID, comment.Text)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		return errors.Wrap(err, "error occurred when executing the query")
	}

	r, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if r == 0 {
		return errors.Wrap(apperror.ErrCreatingComment, "error when replying to the comment in database")
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
		constants.ArticlesCommentsTable,
		constants.UsersTable,
		constants.ArticlesTable)

	rows, err := a.db.Queryx(query, articleID, parentCommentID)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing query")
	}

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
			return nil, errors.Wrap(err, "error when scanning row")
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

		return errors.Wrap(err, "error occurred when beginning the transaction")
	}

	query := fmt.Sprintf(`UPDATE %s SET upvotes = upvotes + 1 WHERE id = $1 AND article_id = $2 AND user_id = $3;`,
		constants.ArticlesCommentsTable)

	res, err := tx.Exec(query, commentID, articleID, userID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when getting rows affected")
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(apperror.ErrUpvoteComment, "error occurred when up-voting the comment in database")
	}

	querySecond := fmt.Sprintf(`INSERT INTO %s (articleID, comment_id, user_id, vote_type) 
														VALUES ($1, $2, $3, $4) 
														RETURNING id`,
		constants.ArticleCommentVotesTable)

	res, err = tx.Exec(querySecond, articleID, commentID, userID, constants.UpvoteType)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when executing the query")
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when getting rows affected")
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(apperror.ErrUpvoteComment, "error occurred when inserting upvote of the comment in database")
	}

	err = tx.Commit()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "Error when committing the transaction")
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

		return errors.Wrap(err, "error occurred when beginning the transaction")
	}

	query := fmt.Sprintf(`UPDATE %s SET upvotes = upvotes - 1 WHERE id = $1 AND article_id = $2 AND user_id = $3;`,
		constants.ArticlesCommentsTable)

	res, err := tx.Exec(query, commentID, articleID, userID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when getting rows affected")
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(apperror.ErrDownvoteComment, "error occurred when down-voting the comment in database")
	}

	querySecond := fmt.Sprintf(`INSERT INTO %s (articleID, comment_id, user_id, vote_type) 
														VALUES ($1, $2, $3, $4) 
														RETURNING id`,
		constants.ArticleCommentVotesTable)

	res, err = tx.Exec(querySecond, articleID, commentID, userID, constants.DownvoteType)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when executing the query")
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		l.Error("Error when getting rows affected", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "error occurred when getting rows affected")
	}

	if rowsAffected < 1 {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(apperror.ErrUpvoteComment, "error occurred when inserting downvote of the comment in database")
	}

	err = tx.Commit()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		return errors.Wrap(err, "Error when committing the transaction")
	}

	return nil
}

func (a *ArticlesDatabase) GetVotesByArticleIDAndCommentID(ctx context.Context, articleID, commentID int) (int, error) {
	var upvote, downvote int

	query := fmt.Sprintf("SELECT upvote, downvote FROM %s WHERE id = $1 AND article_id = $2",
		constants.ArticlesCommentsTable)

	err := a.db.QueryRow(query, commentID, articleID).Scan(&upvote, &downvote)
	if err != nil {
		return 0, errors.Wrap(err, "error occurred when executing the query")
	}

	return upvote - downvote, nil
}

func (a *ArticlesDatabase) HasUserVotedForComment(ctx context.Context, userID, commentID int) (bool, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND comment_id = $2",
		constants.ArticleCommentVotesTable)

	err := a.db.QueryRow(query, userID, commentID).Scan(&id)
	if err != nil {
		return false, errors.Wrap(err, "error occurred when executing the query")
	}

	return id > 0, nil
}
