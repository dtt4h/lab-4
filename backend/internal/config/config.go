package config

type Config struct {
	Env      string
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	Secret string
}
