package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/orderservice/infrastructure"
	"orderservice/pkg/orderservice/model"
	"os"
	"os/signal"
	"syscall"
)

const (
	dbUser     = "postgres"
	dbPassword = "gv9y3ytsow"
	dbName     = "orderservice"
	serverPort = ":8091"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
	log.WithFields(log.Fields{"url": serverPort}).Info("starting the server")
	killSignalChan := getKillSignalChan()

	srv := startServer(serverPort)
	waitForKillSignal(killSignalChan)
	_ = srv.Shutdown(context.Background())
}

func startServer(serverPort string) *http.Server {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	newRepo := infrastructure.NewRepo(db)
	newQueryService := infrastructure.NewQueryService(db)
	orderService := model.NewOrderService(newRepo)
	server := infrastructure.Server{OrderService: orderService, OrderQueryService: newQueryService}
	handler := infrastructure.Router(&server)
	srv := &http.Server{Addr: serverPort, Handler: handler}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("get SIGTERM")
	}
}
