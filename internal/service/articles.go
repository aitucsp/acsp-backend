package service

import (
	"context"
	"strconv"

	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type ArticlesService struct {
	repo      repository.Articles
	usersRepo repository.Authorization
}

func NewArticlesService(repo repository.Articles, usersRepo repository.Authorization) *ArticlesService {
	return &ArticlesService{repo: repo, usersRepo: usersRepo}
}

func (s *ArticlesService) Create(ctx context.Context, userID string, dto dto.CreateArticle) (int64, error) {
	userId, _ := strconv.Atoi(userID)
	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return -1, err
	}

	project := model.Article{
		Topic:       dto.Topic,
		Description: dto.Description,
		Author:      user,
	}

	return s.repo.Create(ctx, project)
}

func (s *ArticlesService) Update(ctx context.Context, userID string, articleDto dto.UpdateArticle) (int64, error) {
	userId, _ := strconv.Atoi(userID)
	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return -1, err
	}

	article := model.Article{
		ID:          articleDto.ID,
		Topic:       articleDto.Topic,
		Description: articleDto.Description,
		Author:      user,
	}

	return s.repo.Update(ctx, article)
}

func (s *ArticlesService) GetAll(ctx context.Context, userID string) ([]model.Article, error) {
	userId, _ := strconv.Atoi(userID)
	return s.repo.GetAllByUserId(ctx, userId)
}

func (s *ArticlesService) Delete(ctx context.Context, userID string, projectId string) (int64, error) {
	userId, _ := strconv.Atoi(userID)
	projectID, _ := strconv.Atoi(projectId)
	return s.repo.Delete(ctx, userId, projectID)
}
