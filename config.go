package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tokenized/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	// EnvProduction is the expected environment variable for production.
	//
	// An empty ENV var also indicates production.
	EnvProduction = "prod"

	// masked is the string to use when masking config values.
	masked = "******"

	// EnvParamName is the the environment variable to use when loading
	// config from the AWS ParamStore.
	EnvParamName = "PARAM_NAME"

	// EnvConfigFile is the the environment variable to use when loading
	// JSON config from a file.
	EnvConfigFile = "CONFIG_FILE"
)

// LoadConfig is the common way to load config.
//
// It will attempt to load from param store, then config file, falling back
// to the environment.
//
// To load from the ParamStore, the PARAM_STORE env var should be set with the
// name of the item to load.
//
// To load from a JSON config file, the CONFIG_FILE env var should have the name
// of the file to load.
//
// Fallback option is to load config from environment variables.
func LoadConfig(ctx context.Context, cfg interface{}) error {
	// check the PARAM_NAME env var
	paramName := os.Getenv(EnvParamName)

	if len(paramName) > 0 {
		// we have a parameter name, try to load it
		logger.Info(ctx, "Loading config from param store : %s", paramName)
		return LoadParamStore(paramName, cfg)
	}

	// check the CONFIG_FILE env var
	filename := os.Getenv(EnvConfigFile)

	if len(filename) > 0 {
		logger.Info(ctx, "Loading config from file : %s", filename)

		return LoadFromFile(filename, cfg)
	}

	logger.Info(ctx, "Loading config from environment")
	return LoadEnvironment(cfg)
}

// LoadFromFile loads a JSON config from a file.
func LoadFromFile(filename string, cfg interface{}) error {
	// Load default values from environment definitions.
	if err := LoadEnvironment(cfg); err != nil {
		return errors.Wrap(err, "load environment defaults")
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	if err := json.Unmarshal(b, cfg); err != nil {
		return err
	}

	return nil
}

// LoadEnvironment attempts to hydrate a struct with environment variables.
func LoadEnvironment(cfg interface{}) error {
	return envconfig.Process("", cfg)
}

// LoadParamStore returns unmarshals the an AWS ParamStore value into a
// struct.
//
// This function is intended to work for any struct that is to be used for
// service configuration.
//
// A difference between ParamStore is that this function will not flatten the
// JSON config that it loads, which is expected to fit the struct it is being
// marshalled into.
//
// It is intended to eventually replace the usage of ParamStore, which
// requires a specific type.
func LoadParamStore(keyName string, cfg interface{}) error {
	// Load default values from environment definitions.
	if err := LoadEnvironment(cfg); err != nil {
		return errors.Wrap(err, "load environment defaults")
	}

	b, err := fetchFromParamStore(keyName)
	if err != nil {
		return errors.Wrap(err, "fetch param store")
	}

	// Unmarshal param value
	return json.Unmarshal(b, cfg)
}

// fetchFromParamStore loads the data from the AWS ParamStore.
func fetchFromParamStore(keyName string) ([]byte, error) {
	// Locate param value
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, errors.Wrap(err, "new session")
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

	withDecryption := true

	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyName,
		WithDecryption: &withDecryption,
	})

	if err != nil {
		return nil, errors.Wrap(err, "get parameters")
	}

	// Unmarshal param value
	b := []byte(*param.Parameter.Value)

	return b, nil
}

// DumpSafe logs a "safe" version of the config, with sensitive values masked.
func DumpSafe(ctx context.Context, cfg interface{}) {
	logger.Info(ctx, "Config : %+v", Mask(cfg))
}
