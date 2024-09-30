package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"songs-lib/internal/config"
	"songs-lib/internal/delivery/http"
	"songs-lib/internal/repository"
	"songs-lib/internal/repository/pg"
	"songs-lib/internal/server"
	"songs-lib/internal/service"
	"songs-lib/migrations"
	"syscall"
)

// @title Online Song's Library
// @version 1.0
// description API Server for OnlineSongLibrary Application

// @host localhost:8000
// @basePath /api

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := config.InitConfig("configs", "config"); err != nil {
		logrus.Fatalf("init config err: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("load .env file err: %s", err.Error())
	}

	db, err := pg.NewPostgresDB(pg.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("init db err: %s", err.Error())
	}

	if db == nil {
		logrus.Fatal("db is nil")
		return
	}

	if err := migrations.RunMigrations(db.DB); err != nil {
		logrus.Fatalf("migration err: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := http.NewHandler(services)

	srv := new(server.Server)
	go func() {
		router := handlers.InitRoutes()

		if err := srv.Run(viper.GetString("port"), router); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured while shutting down http server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured while closing db: %s", err.Error())

	}
}
