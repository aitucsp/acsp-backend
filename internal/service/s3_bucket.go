package service

import (
	"context"
	"os"

	"acsp/internal/repository"
)

// S3BucketService implements the S3Bucket interface.
type S3BucketService struct {
	repo repository.S3Bucket
}

// NewS3BucketService creates a new instance of S3BucketService.
func NewS3BucketService(repo repository.S3Bucket) *S3BucketService {
	return &S3BucketService{repo: repo}
}

func (s *S3BucketService) UploadFile(ctx context.Context, bucket, key string, file *os.File) error {
	return s.repo.AddObject(bucket, key, file)
}
