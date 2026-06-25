package app

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/mbakhodurov/homeworks/week7/notification/internal/config"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер OrderPaid
	go func() {
		if err := a.runOrderPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("order paid consumer crashed: %v", err)
		}
	}()

	// Консьюмер OrderAssembled
	go func() {
		if err := a.runOrderAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("order assembled consumer crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initTelegramBot,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDi(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	// return logger.Init(
	// 	config.AppConfig().Logger.Level(),
	// 	config.AppConfig().Logger.AsJson(),
	// 	true,
	// )
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.EnableOTLP(),
		config.AppConfig().Logger.OtlpEndpoint(),
		config.AppConfig().Logger.ServiceName(),
		config.AppConfig().Logger.ServiceEnv(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order Paid consumer running")

	err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runOrderAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order Assembled consumer running")

	err := a.diContainer.OrderAssembledConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot := a.diContainer.TelegramBot(ctx)

	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🛸 Cosmo Bot активирован! Теперь вы будете получать уведомления о постройке космических кораблей и сделанных заказах!",
		})
		if err != nil {
			logger.Error(ctx, "Failed to send activation message", zap.Error(err))
		}
	})

	go func() {
		logger.Info(ctx, "🤖 Telegram bot started...")
		telegramBot.Start(ctx)
	}()

	return nil
}
