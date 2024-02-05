package config

type AppConfig struct {
	// -------------------
	// Application configs
	// -------------------
	EnvironmentName  string `envconfig:"APP_ENV" default:"development"`
	LocalEnvFilePath string `envconfig:"APP_LOCAL_ENV_FILE_PATH" default:".env"`
	Debug            bool   `envconfig:"APP_DEBUG" default:"false"`
}

func (c *AppConfig) GetEnvironmentName() string {
	return c.EnvironmentName
}

func (c *AppConfig) GetLocalEnvFilePath() string {
	return c.LocalEnvFilePath
}

func (c *AppConfig) IsDebug() bool {
	return c.Debug
}
