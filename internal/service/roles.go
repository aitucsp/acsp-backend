package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"acsp/internal/apperror"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type RolesService struct {
	repo      repository.Roles
	usersRepo repository.Users
}

func NewRolesService(repo repository.Roles, u repository.Users) *RolesService {
	return &RolesService{repo: repo, usersRepo: u}
}

func (r *RolesService) CreateRole(ctx context.Context, userID, name string) error {
	userId, _ := strconv.Atoi(userID)
	user, err := r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.Wrap(apperror.ErrIncorrectRole, "Incorrect role")
	}

	return r.repo.CreateRole(ctx, name)
}

func (r *RolesService) UpdateRole(ctx context.Context, userID, roleID, newName string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	user, err := r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.Wrap(apperror.ErrIncorrectRole, "Incorrect role")
	}

	roleId, err := strconv.Atoi(roleID)
	if err != nil {
		return err
	}

	return r.repo.UpdateRole(ctx, roleId, newName)
}

func (r *RolesService) DeleteRole(ctx context.Context, userID, roleID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	user, err := r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.Wrap(apperror.ErrIncorrectRole, "Incorrect role")
	}

	roleId, err := strconv.Atoi(roleID)
	if err != nil {
		return err
	}

	return r.repo.DeleteRole(ctx, roleId)
}

func (r *RolesService) SaveUserRole(ctx context.Context, userID, roleID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	user, err := r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.Wrap(apperror.ErrIncorrectRole, "Incorrect role")
	}

	roleId, err := strconv.Atoi(roleID)
	if err != nil {
		return err
	}

	return r.repo.SaveUserRole(ctx, userId, roleId)
}

func (r *RolesService) DeleteUserRole(ctx context.Context, userID, roleID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	user, err := r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.Wrap(apperror.ErrIncorrectRole, "Incorrect role")
	}

	roleId, err := strconv.Atoi(roleID)
	if err != nil {
		return err
	}

	return r.repo.DeleteUserRole(ctx, userId, roleId)
}

func (r *RolesService) GetUserRoles(ctx context.Context, userID string) ([]model.Role, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.Role{}, err
	}

	_, err = r.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return []model.Role{}, errors.Wrap(err, "Error when getting roles from user repository")
	}

	return r.repo.GetUserRoles(ctx, userId)
}
