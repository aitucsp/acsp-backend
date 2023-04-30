package service

import (
	"context"
	"mime/multipart"

	"github.com/pkg/errors"

	"acsp/internal/repository"
)

// S3BucketService implements the S3Bucket interface.
type S3BucketService struct {
	repo   repository.S3Bucket
	bucket string
}

// NewS3BucketService creates a new instance of S3BucketService.
func NewS3BucketService(repo repository.S3Bucket, b string) *S3BucketService {
	return &S3BucketService{repo: repo, bucket: b}
}

func (s *S3BucketService) UploadFile(ctx context.Context, key string, file *multipart.FileHeader) error {
	// Open the file
	f, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "Error occurred when opening file")
	}

	// Read the file contents into a byte slice
	fileBytes := make([]byte, file.Size)
	_, err = f.Read(fileBytes)
	if err != nil {
		return errors.Wrap(err, "Error occurred when reading file")
	}

	return s.repo.PutObject(s.bucket, key, fileBytes)
}
