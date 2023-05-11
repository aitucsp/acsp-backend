package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type MaterialsService struct {
	repo      repository.Materials
	usersRepo repository.Authorization
}

func NewMaterialsService(repo repository.Materials, usersRepo repository.Authorization) *MaterialsService {
	return &MaterialsService{repo: repo, usersRepo: usersRepo}
}

func (m *MaterialsService) Create(ctx context.Context, userID string, dto dto.CreateMaterial) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	user, err := m.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	material := model.Material{
		Topic:       dto.Topic,
		Description: dto.Description,
		Author:      user,
	}

	return m.repo.Create(ctx, material)
}

func (m *MaterialsService) Update(ctx context.Context, materialID, userID string, materialDto dto.UpdateMaterial) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	user, err := m.usersRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	materialId, err := strconv.Atoi(materialID)
	if err != nil {
		return errors.Wrap(err, "error converting material id to int")
	}

	material := model.Material{
		ID:          materialId,
		Topic:       materialDto.Topic,
		Description: materialDto.Description,
		Author:      user,
	}

	return m.repo.Update(ctx, material)
}

func (m *MaterialsService) Delete(ctx context.Context, userID, materialID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(err, "error converting user id to int")
	}

	projectID, err := strconv.Atoi(materialID)
	if err != nil {
		return errors.Wrap(err, "error converting material id to int")
	}

	return m.repo.Delete(ctx, userId, projectID)
}

func (m *MaterialsService) GetAll(ctx context.Context) ([]model.Material, error) {
	return m.repo.GetAll(ctx)
}

func (m *MaterialsService) GetAllByUserID(ctx context.Context, userID string) ([]model.Material, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return []model.Material{}, err
	}

	return m.repo.GetAllByUserID(ctx, userId)
}

func (m *MaterialsService) GetByID(ctx context.Context, materialID string) (model.Material, error) {
	materialId, err := strconv.Atoi(materialID)
	if err != nil {
		return model.Material{}, err
	}

	return m.repo.GetByID(ctx, materialId)
}
