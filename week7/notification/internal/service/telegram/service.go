package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"github.com/mbakhodurov/homeworks/week7/notification/internal/client/http"
	"github.com/mbakhodurov/homeworks/week7/notification/internal/model"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

const chatID = 886449792

//go:embed templates/assembled_notification.tmpl
//go:embed templates/paid_notification.tmpl
var templateFS embed.FS

type assembledTemplateData struct {
	EventUUID    string // Уникальный идентификатор события (для идемпотентности)
	OrderUUID    string // Идентификатор собранного заказа
	UserUUID     string // Идентификатор пользователя
	BuildTimeSec int64  // Время (в секундах), потраченное на сборку корабля
}

type paidTemplateData struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

var (
	assembledTemplate = template.Must(template.ParseFS(templateFS, "templates/assembled_notification.tmpl"))
	paidTemplate      = template.Must(template.ParseFS(templateFS, "templates/paid_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembled) error {
	message, err := s.buildShipAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildShipAssembledMessage(event model.ShipAssembled) (string, error) {
	data := assembledTemplateData{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaid) error {
	message, err := s.buildOrderPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}
	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(event model.OrderPaid) (string, error) {
	data := paidTemplateData{
		EventUUID:       event.EventUUID,
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}

	var buf bytes.Buffer

	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
