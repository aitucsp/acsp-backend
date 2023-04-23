package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"acsp/internal/model"
)

type Repository struct {
	Authorization
	Roles
	Articles
	Cards
	Materials
}

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, username, password string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
	ExistsUserByID(ctx context.Context, id int) (bool, error)
}

type Roles interface {
	CreateRole(ctx context.Context, name string) error
	UpdateRole(ctx context.Context, roleID int, newName string) error
	DeleteRole(ctx context.Context, roleID int) error
	SaveUserRole(ctx context.Context, userID, roleID int) error
	DeleteUserRole(ctx context.Context, userID, roleID int) error
	GetUserRoles(ctx context.Context, userID int) ([]model.Role, error)
}

type Articles interface {
	Create(ctx context.Context, article model.Article) error
	Update(ctx context.Context, article model.Article) error
	Delete(ctx context.Context, userID int, articleID int) error
	GetAllByUserID(ctx context.Context, userID int) (*[]model.Article, error)
	GetArticleByIDAndUserID(ctx context.Context, articleID, userID int) (*model.Article, error)
	CreateComment(ctx context.Context, articleID, userID int, comment model.Comment) error
	GetCommentsByArticleID(ctx context.Context, articleID int) ([]model.Comment, error)
	ReplyToComment(ctx context.Context, articleID, userID, parentCommentID int, comment model.Comment) error
	GetRepliesByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) ([]model.Comment, error)
	UpvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) error
	DownvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) error
	GetVotesByArticleIDAndCommentID(ctx context.Context, articleID, commentID int) (int, error)
	HasUserVotedForComment(ctx context.Context, userID, commentID int) (bool, error)
}

type Materials interface {
	Create(ctx context.Context, material model.Material) error
	Update(ctx context.Context, material model.Material) error
	Delete(ctx context.Context, userID int, materialID int) error
	GetByID(ctx context.Context, materialID int) (*model.Material, error)
	GetAllByUserID(ctx context.Context, userID int) (*[]model.Material, error)
	GetAll(ctx context.Context) ([]model.Material, error)
}

type Cards interface {
	Create(ctx context.Context, card model.Card) error
	Update(ctx context.Context, card model.Card) error
	Delete(ctx context.Context, userID int, cardID int) error
	GetByID(ctx context.Context, cardID int) (*model.Card, error)
	GetAllByUserID(ctx context.Context, userID int) (*[]model.Card, error)
	GetAll(ctx context.Context) ([]model.Card, error)
	CreateInvitation(ctx context.Context, inviterID int, card model.Card) error
	GetInvitationsByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error)
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Roles:         NewRolesRepository(db),
		Articles:      NewArticlesRepository(db),
		Cards:         NewCardsRepository(db),
		Materials:     NewMaterialsRepository(db),
	}
}
