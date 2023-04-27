package service

import (
	"context"
	"os"

	"github.com/go-redis/redis/v9"
	_ "github.com/golang/mock/gomock"

	"acsp/internal/config"
	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

// Mockgen is a tool for generating mock Go source files for interfaces.
//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	Authorization
	Articles
	Roles
	Cards
	Materials
	Contests
	S3Bucket
}

type Authorization interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) error
	GetUserByID(ctx context.Context, userID string) (model.User, error)
	GenerateTokenPair(ctx context.Context, email, password string) (*TokenDetails, error)
	ParseToken(token string) (string, error)
	SaveRefreshToken(ctx context.Context, userID string, details *TokenDetails) error
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type Roles interface {
	CreateRole(ctx context.Context, userID, name string) error
	UpdateRole(ctx context.Context, userID, roleID, newName string) error
	DeleteRole(ctx context.Context, userID, roleID string) error
	SaveUserRole(ctx context.Context, userID, roleID string) error
	DeleteUserRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]model.Role, error)
}

type Articles interface {
	Create(ctx context.Context, userID string, dto dto.CreateArticle) error
	GetAll(ctx context.Context) ([]model.Article, error)
	GetAllByUserID(ctx context.Context, userID string) ([]model.Article, error)
	GetByID(ctx context.Context, articleID, userID string) (*model.Article, error)
	Update(ctx context.Context, articleID, userID string, article dto.UpdateArticle) error
	Delete(ctx context.Context, userID, projectId string) error
	CommentByID(ctx context.Context, articleID, userID string, comment dto.CreateComment) error
	GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
	ReplyToCommentByArticleIDAndCommentID(
		ctx context.Context, articleID, userID, parentCommentID string, comment dto.ReplyToComment) error
	GetRepliesByArticleIDAndCommentID(
		ctx context.Context, articleID, userID, commentID string) ([]model.Comment, error)
	UpvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error
	DownvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error
	GetVotesByArticleIDAndCommentID(userContext context.Context, articleID, commentID string) (int, error)
}

type Materials interface {
	Create(ctx context.Context, userID string, dto dto.CreateMaterial) error
	GetAll(ctx context.Context) ([]model.Material, error)
	GetByID(ctx context.Context, materialID string) (*model.Material, error)
	Update(ctx context.Context, materialID, userID string, material dto.UpdateMaterial) error
	Delete(ctx context.Context, userID, materialID string) error
	GetByUserID(ctx context.Context, userID string) (*[]model.Material, error)
}

type Cards interface {
	Create(ctx context.Context, userID string, dto dto.CreateCard) error
	Update(ctx context.Context, userID string, cardID int, dto dto.UpdateCard) error
	Delete(ctx context.Context, userID string, cardID int) error
	GetAllByUserID(ctx context.Context, userID string) (*[]model.Card, error)
	GetAll(ctx context.Context) ([]model.Card, error)
	GetByID(ctx context.Context, cardID int) (*model.Card, error)
	CreateInvitation(ctx context.Context, userID string, cardID int) error
	GetInvitationsByUserID(ctx context.Context, userID string) ([]model.InvitationCard, error)
}

type Contests interface {
	GetByID(ctx context.Context, contestID string) (model.Contest, error)
	GetAll(ctx context.Context) ([]model.Contest, error)
}

type S3Bucket interface {
	UploadFile(ctx context.Context, bucket, key string, file *os.File) error
}

func NewService(repo *repository.Repository, r *redis.Client, c config.AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, repo.Roles, r, c),
		Articles:      NewArticlesService(repo.Articles, repo.Authorization),
		Roles:         NewRolesService(repo.Roles, repo.Authorization),
		Cards:         NewCardsService(repo.Cards, repo.Authorization),
		Materials:     NewMaterialsService(repo.Materials, repo.Authorization),
		Contests:      NewContestsService(repo.Contests),
		S3Bucket:      NewS3BucketService(repo.S3Bucket),
	}
}
