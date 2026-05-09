package migrator

import "context"

type Migrator interface {
	// Up применяет все новые миграции
	Up(ctx context.Context) error

	// Down откатывает последнюю миграцию
	Down(ctx context.Context) error

	// Reset полностью сбрасывает БД и применяет миграции заново
	Reset(ctx context.Context) error

	// Version возвращает текущую версию миграций
	Version(ctx context.Context) (int64, error)
}
