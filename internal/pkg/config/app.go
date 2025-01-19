package config

type App struct {
	Name    string `env:"APP_NAME" env-default:"Golang-Boilerplate"`
	Version string `env:"APP_VERSION" env-required:"true"`
}
