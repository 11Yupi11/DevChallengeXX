package config

type Config struct {
	Port  int  `yaml:"port" env:"APP_PORT"`
	Debug bool `yaml:"debug" env:"APP_DEBUG"`
}
