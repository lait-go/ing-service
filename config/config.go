package conf

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var Cfg *Config
var Token string

type Config struct {
	Env     string `yaml:"env"`
	DbToken string
	Path    Ways   `yaml:"path"`
}

type Ways struct {
	LogPath     string `yaml:"log_path"`
	ErrorPath   string `yaml:"error_path"`
}

func Config_work()*Config{
	var conf Config
	
	godotenv.Load()

	defPath := os.Getenv("CONFIG_PATH")
	conf.DbToken = os.Getenv("DB_TOKEN")
	Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	if defPath == "" {
		defPath = "../config/config.yaml"
	}
	
	cleanenv.ReadConfig(defPath, &conf)

	return &conf
}