package services

import (
	"WeatherSubs/internal/models"
	"WeatherSubs/internal/patterns"
	"WeatherSubs/internal/repository"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
)

type SubscriptionService interface {
	CreateSubscription(subscription *models.Subscription) error
	GetSubscriptions(userID uint) ([]models.Subscription, error)
	DeleteSubscription(id uint) error
}

type subscriptionService struct {
	repo       repository.SubscriptionRepository
	factory    *patterns.NotificationFactory
	rabbitMQ   *amqp.Channel
	exchange   string
	routingKey string
	schedulers map[string]*gocron.Scheduler
	weatherAPI WeatherAPIClient
}

func NewSubscriptionService(repo repository.SubscriptionRepository, rabbitMQ *amqp.Channel, exchange, routingKey string, weatherAPI WeatherAPIClient) SubscriptionService {
	service := &subscriptionService{
		repo:       repo,
		factory:    &patterns.NotificationFactory{},
		rabbitMQ:   rabbitMQ,
		exchange:   exchange,
		routingKey: routingKey,
		schedulers: make(map[string]*gocron.Scheduler),
		weatherAPI: weatherAPI,
	}

	// Запускаем планировщики для существующих подписок
	service.scheduleNotifications()

	return service
}

func (s *subscriptionService) scheduleNotifications() {
	subscriptions, err := s.repo.GetAll()
	if err != nil {
		fmt.Println("Ошибка получения подписок:", err)
		return
	}

	for _, sub := range subscriptions {
		s.scheduleNotification(sub)
	}
}

func (s *subscriptionService) scheduleNotification(subscription models.Subscription) {
	location, err := time.LoadLocation(subscription.Timezone)
	if err != nil {
		fmt.Printf("Ошибка загрузки часового пояса %s: %v\n", subscription.Timezone, err)
		return
	}

	scheduler := gocron.NewScheduler(location)
	scheduler.Every(1).Day().At("08:00").Do(func() {
		s.sendNotification(subscription)
	})
	scheduler.StartAsync()

	// Сохраняем планировщик
	key := strconv.FormatUint(uint64(subscription.ID), 10)
	s.schedulers[key] = scheduler
}

func (s *subscriptionService) sendNotification(subscription models.Subscription) {
	// Получаем прогноз погоды для города
	weatherData, err := s.weatherAPI.GetWeather(subscription.City)
	if err != nil {
		fmt.Printf("Ошибка получения погоды для города %s: %v\n", subscription.City, err)
		return
	}

	// Генерируем совет по одежде
	advice := generateClothingAdvice(weatherData)

	// Формируем сообщение
	message := fmt.Sprintf("Прогноз погоды для %s: %s. Совет: %s", subscription.City, weatherData.Summary, advice)

	// Публикуем сообщение в RabbitMQ
	err = s.publishNotification(subscription.Type, message)
	if err != nil {
		fmt.Printf("Ошибка публикации уведомления: %v\n", err)
	}
}

func generateClothingAdvice(weatherData WeatherData) string {
	if weatherData.Temperature < 0 {
		return "Наденьте теплую куртку и шапку"
	} else if weatherData.Temperature < 10 {
		return "Наденьте куртку"
	} else if weatherData.Temperature < 20 {
		return "Легкая куртка будет достаточно"
	} else {
		return "Одевайтесь по-летнему"
	}
}

func (s *subscriptionService) publishNotification(notificationType, message string) error {
	notificationMessage := fmt.Sprintf("%s:%s", notificationType, message)
	err := s.rabbitMQ.Publish(
		s.exchange,
		s.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(notificationMessage),
		},
	)
	return err
}

func (s *subscriptionService) CreateSubscription(subscription *models.Subscription) error {
	err := s.repo.Create(subscription)
	if err != nil {
		return err
	}

	// Планируем уведомление для новой подписки
	s.scheduleNotification(*subscription)

	return nil
}

func (s *subscriptionService) GetSubscriptions(userID uint) ([]models.Subscription, error) {
	return s.repo.GetByUserID(userID)
}

func (s *subscriptionService) DeleteSubscription(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Останавливаем и удаляем планировщик
	key := strconv.FormatUint(uint64(id), 10)
	if scheduler, exists := s.schedulers[key]; exists {
		scheduler.Stop()
		delete(s.schedulers, key)
	}

	return nil
}
