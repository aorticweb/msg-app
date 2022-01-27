package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/aorticweb/msg-app/app/handlers"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConnectionWaitTime time.Duration = 5 * time.Minute
var shutdownTimeout time.Duration = 3 * time.Second
var serverListenAddr string = ":8080"

func dbConnection(url string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(url), &gorm.Config{})
}

func dbUrl() (string, error) {
	url, exist := os.LookupEnv("POSTGRES_URL")
	if !exist {
		return "", errors.New("Env varialbe POSTGRES_URL is not set")
	}
	return url, nil
}

// waitForDB ... wait up to timeout for database to be up then return connection
func waitForDB(timeout time.Duration) (*gorm.DB, error) {
	url, err := dbUrl()
	if err != nil {
		return nil, err
	}
	db, err := dbConnection(url)
	for start := time.Now(); time.Since(start) < timeout; {
		if err == nil {
			break
		}
		db, err = dbConnection(url)
		time.Sleep(1 * time.Second)
	}
	return db, err
}

func registerKillSwitch() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return shutdown
}

func gracefullyShutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err == nil {
		return nil
	}
	logrus.Error("Server shutdown did not succeed in %v: %v", shutdownTimeout, err)
	err = server.Close()

	if err != nil {
		return fmt.Errorf("Server close failed: %v", err)
	}
	logrus.Info("Bye Bye")
	return nil
}

func waitForKillSwitch(kill chan os.Signal, server *http.Server) {
	<-kill
	gracefullyShutdown(server)
}

func setupServer(logger *logrus.Logger) (*http.Server, error) {
	db, err := waitForDB(dbConnectionWaitTime)
	if err != nil {
		logger.Error("Failed to fetch database connection")
		return nil, err
	}
	killSwitch := registerKillSwitch()
	API := api.NewAPI(db, logger)
	server := http.Server{
		Addr:         serverListenAddr,
		Handler:      API,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go waitForKillSwitch(killSwitch, &server)
	return &server, nil
}

func must(err error, logger *logrus.Logger) {
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	logger := logrus.New()
	server, err := setupServer(logger)
	must(err, logger)
	logger.Info("API says Hello\n")
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		logger.Info("API says Goodbye\n")
		os.Exit(0)
	}
	must(err, logger)
}
