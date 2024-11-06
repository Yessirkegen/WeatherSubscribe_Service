package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	RabbitMQURL       string
	WeatherServiceURL string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		// Загружаем переменные окружения из .env
		err := godotenv.Load()
		if err != nil {
			log.Println("Ошибка загрузки .env файла, используем переменные окружения")
		}

		instance = &Config{
			DatabaseURL:       os.Getenv("DATABASE_URL"),
			RabbitMQURL:       os.Getenv("RABBITMQ_URL"),
			WeatherServiceURL: os.Getenv("WEATHER_SERVICE_URL"),
		}
	})

	return instance
}
