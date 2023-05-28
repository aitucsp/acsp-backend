package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/constants"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type ArticlesService struct {
	repo      repository.Articles
	usersRepo repository.Users
	txManager repository.Transactional
	s3Bucket  S3Bucket
}

func NewArticlesService(r repository.Articles, a repository.Users, s S3Bucket, t repository.Transactional) *ArticlesService {
	return &ArticlesService{
		repo:      r,
		usersRepo: a,
		s3Bucket:  s,
		txManager: t,
	}
}

func (s *ArticlesService) Create(ctx context.Context, userID string, dto dto.CreateArticle) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("topic", dto.Topic))

	userId, err := strconv.Atoi(userID)
	if err != nil {
		l.Error("Error occurred when converting string to int", zap.Error(err))

		return errors.Wrap(err, "Error occurred when converting string to int")
	}

	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		l.Error("Error occurred when getting user by id", zap.Error(err))

		return err
	}

	article := model.Article{
		Topic:       dto.Topic,
		Description: dto.Description,
		Author:      &user,
	}

	tx, err := s.txManager.Begin(ctx, nil)
	if err != nil {
		l.Error("Error occurred when starting transaction", zap.Error(err))

		return err
	}

	defer func() {
		if err != nil {
			l.Error("Error occurred when rolling back transaction", zap.Error(err))

			err = s.txManager.Rollback(ctx, tx)
			return
		}

		err = s.txManager.Commit(ctx, tx)
		if err != nil {
			l.Error("Error occurred when committing transaction", zap.Error(err))

			return
		}
	}()

	err = s.repo.Create(ctx, tx, &article)
	if err != nil {
		l.Error("Error occurred when creating article", zap.Error(err))

		return errors.Wrap(err, "Error occurred when creating article")
	}

	if dto.Image != nil {
		err = s.s3Bucket.UploadFile(ctx, constants.ArticlesImagesFolder+"/"+strconv.Itoa(article.ID), dto.Image)
		if err != nil {
			l.Info("Error occurred when uploading file to s3 bucket", zap.Error(err))

			return err
		}

		err = s.repo.UpdateImageURL(ctx, tx, article.ID)
		if err != nil {
			l.Info("Error occurred when updating image URL", zap.Error(err))

			return err
		}
	}

	return nil
}

func (s *ArticlesService) GetAll(ctx context.Context) ([]model.Article, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return []model.Article{}, err
	}

	if articles == nil {
		return []model.Article{}, nil
	}

	articles = s.getFullURLForArticles(articles)

	return articles, nil
}

func (s *ArticlesService) GetAllByUserID(ctx context.Context, userID string) ([]model.Article, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("userID", userID))

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.Article{}, errors.Wrap(err, "failed to convert user id to int")
	}

	articles, err := s.repo.GetAllByUserID(ctx, userId)
	if err != nil {
		l.Error("Error occurred when getting all articles by user id", zap.Error(err))
		return []model.Article{}, err
	}

	if articles == nil {
		return []model.Article{}, nil
	}

	articles = s.getFullURLForArticles(articles)

	return articles, nil
}

func (s *ArticlesService) GetByID(ctx context.Context, articleID string) (model.Article, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("articleID", articleID))

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return model.Article{}, err
	}

	l.Info("Get article by id and user id")
	article, err := s.repo.GetArticleByID(ctx, articleId)
	if err != nil {
		l.Error("Error occurred when getting article by id and user id", zap.Error(err))

		return model.Article{}, err
	}

	l.Info("Get comments by article id")
	c, err := s.repo.GetCommentsByArticleID(ctx, articleId)
	if err != nil {
		l.Error("Error occurred when getting comments by article id", zap.Error(err))

		return model.Article{}, err
	}

	article.Comments = c
	article = s.getFullURLForArticle(article)

	return article, nil
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
		Author:      &user,
	}

	return s.repo.Update(ctx, article)
}

func (s *ArticlesService) Delete(ctx context.Context, userID string, projectId string) error {
	userId, _ := strconv.Atoi(userID)
	projectID, _ := strconv.Atoi(projectId)

	return s.repo.Delete(ctx, userId, projectID)
}

func (s *ArticlesService) CommentByID(ctx context.Context, articleID, userID string, commentDTO dto.CreateComment) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "failed to convert user id")
	}

	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return errors.Wrap(err, "failed to convert article id")
	}

	user, err := s.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "failed to get user by id")
	}

	comment := model.Comment{
		UserID:    userId,
		ArticleID: articleId,
		Text:      commentDTO.Text,
		Author:    user,
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
		return errors.Wrap(err, "failed to convert article id")
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "failed to convert user id")
	}

	parentCommentId, err := strconv.Atoi(parentCommentID)
	if err != nil {
		return errors.Wrap(err, "failed to convert parent comment id")
	}

	replyComment := model.Comment{
		UserID:    userId,
		ArticleID: articleId,
		ParentID:  parentCommentId,
		Text:      comment.Text,
	}

	return s.repo.ReplyToComment(ctx, articleId, userId, parentCommentId, replyComment)
}

func (s *ArticlesService) GetRepliesByArticleIDAndCommentID(ctx context.Context, articleID, commentID string) ([]model.Comment, error) {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return []model.Comment{}, errors.Wrap(err, "failed to convert article id")
	}

	parentCommentId, err := strconv.Atoi(commentID)
	if err != nil {
		return []model.Comment{}, errors.Wrap(err, "failed to convert parent comment id")
	}

	comments, err := s.repo.GetRepliesByArticleIDAndCommentID(ctx, articleId, parentCommentId)
	if err != nil {
		return []model.Comment{}, errors.Wrap(err, "failed to get replies by article id and comment id")
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
		return errors.Wrap(err, "failed to convert article id")
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "failed to convert user id")
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return errors.Wrap(err, "failed to convert comment id")
	}

	voted, err := s.repo.HasUserVotedForComment(userContext, userId, commentId)
	if err != nil {
		return errors.Wrap(err, "failed to check if user has voted for comment")
	}

	if voted {
		return errors.New("user has already voted for this comment")
	}

	return s.repo.UpvoteCommentByArticleIDAndCommentID(userContext, articleId, userId, commentId)
}

func (s *ArticlesService) DownvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return errors.Wrap(err, "failed to convert article id")
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "failed to convert user id")
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return errors.Wrap(err, "failed to convert comment id")
	}

	voted, err := s.repo.HasUserVotedForComment(userContext, userId, commentId)
	if err != nil {
		return errors.Wrap(err, "failed to check if user has voted for comment")
	}

	if voted {
		return errors.New("user has already voted for this comment")
	}

	return s.repo.DownvoteCommentByArticleIDAndCommentID(userContext, articleId, userId, commentId)
}

func (s *ArticlesService) GetVotesByArticleIDAndCommentID(userContext context.Context, articleID, commentID string) (int, error) {
	articleId, err := strconv.Atoi(articleID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert article id")
	}

	commentId, err := strconv.Atoi(commentID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert comment id")
	}

	return s.repo.GetVotesByArticleIDAndCommentID(userContext, articleId, commentId)
}

// getFullURLForArticles function gets a slice of articles and changes every article's image_url to a full url
func (s *ArticlesService) getFullURLForArticles(articles []model.Article) []model.Article {
	for i := range articles {
		articles[i].ImageURL = constants.BucketName + "." +
			constants.EndPoint + "/" +
			constants.ArticlesImagesFolder +
			articles[i].ImageURL
	}

	return articles
}

// getFullURLForArticle function gets an article and changes its image_url to a full url
func (s *ArticlesService) getFullURLForArticle(article model.Article) model.Article {
	article.ImageURL = constants.BucketName + "." +
		constants.EndPoint + "/" +
		constants.ArticlesImagesFolder +
		article.ImageURL

	return article
}
