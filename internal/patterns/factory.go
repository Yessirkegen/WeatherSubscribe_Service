package patterns

import "fmt"

type Notification interface {
	Send(message string) error
}

type EmailNotification struct{}

func (n *EmailNotification) Send(message string) error {
	// Логика отправки email
	fmt.Println("Отправка email:", message)
	return nil
}

type SMSNotification struct{}

func (n *SMSNotification) Send(message string) error {
	// Логика отправки SMS
	fmt.Println("Отправка SMS:", message)
	return nil
}

type NotificationFactory struct{}

func (f *NotificationFactory) CreateNotification(notificationType string) (Notification, error) {
	switch notificationType {
	case "email":
		return &EmailNotification{}, nil
	case "sms":
		return &SMSNotification{}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип уведомления: %s", notificationType)
	}
}
