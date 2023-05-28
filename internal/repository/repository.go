package repository

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jmoiron/sqlx"

	"acsp/internal/dto"
	"acsp/internal/model"
)

type Repository struct {
	Users
	Roles
	Articles
	Cards
	Materials
	Contests
	Disciplines
	Projects
	ProjectModules
	Transactional
	S3Bucket
}

// Users interface provides methods for working with users
type Users interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, username, password string) (*model.User, error)
	GetByID(ctx context.Context, id int) (model.User, error)
	GetUserDetailsByUserId(ctx context.Context, id int) (*model.UserDetails, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	UpdateDetails(ctx context.Context, userID int, userDetails model.UserDetails) error
	ExistsUserByID(ctx context.Context, id int) (bool, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	UpdateImageURL(ctx context.Context, userID int) error
}

// Roles interface provides methods for working with roles
type Roles interface {
	CreateRole(ctx context.Context, name string) error
	UpdateRole(ctx context.Context, roleID int, newName string) error
	DeleteRole(ctx context.Context, roleID int) error
	SaveUserRole(ctx context.Context, userID, roleID int) error
	DeleteUserRole(ctx context.Context, userID, roleID int) error
	GetUserRoles(ctx context.Context, userID int) ([]model.Role, error)
}

type Articles interface {
	Create(ctx context.Context, tx *sqlx.Tx, article *model.Article) error
	Update(ctx context.Context, article model.Article) error
	UpdateImageURL(ctx context.Context, tx *sqlx.Tx, articleID int) error
	Delete(ctx context.Context, userID int, articleID int) error
	GetAll(ctx context.Context) ([]model.Article, error)
	GetAllByUserID(ctx context.Context, userID int) ([]model.Article, error)
	GetArticleByID(ctx context.Context, articleID int) (model.Article, error)
	GetArticleByIDAndUserID(ctx context.Context, articleID, userID int) (model.Article, error)
	CreateComment(ctx context.Context, articleID, userID int, comment model.Comment) error
	GetCommentsByArticleID(ctx context.Context, articleID int) ([]model.Comment, error)
	ReplyToComment(ctx context.Context, articleID, userID, parentCommentID int, comment model.Comment) error
	GetRepliesByArticleIDAndCommentID(ctx context.Context, articleID, parentCommentID int) ([]model.Comment, error)
	UpvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) error
	DownvoteCommentByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) error
	GetVotesByArticleIDAndCommentID(ctx context.Context, articleID, commentID int) (int, error)
	HasUserVotedForComment(ctx context.Context, userID, commentID int) (bool, error)
}

// Materials interface provides methods for working with materials.
type Materials interface {
	Create(ctx context.Context, material model.Material) error
	Update(ctx context.Context, material model.Material) error
	Delete(ctx context.Context, userID int, materialID int) error
	GetByID(ctx context.Context, materialID int) (model.Material, error)
	GetAllByUserID(ctx context.Context, userID int) ([]model.Material, error)
	GetAll(ctx context.Context) ([]model.Material, error)
}

// Cards interface provides methods for working with cards.
type Cards interface {
	Create(ctx context.Context, card model.Card) error
	Update(ctx context.Context, card model.Card) error
	Delete(ctx context.Context, userID int, cardID int) error
	GetByID(ctx context.Context, cardID int) (*model.Card, error)
	GetByIdAndUserID(ctx context.Context, userID, cardID int) (model.Card, error)
	GetAllByUserID(ctx context.Context, userID int) (*[]model.Card, error)
	GetAll(ctx context.Context) ([]model.Card, error)
	CreateInvitation(ctx context.Context, inviterID int, card model.Card) error
	GetInvitationsByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error)
	GetInvitationByID(ctx context.Context, userID, cardID, invitationID int) (model.InvitationCard, error)
	GetInvitationsByCardID(ctx context.Context, cardID int) ([]model.InvitationCard, error)
	GetResponsesByUserID(ctx context.Context, userID int) ([]model.InvitationCard, error)
	AcceptCardInvitation(ctx context.Context, userID, cardID, invitationID int, input dto.AnswerInvitation) error
	DeclineCardInvitation(ctx context.Context, userID, cardID, invitationID int, input dto.AnswerInvitation) error
}

type Contests interface {
	Create(ctx context.Context, contest model.Contest) error
	Update(ctx context.Context, contest model.Contest) error
	Delete(ctx context.Context, contestID int) error
	GetByID(ctx context.Context, contestID int) (model.Contest, error)
	GetAll(ctx context.Context) ([]model.Contest, error)
}

type Disciplines interface {
	Create(ctx context.Context, discipline model.Discipline) error
	Update(ctx context.Context, discipline model.Discipline) error
	Delete(ctx context.Context, disciplineID int) error
	GetByID(ctx context.Context, disciplineID int) (model.Discipline, error)
	GetAll(ctx context.Context) ([]model.Discipline, error)
}

type Projects interface {
	Create(ctx context.Context, disciplineID int, project model.Project) error
	Update(ctx context.Context, project model.Project) error
	Delete(ctx context.Context, projectID int) error
	GetAll(ctx context.Context) ([]model.Project, error)
	GetByID(ctx context.Context, projectID int) (model.Project, error)
	GetAllByDisciplineID(ctx context.Context, disciplineID int) ([]model.Project, error)
}

type ProjectModules interface {
	Create(ctx context.Context, module model.ProjectModule) error
	Update(ctx context.Context, module model.ProjectModule) error
	Delete(ctx context.Context, projectID, moduleID int) error
	GetAllByProjectID(ctx context.Context, projectID int) ([]model.ProjectModule, error)
	GetByID(ctx context.Context, moduleID int) (model.ProjectModule, error)
}

// S3Bucket interface provides methods for storing and retrieving objects from an S3 bucket.
type S3Bucket interface {
	PutObject(bucketName string, objectName string, fileBytes []byte) error
}

type Transactional interface {
	Begin(ctx context.Context, o *sql.TxOptions) (*sqlx.Tx, error)
	Commit(ctx context.Context, tx *sqlx.Tx) error
	Rollback(ctx context.Context, tx *sqlx.Tx) error
}

func NewRepository(db *sqlx.DB, sess *session.Session) *Repository {
	return &Repository{
		Users:          NewUsersRepository(db),
		Roles:          NewRolesRepository(db),
		Articles:       NewArticlesRepository(db),
		Cards:          NewCardsRepository(db),
		Materials:      NewMaterialsRepository(db),
		Contests:       NewContestsRepository(db),
		Disciplines:    NewDisciplinesRepository(db),
		Projects:       NewProjectsRepository(db),
		ProjectModules: NewProjectModulesRepository(db),
		S3Bucket:       NewS3BucketRepository(sess),
		Transactional:  NewTransactionManager(db),
	}
}
