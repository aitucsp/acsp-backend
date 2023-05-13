package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"acsp/internal/model"
	"acsp/internal/repository"
)

type ContestsService struct {
	repo repository.Contests
}

func NewContestsService(r repository.Contests) *ContestsService {
	return &ContestsService{
		repo: r,
	}
}

func (c *ContestsService) GetByID(ctx context.Context, contestID string) (model.Contest, error) {
	contestId, err := strconv.Atoi(contestID)
	if err != nil {
		return model.Contest{}, errors.Wrap(err, "error when converting id to int")
	}

	contest, err := c.repo.GetByID(ctx, contestId)
	if err != nil {
		return model.Contest{}, errors.Wrap(err, "error when getting the contest")
	}

	return contest, nil
}

func (c *ContestsService) GetAll(ctx context.Context) ([]model.Contest, error) {
	return c.repo.GetAll(ctx)
}
