package service

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"acsp/internal/constants"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type UserService struct {
	repo repository.Users
}

func (u UserService) DeleteUser(ctx context.Context, userID string) error {
	// TODO implement me
	panic("implement me")
}

func NewUsersService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (u UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	l := logging.LoggerFromContext(ctx)

	users, err := u.repo.GetAll(ctx)
	if err != nil {
		l.Error("Error getting all users from database", zap.Error(err))

		return nil, err
	}

	users = u.getFullURLForUsers(users)

	return users, nil
}

func (u UserService) CreateUser(ctx context.Context, dto dto.CreateUser) error {
	// TODO implement me
	panic("implement me")
}

func (u UserService) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("userID", userID))
	userId, err := strconv.Atoi(userID)
	if err != nil {
		l.Error("Error converting string to int", zap.Error(err))

		return model.User{}, err
	}

	user, err := u.repo.GetByID(ctx, userId)
	if err != nil {
		l.Error("Error getting user by id from database", zap.Error(err))

		return model.User{}, err
	}

	userInfo, err := u.repo.GetUserDetailsByUserId(ctx, userId)
	if err != nil {
		l.Error("Error getting user details by id from database", zap.Error(err))

		return model.User{}, err
	}

	user.UserInfo = userInfo

	user = u.getFullUrl(user)

	return user, nil
}

func (u UserService) UpdateUserImageURL(ctx context.Context, userID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	return u.repo.UpdateImageURL(ctx, userId)
}

func (u UserService) UpdateUser(ctx context.Context, userID string, dto dto.UpdateUser) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	userDetails := model.UserDetails{
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		PhoneNumber:    dto.PhoneNumber,
		Specialization: dto.Specialization,
	}

	return u.repo.UpdateDetails(ctx, userId, userDetails)
}

// getFullURLForArticles function gets a slice of articles and changes every article's image_url to a full url
func (u UserService) getFullURLForUsers(users []model.User) []model.User {
	for i := range users {
		users[i].ImageURL = constants.BucketName + "." +
			constants.EndPoint + "/" +
			constants.UsersAvatarsFolder +
			users[i].ImageURL
	}

	return users
}

// getFullUrl function gets an article and changes its image_url to a full url
func (u UserService) getFullUrl(user model.User) model.User {
	user.ImageURL = constants.BucketName + "." +
		constants.EndPoint + "/" +
		constants.UsersAvatarsFolder +
		user.ImageURL

	return user
}
