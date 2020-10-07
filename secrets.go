package config

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

// SecretResolver looks up secrets from a source.
type SecretResolver struct {
	Fetcher Fetcher
}

// NewSecretResolver returns new SecretResolver with an AWSSecretManager as the fetcher.
func NewSecretResolver() *SecretResolver {
	return &SecretResolver{
		Fetcher: NewAWSSecretManager(),
	}
}

// Resolve resolves the secret if needed, and returns the value.
func (r *SecretResolver) Resolve(ctx context.Context, val string) (*string, error) {
	u, err := url.Parse(val)
	if err != nil {
		return nil, err
	}

	if len(u.Scheme) == 0 || strings.HasPrefix(u.Scheme, "postgres") {
		// we can just use the datasource as is.
		return &val, nil
	}

	if u.Scheme != "secretsmanager" {
		// we do not know how to resolve this url
		return nil, errors.New("Unresolvable")
	}

	// We need to retreive the connection string details from the secret manager.
	//
	// The secret name is in the host part of the url.
	name := u.Host

	// get the value from the secret mananager
	b, err := r.Fetcher.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	// the DB secrets are stored in JSON format for RDS, so we need to parse the value.
	secret := RDSSecret{}
	if err := json.Unmarshal(b, &secret); err != nil {
		return nil, err
	}

	// build the connection string from the connection string elements.
	connStr := secret.ConnectionString()

	return &connStr, nil
}
