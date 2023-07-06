package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	consumer "github.com/yoratyo/go-redpanda/kafka"
	"github.com/yoratyo/go-redpanda/repository"
	"github.com/yoratyo/go-redpanda/service"
	"github.com/yoratyo/go-redpanda/storage"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func runSignalListener(cancel context.CancelFunc) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-sign
		fmt.Printf("Received %s, canceling context", s.String())
		cancel()
	}()
}

func main() {
	os.Exit(run())
}

func initStorage(ctx context.Context) *storage.Client {
	st, err := storage.New(ctx, storage.Options{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	return st
}

func run() int {
	ctx, cancel := context.WithCancel(context.Background())
	runSignalListener(cancel)
	errCh := make(chan error, 20)

	brokers := []string{os.Getenv("BROKER_HOST")}
	groupID := os.Getenv("CONSUMER_GROUP_ID")

	st := initStorage(ctx)
	defer st.Close()

	repo := repository.NewRepository(*st)
	svc := service.NewService(repo, nil)

	subscription := []consumer.Subscription{
		{
			Topic:    os.Getenv("TOPIC_NAME"),
			Listener: svc.ConsumeCrypto,
		},
	}

	log.Println("Starting Consumer...")
	client := consumer.NewConsumer(brokers, groupID, subscription)
	client.RunSubscription(ctx, errCh)

	select {
	case err := <-errCh:
		fmt.Printf("App ended with err: %s", err)
		cancel()
		return 1
	case <-ctx.Done():
		fmt.Println("Context cancelled.")
		return 0
	}
}
