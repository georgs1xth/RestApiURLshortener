package config

// env: "local" #* local, dev, prod
// storage_path: "./storage/storage.db"
// http_server:
// address: "localhost:8082"
// timeout: 4s
// idle_timeout: 60s
type Config struct {
	Env string `yaml:"env" env:"ENV" env-default:"local" env-required:"true`
}