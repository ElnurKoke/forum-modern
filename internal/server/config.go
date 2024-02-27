package server

import (
	"encoding/json"
	"forum/internal/models"
	"os"
)

type Config struct {
	Port string
	DB   struct {
		Dsn    string
		Driver string
	}
}

func NewConfig() (Config, error) {
	// Открываем JSON-файл с конфигурацией.
	configFile, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	// Декодируем JSON-файл в структуру Config.
	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	models.InfoLog.Println("Configuration extraction successful")
	return config, nil
}
