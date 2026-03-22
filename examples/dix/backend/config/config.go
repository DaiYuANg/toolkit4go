package config

type AppConfig struct {
	Server struct {
		Port int `koanf:"port"`
	} `koanf:"server"`
	DB struct {
		DSN string `koanf:"dsn"`
	} `koanf:"db"`
}
