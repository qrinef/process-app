package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	AppEnvironmentNameVariable = "APP_ENV"
	AppEnvFilePathVariableName = "APP_LOCAL_ENV_FILE_PATH"

	EnvDev = "development"
)

var (
	ErrVariableEmptyButRequired = errors.New("env variables is empty and has required tag")
)

func LoadLocalEnvIfDev() error {
	value, isEnvVariableExists := os.LookupEnv(AppEnvironmentNameVariable)
	if !isEnvVariableExists {
		return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, AppEnvironmentNameVariable)
	}

	if value == EnvDev {
		envFilePath, isExists := os.LookupEnv(AppEnvFilePathVariableName)
		if !isExists {
			return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, AppEnvFilePathVariableName)
		}

		loadErr := godotenv.Load(envFilePath)
		if loadErr != nil {
			return loadErr
		}
	}

	return nil
}

func Prepare(config interface{}) error {
	err := LoadLocalEnvIfDev()
	if err != nil {
		return err
	}

	err = envconfig.Process("", config)
	if err != nil {
		return err
	}

	return nil
}
