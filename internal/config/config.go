package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	loggerConf struct {
		Level string `yaml:"level"`
	}

	memoryConf struct {
		CleanUpInterval time.Duration `yaml:"cleanup_interval"`
	}

	databaseConf struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	}

	storageConf struct {
		Memory memoryConf   `yaml:"memory"`
		DB     databaseConf `yaml:"db"`
	}

	grpcConf struct {
		Port string `yaml:"port"`
	}

	serverConf struct {
		GRPC grpcConf `yaml:"grpc"`
	}

	constraintsConf struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
		IP       string `yaml:"ip"`
	}

	Config struct {
		Logger      loggerConf      `yaml:"logger"`
		Storage     storageConf     `yaml:"storage"`
		Server      serverConf      `yaml:"server"`
		Constraints constraintsConf `yaml:"constraints"`
	}
)

func New(path string) (Config, error) {
	var cfg Config
	err := config(path, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func config(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}
