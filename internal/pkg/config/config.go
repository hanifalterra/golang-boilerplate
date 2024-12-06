package config

type Config struct {
	App    App
	DB     DB
	Redis  Redis
	Logger Logger
}

type App struct {
	Name    string `env-required:"true" env:"APP_NAME"`
	Version string `env-required:"true" env:"APP_VERSION"`
}
