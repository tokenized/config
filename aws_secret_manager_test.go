package config

import (
	"context"
	"testing"
)

func TestAWSSecretManager_Get(t *testing.T) {
	t.Skip("Only run when testing AWS connectivity")

	// name of the secret to fetch
	name := "a-secret"

	f := NewAWSSecretManager()

	ctx := context.Background()

	got, err := f.Get(ctx, name)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 {
		t.Errorf("len should have been > 0 bytes")
	}
}
