package main

import (
	_ "Yakudza/docs"
	"Yakudza/iternal/migrator"
	"Yakudza/iternal/server"
	"Yakudza/pkg/config"
	"Yakudza/pkg/database"
	"Yakudza/pkg/logger"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	//Инициализация конфигурации
	cfg := config.MustLoad()

	//Инициализация логгер
	if err := logger.New(cfg); err != nil {
		panic(any(fmt.Errorf("Ошибка при инициализации логгера: %w", err)))
	}

	//Инициализация БД
	db, err := database.Init(cfg)
	if err != nil {
		panic(any(fmt.Errorf("ошибка инициализации базы данных: %w", err)))
	}

	//Миграции
	{
		err := migrator.Migrations(db)
		if err != nil {
			panic(any(fmt.Errorf("Ошибка при миграции: %w", err)))
		}
	}

	//Создание сервера маршрутизации
	srv := server.New(cfg)

	//Правильное завершение сервиса
	{
		wait := time.Second * 15

		// Запуск сервера в отдельном потоке
		go func() {
			logger.Info("Сервер запущен на адресе: %s", srv.Addr)
			switch cfg.Env {
			case "local":
				{
					if err := srv.ListenAndServe(); err != nil {
						logger.Error("Ошибка при прослушивании сервера: %v", err)
					}
				}
			case "prod":
				{
					if err := srv.ListenAndServeTLS("cert.crt", "key.key"); err != nil {
						logger.Error("Ошибка при прослушивании сервера: %v", err)
					}
				}
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		_ = srv.Shutdown(ctx)
		_ = srv.Close()
		database.CloseConnection()
		logger.Warn("Выключение сервера")
		os.Exit(0)
	}
}
