package service

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/dto"
	"acsp/internal/logging"
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

func (c *ContestsService) Create(ctx context.Context, contest dto.CreateContest) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("contestName", contest.Name))

	contestModel := model.Contest{
		Name:        contest.Name,
		Description: contest.Description,
		Link:        contest.Link,
		StartDate:   contest.StartDate,
		EndDate:     contest.EndDate,
	}

	err := c.repo.Create(ctx, contestModel)
	if err != nil {
		l.Error("Error when creating the contest", zap.Error(err))

		return errors.Wrap(err, "error when creating the contest")
	}

	return nil
}

func (c *ContestsService) Update(ctx context.Context, contestID string, contest dto.UpdateContest) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("contestID", contestID),
		zap.String("contestName", contest.Name))

	contestModel := model.Contest{
		Name:        contest.Name,
		Description: contest.Description,
		Link:        contest.Link,
		StartDate:   contest.StartDate,
		EndDate:     contest.EndDate,
	}

	err := c.repo.Update(ctx, contestModel)
	if err != nil {
		l.Error("Error when updating the contest", zap.Error(err))

		return errors.Wrap(err, "error when updating the contest")
	}

	return nil
}
func (c *ContestsService) Delete(ctx context.Context, contestID string) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("contestID", contestID))

	contestId, err := strconv.Atoi(contestID)
	if err != nil {
		l.Error("Error when converting id to int", zap.Error(err))

		return errors.Wrap(err, "error when converting id to int")
	}

	err = c.repo.Delete(ctx, contestId)
	if err != nil {
		l.Error("Error when deleting the contest", zap.Error(err))

		return errors.Wrap(err, "error when deleting the contest")
	}

	return nil
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
