package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

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

func (s *ArticlesService) Create(ctx context.Context, userID string, dto dto.CreateArticle) error {
	userId, _ := strconv.Atoi(userID)
	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	article := model.Article{
		Topic:       dto.Topic,
		Description: dto.Description,
		Author:      user,
	}

	return s.repo.Create(ctx, article)
}

func (s *ArticlesService) GetAll(ctx context.Context, userID string) (*[]model.Article, error) {
	userId, _ := strconv.Atoi(userID)
	return s.repo.GetAllByUserID(ctx, userId)
}

func (s *ArticlesService) GetByID(ctx context.Context, articleID, userID string) (*model.Article, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return &model.Article{}, err
	}

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return &model.Article{}, err
	}

	return s.repo.GetArticleByIDAndUserID(ctx, articleId, userId)
}

func (s *ArticlesService) Update(ctx context.Context, articleID string, userID string, articleDto dto.UpdateArticle) error {
	userId, _ := strconv.Atoi(userID)
	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	article := model.Article{
		ID:          articleId,
		Topic:       articleDto.Topic,
		Description: articleDto.Description,
		Author:      user,
	}

	return s.repo.Update(ctx, article)
}

func (s *ArticlesService) Delete(ctx context.Context, userID string, projectId string) error {
	userId, _ := strconv.Atoi(userID)
	projectID, _ := strconv.Atoi(projectId)

	return s.repo.Delete(ctx, userId, projectID)
}

func (s *ArticlesService) CommentByID(ctx context.Context, articleID, userID string, commentDTO dto.CreateComment) error {
	userId, _ := strconv.Atoi(userID)
	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	comment := model.Comment{
		UserID:    userId,
		ArticleID: articleId,
		Text:      commentDTO.Text,
		Author:    *user,
	}

	return s.repo.CreateComment(ctx, articleId, userId, comment)
}

func (s *ArticlesService) GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return []model.Comment{}, err
	}

	return s.repo.GetCommentsByArticleID(ctx, articleId)
}

func (s *ArticlesService) ReplyToCommentByArticleIDAndCommentID(ctx context.Context,
	articleID string,
	userID string,
	parentCommentID string,
	comment dto.ReplyToComment) error {

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	parentCommentId, err := strconv.Atoi(parentCommentID)
	if err != nil {
		return err
	}

	// convertedCommentID := sql.NullInt64{
	// 	Int64: int64(parentCommentId),
	// 	Valid: true,
	// }

	replyComment := model.Comment{
		UserID:    userId,
		ArticleID: articleId,
		ParentID:  parentCommentId,
		Text:      comment.Text,
	}

	return s.repo.ReplyToComment(ctx, articleId, userId, parentCommentId, replyComment)
}

func (s *ArticlesService) GetRepliesByArticleIDAndCommentID(ctx context.Context, articleID, userID, commentID string) ([]model.Comment, error) {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return []model.Comment{}, err
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.Comment{}, err
	}

	parentCommentId, err := strconv.Atoi(commentID)
	if err != nil {
		return []model.Comment{}, err
	}

	comments, err := s.repo.GetRepliesByArticleIDAndCommentID(ctx, articleId, userId, parentCommentId)

	if err != nil {
		return []model.Comment{}, err
	}

	// check if comments are empty and if empty, return empty slice
	if len(comments) == 0 {
		return []model.Comment{}, nil
	}

	return comments, nil
}

func (s *ArticlesService) UpvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return err
	}

	voted, err := s.repo.HasUserVotedForComment(userContext, userId, commentId)
	if err != nil {
		return err
	}

	if voted {
		return errors.New("user has already voted for this comment")
	}

	return s.repo.UpvoteCommentByArticleIDAndCommentID(userContext, articleId, userId, commentId)
}

func (s *ArticlesService) DownvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return err
	}

	voted, err := s.repo.HasUserVotedForComment(userContext, userId, commentId)
	if err != nil {
		return err
	}

	if voted {
		return errors.New("user has already voted for this comment")
	}

	return s.repo.DownvoteCommentByArticleIDAndCommentID(userContext, articleId, userId, commentId)
}

func (s *ArticlesService) GetVotesByArticleIDAndCommentID(userContext context.Context, articleID, commentID string) (int, error) {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return 0, err
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return 0, err
	}

	return s.repo.GetVotesByArticleIDAndCommentID(userContext, articleId, commentId)
}
