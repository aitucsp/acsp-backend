package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"acsp/internal/config"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type ArticlesDatabase struct {
	db *sqlx.DB
}

func NewArticlesPostgres(db *sqlx.DB) *ArticlesDatabase {
	return &ArticlesDatabase{
		db: db,
	}
}

func (a *ArticlesDatabase) Create(ctx context.Context, article model.Article) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", article.ID))

	query := fmt.Sprint("INSERT INTO $1 (user_id, topic, description) VALUES ($2, $3, $4) RETURNING id")

	_, err := a.db.Exec(query, config.ArticlesTable, article.Author.ID, article.Topic, article.Description)
	if err != nil {
		l.Error("Error when getting the user from database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) Update(ctx context.Context, article model.Article) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("articleID", article.ID))

	query := fmt.Sprintf("UPDATE $1 SET topic = $2, description = $3, updated_at = now() WHERE id = $4 AND user_id = $5")

	_, err := a.db.Exec(query, config.ArticlesTable, article.Topic, article.Description, article.ID, article.Author.ID)
	if err != nil {
		l.Error("Error when update the article in database", zap.Error(err))

		return err
	}

	return nil
}

func (a *ArticlesDatabase) Delete(ctx context.Context, userID int, articleID int) error {
	query := fmt.Sprintf("DELETE FROM $1 WHERE id = $2 AND user_id = $3")

	_, err := a.db.Exec(query, config.ArticlesTable, articleID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesDatabase) GetAllByUserId(ctx context.Context, userID int) ([]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf("SELECT * FROM $1 WHERE user_id = $2")

	err := a.db.Select(&articles, query, config.RolesTable, userID)
	if err != nil {
		return []model.Article{}, err
	}

	return articles, nil
}
