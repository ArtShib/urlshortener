package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileConfig структура фалового конфига
type FileConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// LoadConfigFile загрузка конфига из файла
func LoadConfigFile(path string) (*FileConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *FileConfig

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *FileConfig) LoadConfig(config *Config) {
	if config.HTTPServer.ServerAddress == "" {
		config.HTTPServer.ServerAddress = cfg.ServerAddress
	}
	if config.ShortService.BaseURL == "" {
		config.ShortService.BaseURL = cfg.BaseURL
	}
	if config.RepoConfig.FileStoragePath == "" {
		config.RepoConfig.FileStoragePath = cfg.FileStoragePath
	}
	if config.RepoConfig.DatabaseDSN == "" {
		config.RepoConfig.DatabaseDSN = cfg.DatabaseDSN
	}
	if !config.TLSConfig.Enabled {
		config.TLSConfig.Enabled = cfg.EnableHTTPS
	}
}
