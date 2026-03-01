package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	inventoryAddress = "localhost:50052"
	paymentAddress   = "localhost:50051"

	httpPort   = "8086"
	configPath = "../deploy/compose/order/.env"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func buildPostgresDSN() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		db,
		sslMode,
	)
}

func main() {
	ctx := context.Background()

	err := godotenv.Load(configPath)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := buildPostgresDSN()
	// dbURI := os.Getenv("DB_URI")

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	// // Инициализируем мигратор
	// migrationsDir := os.Getenv("MIGRATIONS_DIR")
	// migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), "../"+migrationsDir)

	// err = migratorRunner.Up()
	// if err != nil {
	// 	log.Printf("Ошибка миграции базы данных: %v\n", err)
	// 	return
	// }
}
