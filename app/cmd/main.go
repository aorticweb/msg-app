package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aorticweb/msg-app/app/crud"
	api "github.com/aorticweb/msg-app/app/handlers"
)

var dbConnectionWaitTime time.Duration = 5 * time.Minute
var shutdownTimeout time.Duration = 3 * time.Second
var serverListenAddr string = ":3001" // TODO: get from environment

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
	log.Printf("Server shutdown did not succeed in %v: %v\n", shutdownTimeout, err)
	err = server.Close()

	if err != nil {
		return fmt.Errorf("Server close failed: %v", err)
	}
	return nil
}

func waitForKillSwitch(kill chan os.Signal, server *http.Server) {
	<-kill
	gracefullyShutdown(server)
}

func setupServer(logger *log.Logger) (*http.Server, error) {
	db, err := crud.WaitForDB(dbConnectionWaitTime)
	if err != nil {
		logger.Println("Failed to fetch database connection")
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

func must(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	logger := log.New(os.Stdout, "msg-app: ", log.LstdFlags|log.Lshortfile)
	server, err := setupServer(logger)
	must(err, logger)
	logger.Println("API says Hello")
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		logger.Println("API says Goodbye")
		os.Exit(0)
	}
	must(err, logger)
}
