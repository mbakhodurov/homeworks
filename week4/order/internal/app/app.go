package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mbakhodurov/homeworks/week4/order/internal/config"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	router      *chi.Mux
	httpServer  *http.Server
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
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initDatabase,
		a.initRouter,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.ApppConfig().Logger.Level(),
		config.ApppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	_ = a.diContainer.PgxConn(ctx)

	migrationRunner := a.diContainer.MigratorRunner(ctx)

	err := migrationRunner.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info(ctx, "🗄️ Database initialized and migrations applied")
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(2 * time.Second))

	orderServer := a.diContainer.OrderV1Server(ctx)
	r.Mount("/", orderServer)

	a.router = r
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = &http.Server{
		Addr:              config.ApppConfig().OrderHTTP.Address(),
		Handler:           a.router,
		ReadHeaderTimeout: config.ApppConfig().OrderHTTP.Readtimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, config.ApppConfig().OrderHTTP.Shutdowntimeout())
		defer cancel()

		err := a.httpServer.Shutdown(shutdownCtx)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.ApppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
