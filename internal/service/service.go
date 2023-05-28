package service

import (
	"context"
	"mime/multipart"

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
	Users
	Articles
	Roles
	Cards
	Materials
	Contests
	Disciplines
	Projects
	ProjectModules
	S3BucketService S3Bucket
}

type Authorization interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) error
	GetUserByID(ctx context.Context, userID string) (model.User, error)
	GenerateTokenPair(ctx context.Context, email, password string) (*TokenDetails, error)
	ParseToken(token string) (string, error)
	SaveRefreshToken(ctx context.Context, userID string, details *TokenDetails) error
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type Users interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) error
	UpdateUser(ctx context.Context, userID string, dto dto.UpdateUser) error
	DeleteUser(ctx context.Context, userID string) error
	UpdateUserImageURL(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (model.User, error)
	GetAllUsers(userContext context.Context) ([]model.User, error)
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
	GetByID(ctx context.Context, articleID string) (model.Article, error)
	Update(ctx context.Context, articleID, userID string, article dto.UpdateArticle) error
	Delete(ctx context.Context, userID, projectId string) error
	CommentByID(ctx context.Context, articleID, userID string, comment dto.CreateComment) error
	GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
	ReplyToCommentByArticleIDAndCommentID(
		ctx context.Context, articleID, userID, parentCommentID string, comment dto.ReplyToComment) error
	GetRepliesByArticleIDAndCommentID(
		ctx context.Context, articleID, commentID string) ([]model.Comment, error)
	UpvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error
	DownvoteCommentByArticleIDAndCommentID(userContext context.Context, articleID, commentID, userID string) error
	GetVotesByArticleIDAndCommentID(userContext context.Context, articleID, commentID string) (int, error)
}

type Materials interface {
	Create(ctx context.Context, userID string, dto dto.CreateMaterial) error
	GetAll(ctx context.Context) ([]model.Material, error)
	GetByID(ctx context.Context, materialID string) (model.Material, error)
	Update(ctx context.Context, materialID, userID string, material dto.UpdateMaterial) error
	Delete(ctx context.Context, userID, materialID string) error
	GetAllByUserID(ctx context.Context, userID string) ([]model.Material, error)
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
	GetInvitationByID(ctx context.Context, userID, cardID, invitationID string) (model.InvitationCard, error)
	GetInvitationsByCardID(ctx context.Context, userID, cardID string) ([]model.InvitationCard, error)
	GetResponsesByUserID(ctx context.Context, userID string) ([]model.InvitationCard, error)
	AcceptInvitation(ctx context.Context, userID, cardID, invitationID string, input dto.AnswerInvitation) error
	DeclineInvitation(ctx context.Context, userID, cardID, invitationID string, input dto.AnswerInvitation) error
}

type Contests interface {
	Create(ctx context.Context, contest dto.CreateContest) error
	Update(ctx context.Context, contestID string, contest dto.UpdateContest) error
	Delete(ctx context.Context, contestID string) error
	GetByID(ctx context.Context, contestID string) (model.Contest, error)
	GetAll(ctx context.Context) ([]model.Contest, error)
}

type Disciplines interface {
	Create(ctx context.Context, discipline dto.CreateDiscipline) error
	Update(ctx context.Context, disciplineID int, discipline dto.UpdateDiscipline) error
	Delete(ctx context.Context, disciplineID int) error
	GetAll(ctx context.Context) ([]model.Discipline, error)
	GetByID(ctx context.Context, disciplineID int) (model.Discipline, error)
}

type Projects interface {
	Create(ctx context.Context, disciplineID int, project dto.CreateProject) error
	Update(ctx context.Context, disciplineID, projectID int, project dto.UpdateProject) error
	Delete(ctx context.Context, disciplineID, projectID int) error
	GetAll(ctx context.Context) ([]model.Project, error)
	GetByID(ctx context.Context, disciplineID, projectID int) (model.Project, error)
	GetAllByDisciplineID(ctx context.Context, disciplineID int) ([]model.Project, error)
}

type ProjectModules interface {
	Create(ctx context.Context, projectID int, module dto.CreateProjectModule) error
	Update(ctx context.Context, projectID, moduleID int, module dto.UpdateProjectModule) error
	Delete(ctx context.Context, projectID, moduleID int) error
	GetAll(ctx context.Context, projectID int) ([]model.ProjectModule, error)
	GetByID(ctx context.Context, moduleID int) (model.ProjectModule, error)
}

type S3Bucket interface {
	UploadFile(ctx context.Context, key string, file *multipart.FileHeader) error
}

func NewService(repo *repository.Repository, r *redis.Client, c config.AuthConfig, bucketName string) *Service {
	service := &Service{
		Authorization:   NewAuthService(repo.Users, repo.Roles, r, c),
		S3BucketService: NewS3BucketService(repo.S3Bucket, bucketName),
		Users:           NewUsersService(repo.Users),
		Roles:           NewRolesService(repo.Roles, repo.Users),
		Cards:           NewCardsService(repo.Cards, repo.Users),
		Materials:       NewMaterialsService(repo.Materials, repo.Users),
		Disciplines:     NewDisciplinesService(repo.Disciplines, repo.Projects),
		Projects:        NewProjectsService(repo.Projects, repo.ProjectModules),
		ProjectModules:  NewProjectModulesService(repo.ProjectModules),
		Contests:        NewContestsService(repo.Contests),
	}

	service.Articles = NewArticlesService(repo.Articles, repo.Users, service.S3BucketService, repo.Transactional)

	return service
}
