package config

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// AWSSecretManager is used to connect to the AWS Secrets Manager service.
type AWSSecretManager struct {
}

// NewAWSSecretManager returns an AWSSecretManager
func NewAWSSecretManager() *AWSSecretManager {
	return &AWSSecretManager{}
}

// Get returns the contents of a secret.
func (s *AWSSecretManager) Get(ctx context.Context, name string) ([]byte, error) {
	// build the input
	in := secretsmanager.GetSecretValueInput{
		SecretId: &name,
	}

	// make the request
	out, err := s.client().GetSecretValue(&in)
	if err != nil {
		return nil, err
	}

	// the secret may represent JSON, or primitive values, so return bytes
	return []byte(*out.SecretString), nil
}

// client returns a client for communicating with the AWS Secrets Manager
// service.
func (s *AWSSecretManager) client() *secretsmanager.SecretsManager {
	return secretsmanager.New(session.New())
}
