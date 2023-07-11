package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
	"os"
)

const Namespace = "service_api_test"

type ServerData struct {
	Port   string `yaml:"port"`
	ApiKey string `yaml:"api_key"`
}

type DBSQLConnection struct {
	Uri                   string `yaml:"uri"`
	ConnMaxLifetimeMinute int    `yaml:"conn_max_lifetime_minute" default:"1"`
	MaxOpenConn           int    `yaml:"max_open_conn" default:"10"`
}
type Cache struct {
	TtlSecond int `yaml:"ttl_second" default:"60"`
	Size      int `yaml:"size" default:"1000"`
}
type Config struct {
	Server          ServerData      `yaml:"server"`
	DBSQLConnection DBSQLConnection `yaml:"postgres"`
	Cache           Cache           `yaml:"cache"`
}

func NewConfig(configFile string) (*Config, error) {
	var config Config
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failure upload yaml file. err %w", err)
	}

	err = defaults.Set(&config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &config, nil
}
