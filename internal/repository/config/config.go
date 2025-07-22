package config

type Config struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	FileStorageName string `env:"FILE_STORAGE_NAME"`
}
