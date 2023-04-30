package service

import (
	"context"
	"strconv"

	"acsp/internal/dto"
	"acsp/internal/repository"
)

type UserService struct {
	repo repository.Authorization
}

func NewUsersService(repo repository.Authorization) *UserService {
	return &UserService{repo: repo}
}

func (u UserService) CreateUser(ctx context.Context, dto dto.CreateUser) error {
	// TODO implement me
	panic("implement me")
}

// func (u UserService) GetUserByID(ctx context.Context, userID string) (model.User, error) {
// 	// TODO implement me
// 	panic("implement me")
// }

func (u UserService) UpdateUserImageURL(ctx context.Context, userID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	return u.repo.UpdateImageURL(ctx, userId)
}

// func (u UserService) UpdateUser(ctx context.Context, dto dto.UpdateUser) error {
// 	// TODO implement me
// 	panic("implement me")
// }
