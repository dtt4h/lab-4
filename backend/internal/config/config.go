package config

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	URL    string
	Driver string
}

type JWTConfig struct {
	SecretKey string
}
