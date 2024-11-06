package main

import (
	"WeatherSubs/internal/config"
	"WeatherSubs/internal/handlers"
	"WeatherSubs/internal/models"
	"WeatherSubs/internal/repository"
	"WeatherSubs/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Получаем конфигурацию
	cfg := config.GetConfig()

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// Миграции
	err = db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Подключение к RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала RabbitMQ: %s", err)
	}
	defer ch.Close()

	// Объявление обменника
	err = ch.ExchangeDeclare(
		"notifications_exchange", // имя обменника
		"fanout",                 // тип
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления обменника: %s", err)
	}

	// Инициализация репозитория и сервиса
	repo := repository.NewSubscriptionRepository(db)

	// Инициализация клиента погодного API
	weatherAPI := services.NewWeatherServiceClient(cfg.WeatherServiceURL)

	service := services.NewSubscriptionService(repo, ch, "notifications_exchange", "", weatherAPI)

	// Инициализация хендлеров
	handler := handlers.NewSubscriptionHandler(service)

	// Настройка роутера
	router := gin.Default()

	// Middleware для аутентификации (заглушка)
	router.Use(dummyAuthMiddleware)

	// Маршруты
	router.POST("/subscriptions", handler.CreateSubscription)
	router.GET("/subscriptions", handler.GetSubscriptions)
	router.DELETE("/subscriptions/:id", handler.DeleteSubscription)

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}

// Заглушка для Middleware аутентификации
func dummyAuthMiddleware(c *gin.Context) {
	// В реальном приложении вы должны проверять JWT или сессии
	// Здесь мы просто устанавливаем userID в контекст
	c.Set("userID", uint(1))
	c.Next()
}
