package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	LOCAL = "local"
	PROD  = "prod"
)

type Paths struct {
	LogDir  string `yaml:"logDir"`
	LogName string `yaml:"logName"`
}

type Server struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type JWT struct {
	Secret  string        `yaml:"secret"`
	Timeout time.Duration `yaml:"timeout"`
}

type Database struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Env      string    `yaml:"env"`
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"database"`
	JWT      *JWT      `yaml:"jwt"`
	Paths    *Paths    `yaml:"paths"`
}

// MustLoad - загружает конфигурацию
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic(any("Файл конфигурации по указанному пути отсутствует"))
	}

	return MustLoadByPath(path)
}

// MustLoadByPath - загружает конфигурацию по указанному пути
func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(any(fmt.Sprintf("Файл конфигурации не найден: %s", configPath)))
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(any(fmt.Sprintf("Ошибка чтения файла конфигурации: %v", err)))
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		path := "local.yaml"

		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = "./config/local.yaml"
			res = path
		}
	}

	return res
}
