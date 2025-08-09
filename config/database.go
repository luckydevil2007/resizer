package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	PostgresConfig   DbConfig `yaml:"postgres"`
	ClickhouseConfig DbConfig `yaml:"clickhouse"`
}

func NewConfig() *Config {

	var config Config

	content, err := os.ReadFile("./config/db.yaml")
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil
	}

	return &config
}

func FormatDSN(dbCponfig DbConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCponfig.Host, dbCponfig.Port, dbCponfig.User, dbCponfig.Password, dbCponfig.DBName)
}
