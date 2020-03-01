package common

import (
	"gopkg.in/ini.v1"
)

// Config config
type Config struct {
	Common `ini:"common"`
	Redis  `ini:"redis"`
}

// Common config
type Common struct {
	DataFolder  string `ini:"dataFolder"`
	Port        int    `ini:"port"`
	Repo        string `ini:"repo"`
	Mode        string `ini:"mode"`
	LogFilePath string `ini:"logFilePath"`
	LogFileName string `ini:"logFileName"`
	LogLevel    int    `ini:"logLevel"`
}

// Redis config
type Redis struct {
	Host         string `ini:"host"`
	DB           int    `ini:"db"`
	CachePrefix  string `ini:"cachePrefix"`
	SortedPrefix string `ini:"sortedPrefix"`
}

// InitConfig read config from file
func InitConfig(configFileMame string) (*Config, error) {
	cfg, err := ini.Load(configFileMame)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = cfg.MapTo(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
