package config

import (
	"encoding/json"
	"io"
	"os"
)

type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

func parseJSON(configPath string) (JSONConfig, error) {
	var cfg JSONConfig
	if configPath == "" {
		return cfg, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return cfg, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil && err != io.EOF {
		return cfg, err
	}

	return cfg, nil
}

func applyJSONConfig(conf *Config, jsonCfg JSONConfig) {
	setJSONStringValue(&conf.Handlers.ServerAddr, DefaultServerAddr, jsonCfg.ServerAddress)
	setJSONStringValue(&conf.Service.ServerURL, DefaultServerURL, jsonCfg.BaseURL)
	setJSONStringValue(&conf.FileStorage.FileStoragePath, DefaultFileStorage, jsonCfg.FileStoragePath)
	setJSONStringValue(&conf.DB.DatabaseDsn, DefaultDatabaseDsn, jsonCfg.DatabaseDSN)
	setJSONBoolValue(&conf.Handlers.EnableHTTPS, DefaultHTTPS, jsonCfg.EnableHTTPS)
}

func setJSONStringValue(confValue *string, defaultValue string, jsonValue string) {
	if *confValue == defaultValue && jsonValue != "" {
		*confValue = jsonValue
	}
}

func setJSONBoolValue(confValue *bool, defaultValue bool, jsonValue bool) {
	if *confValue == defaultValue {
		*confValue = jsonValue
	}
}
