package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"acsp/internal/config"
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

func (a *ArticlesDatabase) Create(ctx context.Context, article model.Article) (int64, error) {
	query := fmt.Sprintf("INSERT INTO articles (user_id, topic, description) VALUES ($1, $2, $3) RETURNING id")

	result, err := a.db.Exec(query, article.Author.ID, article.Topic, article.Description)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

func (a *ArticlesDatabase) Update(ctx context.Context, article model.Article) (int64, error) {
	query := fmt.Sprintf("UPDATE $1 SET topic = $2, description = $3, updated_at = now() WHERE id = $4 AND user_id = $5")

	result, err := a.db.Exec(query, config.ArticlesTable, article.Topic, article.Description, article.ID, article.Author.ID)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

func (a *ArticlesDatabase) Delete(ctx context.Context, userID int, articleID int) (int64, error) {
	query := fmt.Sprintf("DELETE FROM $1 WHERE id = $2 AND user_id = $3")

	result, err := a.db.Exec(query, config.ArticlesTable, articleID, userID)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

func (a *ArticlesDatabase) GetAllByUserId(ctx context.Context, userID int) ([]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf("SELECT * FROM articles WHERE user_id = $1")

	err := a.db.Select(&articles, query, userID)
	if err != nil {
		return []model.Article{}, err
	}

	return articles, nil
}
