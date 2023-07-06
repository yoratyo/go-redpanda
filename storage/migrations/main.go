package main

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/yoratyo/go-redpanda/storage"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	c, err := storage.New(ctx, storage.Options{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	if err := c.RunMigration(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	log.Println("Success migrate...")
}
