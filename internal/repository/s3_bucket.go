package repository

import (
	"bytes"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// S3Repository implements the ObjectRepository interface using the AWS SDK for Go.
type S3Repository struct {
	sess *session.Session
}

// NewS3BucketRepository creates a new instance of S3Repository.
func NewS3BucketRepository(sess *session.Session) *S3Repository {
	return &S3Repository{sess: sess}
}

// PutObject adds an object to an S3 bucket.
func (r *S3Repository) PutObject(bucketName string, objectName string, fileBytes []byte) error {
	svc := s3.New(r.sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectName),
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return errors.Wrap(err, "Error occurred when uploading file to S3")
	}

	return nil
}
