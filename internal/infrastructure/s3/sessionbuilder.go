package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"acsp/internal/config"
)

// SessionBuilder configures a client.ConfigProvider.
type SessionBuilder struct {
	awsConfig *config.S3Config
}

// NewSessionBuilder creates a SessionBuilder.
func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{}
}

// WithAWSConfig adds a config.AWSConfig.
func (b *SessionBuilder) WithAWSConfig(a *config.S3Config) *SessionBuilder {
	b.awsConfig = a

	return b
}

// NewSession creates a client.ConfigProvider.
func (b *SessionBuilder) NewSession() (*session.Session, error) {
	lv := aws.LogDebug

	c := aws.NewConfig()

	c = c.WithLogLevel(lv)
	c = c.WithCredentials(
		credentials.NewStaticCredentials(
			b.awsConfig.AccessToken,
			b.awsConfig.SecretKey,
			""),
	)

	c = c.WithEndpoint(b.awsConfig.Endpoint)

	s, err := session.NewSession(c, aws.NewConfig().WithRegion(b.awsConfig.Region))
	if err != nil {
		return nil, err
	}

	return s, err
}
