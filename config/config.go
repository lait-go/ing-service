package conf

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var Cfg *Config
var Token string

type Config struct {
	Env     string `yaml:"env"`
	Address string `yaml:"address"`
	Path    Ways   `yaml:"path"`
}

type Ways struct {
	LogPath     string `yaml:"log_path"`
	ErrorPath   string `yaml:"error_path"`
}

func Config_work()*Config{
	var conf Config

	defPath := os.Getenv("CONFIG_PATH")
	if defPath == "" {
		defPath = "../config/config.yaml"
	}
	
	cleanenv.ReadConfig(defPath, &conf)

	Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	return &conf
}