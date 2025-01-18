package config

type Cacabot struct {
	URL      string `env:"CACABOT_URL" env-required:"true"`
	Username string `env:"CACABOT_USERNAME" env-required:"true"`
	Password string `env:"CACABOT_PASSWORD" env-required:"true"`
	Enabled  bool   `env:"CACABOT_ENABLED" env-required:"true"`
}
